-- name: CreatePlace :one
INSERT INTO socialmaps.place (
    "name",
    lat,
    lon,
    osm_type,
    osm_id,
    created,
    updated,
    dec_updated_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    @as_of,
    @as_of,
    @as_of
)
RETURNING *;

-- name: ListPlacesByCoord :many
SELECT *
FROM
    socialmaps.place
WHERE
    @lat_min <= lat AND lat <= @lat_max
    AND @lon_min <= lon AND lon <= @lon_max;

-- name: LoadPlace :one
SELECT *
FROM
    socialmaps.place
WHERE
    id = $1;

-- name: CreateReview :one
INSERT INTO socialmaps.review (
    place_id,
    user_id,
    liked,
    "comment",
    created,
    updated,
    dec_updated_at
)
VALUES (
    @place_id,
    @user_id,
    @liked,
    @comment,
    @as_of,
    @as_of,
    @as_of
) RETURNING *;

-- name: DeleteReview :exec
DELETE FROM socialmaps.review
WHERE
    id = $1;

-- name: UpdateReview :one
UPDATE socialmaps.review SET
    liked = $1,
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

-- name: ListHottestApprovedReviewsOfPlace :many
SELECT *
FROM socialmaps.review
WHERE
    place_id = $1
    AND last_decision_approved
    AND dec_n_likes <= @last_dec_n_likes
    AND id < @last_id
ORDER BY
    dec_n_likes DESC,
    id DESC
LIMIT
    $2;

-- name: ListLatestApprovedReviewsOfPlace :many
SELECT
    sqlc.embed(rvw), -- noqa
    sqlc.embed(usr) -- noqa
FROM socialmaps.review AS rvw
INNER JOIN socialmaps."user" AS usr ON rvw.user_id = usr.id
WHERE
    rvw.place_id = @place_id
    AND rvw.last_decision_approved
    AND rvw.created <= to_timestamp(@last_created)
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
    AND rvw.created >= to_timestamp(@first_created)
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

-- name: LikeReview :exec
INSERT INTO socialmaps.review_like (
    review_id,
    user_id,
    created
) VALUES (
    @review_id,
    @user_id,
    @as_of
)
ON CONFLICT DO NOTHING;

-- name: UnlikeReview :exec
DELETE FROM socialmaps.review_like
WHERE
    review_id = $1
    AND user_id = $2;

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
