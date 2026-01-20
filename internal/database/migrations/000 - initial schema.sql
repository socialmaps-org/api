CREATE TABLE user (
    -- same as their OpenStreetMap ID
    id INTEGER PRIMARY KEY,
    created INTEGER NOT NULL DEFAULT (my_unixepoch()),
    updated INTEGER NOT NULL DEFAULT (my_unixepoch()),
    display_name TEXT NOT NULL
);

CREATE TRIGGER enforce_user_constants
BEFORE UPDATE OF created ON user BEGIN
    SELECT raise(FAIL, 'cannot update the created column');
END;

CREATE TRIGGER update_user_on_update
AFTER UPDATE ON user
FOR EACH ROW BEGIN
    UPDATE user SET
        updated = my_unixepoch()
    WHERE
        id = old.id;
END;

CREATE TABLE place (
    id INTEGER PRIMARY KEY,
    created INTEGER NOT NULL DEFAULT (my_unixepoch()),
    updated INTEGER NOT NULL DEFAULT (my_unixepoch()),
    name TEXT NOT NULL,
    lat REAL NOT NULL,
    lon REAL NOT NULL,
    osm_type TEXT NOT NULL,
    osm_id INTEGER NOT NULL,

    n_likes INTEGER NOT NULL DEFAULT 0,
    n_dislikes INTEGER NOT NULL DEFAULT 0,

    -- decayed number of likes
    dec_n_likes REAL NOT NULL DEFAULT 0,
    -- decayed number of dislikes
    dec_n_dislikes REAL NOT NULL DEFAULT 0,
    -- decayed numbers last updated (recalculated) at
    dec_updated_at INTEGER NOT NULL DEFAULT (my_unixepoch()),

    score REAL NOT NULL DEFAULT 0.5,

    CONSTRAINT name CHECK (length(name) <= 256),
    CONSTRAINT lat CHECK (lat BETWEEN -90 AND +90),
    CONSTRAINT lon CHECK (lon BETWEEN -180 AND +180),
    CONSTRAINT osm_type CHECK (osm_type IN ('node', 'way', 'relation')),
    CONSTRAINT osm_id CHECK (osm_id >= 0),
    CONSTRAINT n_likes CHECK (n_likes >= 0),
    CONSTRAINT n_dislikes CHECK (n_dislikes >= 0),
    CONSTRAINT dec_n_likes CHECK (dec_n_likes BETWEEN 0 AND n_likes),
    CONSTRAINT dec_n_dislikes CHECK (
        dec_n_dislikes BETWEEN 0 AND n_dislikes
    )
);

CREATE UNIQUE INDEX place_osm ON place (osm_type, osm_id);

CREATE INDEX place_coord ON place (lat, lon);

CREATE TRIGGER enforce_place_constants
BEFORE UPDATE OF created ON place BEGIN
    SELECT raise(FAIL, 'cannot update the created column');
END;

CREATE TRIGGER update_place_score
AFTER UPDATE OF dec_n_likes, dec_n_dislikes ON place BEGIN
    UPDATE place SET
        score
        = (new.dec_n_likes + 1)
        / ((new.dec_n_likes + 1) + (new.dec_n_dislikes + 1))
    WHERE
        id = new.id;
END;

CREATE TRIGGER update_place_on_update
AFTER UPDATE ON place
FOR EACH ROW BEGIN
    UPDATE place SET
        updated = my_unixepoch()
    WHERE
        id = old.id;
END;


CREATE TABLE review (
    id INTEGER PRIMARY KEY,
    created INTEGER NOT NULL DEFAULT (my_unixepoch()),
    updated INTEGER NOT NULL DEFAULT (my_unixepoch()),
    place_id INTEGER NOT NULL REFERENCES place (id),
    user_id INTEGER NOT NULL REFERENCES user (id),
    liked BOOLEAN NOT NULL,
    comment TEXT,
    n_likes INTEGER NOT NULL DEFAULT 0,
    dec_n_likes REAL NOT NULL DEFAULT 0,
    -- decay updated (recalculated) at
    dec_updated_at INTEGER NOT NULL DEFAULT (my_unixepoch()),

    last_decision_at INTEGER,
    last_decision_by TEXT,
    last_decision_approved BOOLEAN,

    CONSTRAINT liked CHECK (liked IN (TRUE, FALSE)),
    CONSTRAINT n_likes CHECK (n_likes >= 0),
    CONSTRAINT dec_n_likes CHECK (dec_n_likes >= 0),
    CONSTRAINT last_decision_approved CHECK (
        last_decision_approved IN (TRUE, FALSE)
    )
);

CREATE INDEX review_latest_by_place
ON review (
    place_id,
    last_decision_approved,
    created DESC
);

CREATE INDEX review_top_by_place
ON review (
    place_id,
    last_decision_approved,
    n_likes DESC,
    created DESC
);

CREATE INDEX review_hot_by_place
ON review (
    place_id,
    last_decision_approved,
    dec_n_likes DESC,
    created DESC
);

CREATE TRIGGER enforce_review_constants
BEFORE UPDATE OF
created,
place_id,
user_id
ON review BEGIN
    SELECT
        raise(
            FAIL,
            'cannot update created, place_id, and/or user_id column(s)'
        )
    ;
END;

CREATE TRIGGER on_review_update
AFTER UPDATE ON review
FOR EACH ROW BEGIN
    UPDATE review SET
        updated = my_unixepoch()
    WHERE
        id = old.id;
END;

CREATE TRIGGER on_review_insert
AFTER INSERT ON review
FOR EACH ROW BEGIN
    UPDATE place SET
        n_likes
        = n_likes + if(new.liked, 1, 0),
        n_dislikes
        = n_dislikes + if(NOT new.liked, 1, 0),
        dec_n_likes
        = dec_n_likes
        * pow(2, -(new.updated - dec_updated_at) / 15552000.0)
        + if(new.liked, 1, 0),
        dec_n_dislikes
        = dec_n_dislikes
        * pow(2, -(new.updated - dec_updated_at) / 15552000.0)
        + if(NOT new.liked, 1, 0)
    WHERE
        id = new.place_id
    ;
    -- Update dec_updated_at separately so that you can refer to its previous
    -- value while calculating time elapsed (`new.updated-dec_updated_at`).
    UPDATE place SET
        dec_updated_at = new.updated
    WHERE
        id = new.place_id
    ;
END;

CREATE TRIGGER on_review_liked_update
AFTER UPDATE OF liked ON review
FOR EACH ROW BEGIN
    UPDATE places SET
        n_likes
        = n_likes + 1,
        n_dislikes
        = n_dislikes - 1,
        dec_n_likes
        = dec_n_likes
        + pow(2, -(dec_updated_at - old.created) / 15552000.0),
        dec_n_dislikes
        = dec_n_dislikes
        - pow(2, -(dec_updated_at - old.created) / 15552000.0)
    WHERE
        id = new.place_id
        AND NOT old.liked
        AND new.liked
    ;
    UPDATE places SET
        n_likes
        = n_likes - 1,
        n_dislikes
        = n_dislikes + 1,
        dec_n_likes
        = dec_n_likes
        - pow(2, -(dec_updated_at - old.created) / 15552000.0),
        dec_n_dislikes
        = dec_n_dislikes
        + pow(2, -(dec_updated_at - old.created) / 15552000.0)
    WHERE
        id = new.place_id
        AND old.liked
        AND NOT new.liked
    ;
END;

CREATE TRIGGER on_review_delete
AFTER DELETE ON review
FOR EACH ROW BEGIN
    UPDATE place SET
        n_dislikes
        = n_dislikes - 1,
        dec_n_dislikes
        = dec_n_dislikes
        - pow(2, -(dec_updated_at - old.created) / 15552000.0)
    WHERE
        id = old.place_id
        AND NOT old.liked
    ;
    UPDATE place SET
        n_likes
        = n_likes - 1,
        dec_n_likes
        = dec_n_likes - pow(2, -(dec_updated_at - old.created) / 15552000.0)
    WHERE
        id = old.place_id
        AND old.liked
    ;
END;


CREATE TABLE review_decision (
    id INTEGER PRIMARY KEY,
    created INTEGER NOT NULL DEFAULT (my_unixepoch()),
    review_id INTEGER NOT NULL REFERENCES review (id),
    moderator TEXT NOT NULL,
    approved BOOLEAN NOT NULL,
    details TEXT NOT NULL,

    CONSTRAINT approved CHECK (approved IN (TRUE, FALSE))
);

CREATE TRIGGER on_review_decision_insert
AFTER INSERT ON review_decision
FOR EACH ROW BEGIN
    UPDATE review SET
        last_decision_at = new.created,
        last_decision_by = new.moderator,
        last_decision_approved = new.approved
    WHERE
        id = new.review_id
    ;
END;

CREATE TRIGGER prohibit_review_decision_updates
BEFORE UPDATE ON review_decision BEGIN
    SELECT raise(FAIL, 'cannot update a ReviewDecision');
END;

CREATE TABLE review_like (
    id INTEGER PRIMARY KEY,
    created INTEGER NOT NULL DEFAULT (my_unixepoch()),
    review_id INTEGER NOT NULL REFERENCES review (id),
    user_id INTEGER NOT NULL REFERENCES user (id)
);

CREATE UNIQUE INDEX uniq_user_review_like ON review_like (
    review_id, user_id
);

CREATE TRIGGER prohibit_review_like_updates
BEFORE UPDATE ON review_like BEGIN
    SELECT raise(FAIL, 'cannot update a ReviewLike');
END;

CREATE TRIGGER on_review_like_insert
AFTER INSERT ON review_like
FOR EACH ROW BEGIN
    UPDATE review SET
        n_likes
        = n_likes + 1,
        dec_n_likes
        = dec_n_likes * pow(2, -(new.created - dec_updated_at) / 15552000.0)
        + 1
    WHERE
        id = new.review_id
    ;
END;

CREATE TRIGGER on_review_like_delete
AFTER DELETE ON review_like
FOR EACH ROW BEGIN
    UPDATE review SET
        n_likes = n_likes - 1,
        dec_n_likes
        = dec_n_likes - pow(2, -(dec_updated_at - old.created) / 15552000.0)
    WHERE
        id = old.review_id
    ;
END;

PRAGMA user_version = 1;
