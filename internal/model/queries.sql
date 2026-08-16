-- name: CreatePlace :one
INSERT INTO socialmaps.place (
    "name",
    "location",
    osm_type,
    osm_id,
    created,
    updated
)
VALUES (
    @name,
    ST_POINT(@lon::double precision, @lat::double precision, 4326),
    @osm_type,
    @osm_id,
    @as_of,
    @as_of
)
RETURNING *;

-- name: QueryPlaces :many
SELECT DISTINCT ON (elm."name") -- TODO: buggy. we want to allow same named POIs if they are apart enough
    sqlc.embed(elm),  -- noqa
    sqlc.embed(plc)  -- noqa
FROM
    osm2pgsql.element AS elm
LEFT OUTER JOIN socialmaps.optional_computed_place AS plc
    ON (
        ST_DWITHIN(plc."location"::geography, elm."location"::geography, 10)
        AND elm."name" = plc."name"
        AND plc."location" && ST_MAKEENVELOPE(
            @lon_min::double precision,
            @lat_min::double precision,
            @lon_max::double precision,
            @lat_max::double precision,
            4326
        )
    )
WHERE
    elm."location" && ST_MAKEENVELOPE(
        @lon_min::double precision,
        @lat_min::double precision,
        @lon_max::double precision,
        @lat_max::double precision,
        4326
    )
    AND elm.tags @@ @predicate::jsonpath
LIMIT -- noqa: AM09
    100; -- TODO: we'll add order by later

-- name: LookupPlace :one
SELECT
    sqlc.embed(elm),  -- noqa
    sqlc.embed(plc)  -- noqa
FROM
    osm2pgsql.element AS elm
CROSS JOIN
    LATERAL (  -- noqa: ST05
        SELECT SIMILARITY(@name::varchar, tag."value") AS best
        FROM JSONB_EACH_TEXT(elm.tags) AS tag ("key", "value")
        WHERE tag."key" IN ('name', 'int_name') OR tag."key" LIKE 'name:%'
    ) AS sim
LEFT OUTER JOIN socialmaps.optional_computed_place AS plc
    ON (
        ST_DWITHIN(plc."location"::geography, elm."location"::geography, 10)
        AND elm."name" = plc."name"
        AND plc."location" && ST_MAKEENVELOPE(
            @lon_min::double precision,
            @lat_min::double precision,
            @lon_max::double precision,
            @lat_max::double precision,
            4326
        )
    )
WHERE
    elm."location" && ST_MAKEENVELOPE(
        @lon_min::double precision,
        @lat_min::double precision,
        @lon_max::double precision,
        @lat_max::double precision,
        4326
    )
    AND sim.best > 0.3
ORDER BY
    sim.best DESC
LIMIT -- noqa: AM09
    1;

-- name: LoadPlace :one
SELECT
    sqlc.embed(elm),  -- noqa
    sqlc.embed(plc)  -- noqa
FROM
    socialmaps.computed_place AS plc
INNER JOIN osm2pgsql.element AS elm
    ON (
        ST_DWITHIN(plc."location"::geography, elm."location"::geography, 10)
        AND plc."name" = elm."name"
    )
WHERE
    plc.id = $1;

-- name: CreateReview :one
INSERT INTO socialmaps.review (
    place_id,
    user_id,
    rating,
    "comment",
    created,
    reviewed_at,
    updated
)
VALUES (
    @place_id,
    @user_id,
    @rating,
    @comment,
    @as_of,
    @reviewed_at,
    @as_of
) RETURNING *;

-- name: DeleteReview :exec
DELETE FROM socialmaps.review
WHERE
    id = $1;

-- name: UpdateReview :one
UPDATE socialmaps.review SET
    rating = $1,
    "comment" = $2,
    updated = @as_of
WHERE
    id = $3
RETURNING *;


-- name: LoadReview :one
SELECT *
FROM socialmaps.review
WHERE
    id = $1;

-- name: ListLatestApprovedReviewsOfPlace :many
SELECT
    sqlc.embed(rvw), -- noqa
    sqlc.embed(usr) -- noqa
FROM socialmaps.review AS rvw
INNER JOIN socialmaps."user" AS usr ON rvw.user_id = usr.id
WHERE
    rvw.place_id = @place_id
    AND rvw.last_decision_approved
    AND rvw.created <= TO_TIMESTAMP(@last_created)
    AND rvw.id < @last_id
ORDER BY
    rvw.created DESC,
    rvw.id DESC
LIMIT
    @lmt::bigint;

-- name: ListLatestApprovedReviewsOfPlaceReverse :many
SELECT
    sqlc.embed(rvw), -- noqa
    sqlc.embed(usr) -- noqa
FROM socialmaps.review AS rvw
INNER JOIN socialmaps."user" AS usr ON rvw.user_id = usr.id
WHERE
    rvw.place_id = @place_id
    AND rvw.last_decision_approved
    AND rvw.created >= TO_TIMESTAMP(@first_created)
    AND rvw.id > @first_id
ORDER BY
    rvw.created DESC,
    rvw.id DESC
LIMIT
    @lmt::bigint;

-- name: ListEarliestUnapprovedReviews :many
SELECT *
FROM socialmaps.review
WHERE
    last_decision_approved IS NULL
    AND created >= @last_created
    AND id > @last_id
ORDER BY
    created ASC,
    id ASC
LIMIT
    @lim;

-- name: CreateUser :one
INSERT INTO socialmaps."user" (
    created,
    updated,
    id,
    display_name
) VALUES (
    @as_of,
    @as_of,
    @id,
    @display_name
)
ON CONFLICT (id) DO UPDATE SET display_name = excluded.display_name
RETURNING *;

-- name: CreateReviewDecision :one
INSERT INTO socialmaps.review_decision (
    created,
    review_id,
    moderator,
    approved,
    details
) VALUES (
    @as_of,
    @review_id,
    @moderator,
    @approved,
    @details
) RETURNING *;

-- name: LoadLatestDecisionOfReview :one
SELECT *
FROM socialmaps.review_decision
WHERE
    review_id = $1
ORDER BY
    created DESC
LIMIT 1;
