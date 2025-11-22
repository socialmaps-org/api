CREATE TABLE "Users" (
      "created"         INTEGER NOT NULL DEFAULT (my_unixepoch())
    , "updated"         INTEGER NOT NULL DEFAULT (my_unixepoch())
    , "id"              TEXT    PRIMARY KEY
    , "osm_id"          INTEGER NOT NULL
    , "display_name"    TEXT    NOT NULL

    , CONSTRAINT "id"     CHECK ("id" LIKE 'usr_%')
) STRICT;

CREATE UNIQUE INDEX users_osm ON Users (osm_id);

CREATE TRIGGER "enforce_user_constants"
BEFORE UPDATE OF "created" ON "Users" BEGIN
    SELECT raise(FAIL, 'cannot update the "created" column');
END;

CREATE TRIGGER "update_user_on_update"
AFTER UPDATE ON "Users"
FOR EACH ROW BEGIN
    UPDATE "Users" SET
        "updated" = my_unixepoch()
    WHERE
        "id" = NEW."id"
    ;
END;


CREATE TABLE "Places" (
      "created"        INTEGER NOT NULL DEFAULT (my_unixepoch())
    , "updated"        INTEGER NOT NULL DEFAULT (my_unixepoch())
    , "id"             TEXT    PRIMARY KEY
    , "name"           TEXT    NOT NULL
    , "lat"            REAL    NOT NULL
    , "lon"            REAL    NOT NULL
    , "osm_type"       TEXT    NOT NULL
    , "osm_id"         INTEGER NOT NULL
    , "n_likes"        INTEGER NOT NULL DEFAULT 0
    , "n_dislikes"     INTEGER NOT NULL DEFAULT 0
    -- decayed number of likes
    , "dec_n_likes"    REAL    NOT NULL DEFAULT 0
    -- decayed number of dislikes
    , "dec_n_dislikes" REAL    NOT NULL DEFAULT 0
    -- decayed numbers last updated (recalculated) at
    , "dec_updated_at" INTEGER NOT NULL DEFAULT (my_unixepoch())

    , CONSTRAINT "id"             CHECK ("id" LIKE 'plc_%')
    , CONSTRAINT "name"           CHECK (length("name") <= 256)
    , CONSTRAINT "lat"            CHECK ("lat" BETWEEN  -90 AND  +90)
    , CONSTRAINT "lon"            CHECK ("lon" BETWEEN -180 AND +180)
    , CONSTRAINT "osm_type"       CHECK ("osm_type" IN ('node', 'way', 'relation'))
    , CONSTRAINT "osm_id"         CHECK ("osm_id" >= 0)
    , CONSTRAINT "n_likes"        CHECK ("n_likes" >= 0)
    , CONSTRAINT "n_dislikes"     CHECK ("n_dislikes" >= 0)
    , CONSTRAINT "dec_n_likes"    CHECK ("dec_n_likes" BETWEEN 0 AND "n_likes")
    , CONSTRAINT "dec_n_dislikes" CHECK ("dec_n_dislikes" BETWEEN 0 AND "n_dislikes")
) STRICT;

CREATE UNIQUE INDEX places_osm ON Places (osm_type, osm_id);

CREATE INDEX places_by_lat_lon ON Places (lat, lon);

CREATE TRIGGER enforce_place_constants
BEFORE UPDATE OF created ON Places BEGIN
    SELECT raise(FAIL, 'cannot update the "created" column');
END;

CREATE TRIGGER update_place_on_update
AFTER UPDATE ON Places
FOR EACH ROW BEGIN
    UPDATE Places SET
        updated = my_unixepoch()
    WHERE
        id = NEW.id
    ;
END;


CREATE TABLE Reviews (
      created        INTEGER NOT NULL DEFAULT (my_unixepoch())
    , updated        INTEGER NOT NULL DEFAULT (my_unixepoch())
    , id             TEXT    PRIMARY KEY
    , place_id       TEXT    NOT NULL REFERENCES Places (id)
    , user_id        TEXT    NOT NULL REFERENCES Users (id)
    , liked          INTEGER NOT NULL
    , comment        TEXT
    , n_likes        INTEGER NOT NULL DEFAULT 0
    , dec_n_likes    REAL    NOT NULL DEFAULT 0
    -- decay updated (recalculated) at
    , dec_updated_at INTEGER NOT NULL DEFAULT (my_unixepoch())

    , CONSTRAINT "id"          CHECK ("id" LIKE 'rvw_%')
    , CONSTRAINT "liked"       CHECK ("liked" IN (TRUE, FALSE))
    , CONSTRAINT "n_likes"     CHECK ("n_helpful" >= 0)
    , CONSTRAINT "dec_n_likes" CHECK ("dec_n_helpful" >= 0)
) STRICT;

CREATE INDEX latest_reviews_by_place
ON Reviews (
      place_id
    , created DESC
);

CREATE INDEX most_liked_reviews_by_place
ON Reviews (
      place_id
    , n_likes  DESC
    , created  DESC
);

CREATE TRIGGER enforce_review_constants
BEFORE UPDATE OF
      created
    , place_id
    , user_id
ON Reviews BEGIN
    SELECT
        raise(
            FAIL,
            'cannot update "created", "place_id", and/or "user_id" column(s)'
        )
    ;
END;

CREATE TRIGGER on_review_insert
AFTER INSERT ON Reviews
FOR EACH ROW BEGIN
    UPDATE Places SET
          "n_likes"
            = "n_likes" + if(NEW.liked, 1, 0)
        , "n_dislikes"
            = "n_dislikes" + if(NOT NEW.liked, 1, 0)
        , "dec_n_likes"
            = "dec_n_likes" * pow(2, -(my_unixepoch()-"dec_updated_at")/15552000.0) + if(NEW.liked, 1, 0)
        , "dec_n_dislikes"
            = "dec_n_dislikes" * pow(2, -(my_unixepoch()-"dec_updated_at")/15552000.0) + if(NOT NEW.liked, 1, 0)
    WHERE
        id = NEW.place_id
    ;
    -- Update "dec_updated_at" separately so that you can refer to its previous
    -- value while calculating time elapsed (`my_unixepoch()-dec_updated_at`).
    UPDATE Places SET
        dec_updated_at = my_unixepoch()
    WHERE
        id = NEW.place_id
    ;
END;

CREATE TRIGGER on_review_update
AFTER UPDATE ON Reviews
FOR EACH ROW BEGIN
    UPDATE Reviews SET
        updated = my_unixepoch()
    WHERE
        id = NEW.id
    ;
END;

CREATE TRIGGER on_review_liked_update
AFTER UPDATE OF liked ON Reviews
FOR EACH ROW BEGIN
    UPDATE Places SET
          n_likes
            = n_likes + 1
        , n_dislikes
            = n_dislikes - 1
        , dec_n_likes
            = dec_n_likes + pow(2, -(dec_updated_at-OLD.created)/15552000.0)
        , dec_n_dislikes
            = dec_n_dislikes - pow(2, -(dec_updated_at-OLD.created)/15552000.0)
    WHERE
        id = NEW.place_id
        AND OLD.liked = FALSE
        AND NEW.liked = TRUE
    ;
    UPDATE Places SET
          n_likes
            = n_likes - 1
        , n_dislikes
            = n_dislikes + 1
        , dec_n_likes
            = dec_n_likes - pow(2, -(dec_updated_at-OLD.created)/15552000.0)
        , dec_n_dislikes
            = dec_n_dislikes + pow(2, -(dec_updated_at-OLD.created)/15552000.0)
    WHERE
        id = NEW.place_id
        AND OLD.liked = TRUE
        AND NEW.liked = FALSE
    ;
END;

CREATE TRIGGER on_review_delete
AFTER DELETE ON Reviews
FOR EACH ROW BEGIN
    UPDATE Places SET
          n_dislikes
            = n_dislikes - 1
        , dec_n_dislikes
            = dec_n_dislikes - pow(2, -(dec_updated_at-OLD.created)/15552000.0)
    WHERE
        id = OLD.place_id
        AND OLD.liked = FALSE
    ;
    UPDATE Places SET
          n_likes
            = n_likes - 1
        , dec_n_likes
            = dec_n_likes - pow(2, -(dec_updated_at-OLD.created)/15552000.0)
    WHERE
        id = OLD.place_id
        AND OLD.liked = TRUE
    ;
END;


CREATE TABLE ReviewLikes (
      created   INTEGER NOT NULL DEFAULT (my_unixepoch())
    , updated   INTEGER NOT NULL DEFAULT (my_unixepoch())
    , id        INTEGER PRIMARY KEY
    , review_id TEXT    NOT NULL REFERENCES Reviews (id)
    , user_id   TEXT    NOT NULL REFERENCES Users (id)
) STRICT;

CREATE UNIQUE INDEX uniq_user_review_like ON ReviewLikes (review_id, user_id);

CREATE TRIGGER on_reviewlike_insert
AFTER INSERT ON ReviewLikes
FOR EACH ROW BEGIN
    UPDATE Reviews SET
          n_likes
            = n_likes + 1
        , dec_n_likes
            = dec_n_likes * pow(2, -(my_unixepoch()-dec_updated_at)/15552000.0) + 1
    WHERE
        id = NEW.review_id
    ;
END;

CREATE TRIGGER on_reviewlike_delete
AFTER DELETE ON ReviewLikes
FOR EACH ROW BEGIN
    UPDATE Reviews SET
        n_likes = n_likes - 1
    WHERE
        id = OLD.review_id
    ;
END;

PRAGMA user_version = 1;
