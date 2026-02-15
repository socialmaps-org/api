-- name: CreatePlace :one
INSERT INTO place (
    name,
    lat,
    lon,
    osm_type,
    osm_id
)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;

-- name: ListPlacesByCoord :many
SELECT *
FROM
    place
WHERE
    @lat_min <= lat AND lat <= @lat_max
    AND @lon_min <= lon AND lon <= @lon_max;

-- name: LoadPlace :one
SELECT *
FROM
    place
WHERE
    id = ?;

-- name: CreateReview :one
INSERT INTO review (
    place_id,
    user_id,
    liked,
    comment
)
VALUES (
    ?,
    ?,
    ?,
    ?
) RETURNING *;

-- name: DeleteReview :exec
DELETE FROM review
WHERE
    id = ?;

-- name: UpdateReview :one
UPDATE review SET
    liked = ?,
    comment = ?
WHERE
    id = ?
RETURNING *;


-- name: LoadReview :one
SELECT *
FROM review
WHERE
    id = ?;

-- name: ListHottestApprovedReviewsOfPlace :many
SELECT *
FROM review
WHERE
    place_id = ?
    AND last_decision_approved
    AND dec_n_likes <= @last_dec_n_likes
    AND id < @last_id
ORDER BY
    dec_n_likes DESC,
    id DESC
LIMIT
    ?;

-- name: ListLatestApprovedReviewsOfPlace :many
SELECT
    sqlc.embed(review), -- noqa
    sqlc.embed(user) -- noqa
FROM review
INNER JOIN user ON review.user_id = user.id
WHERE
    review.place_id = ?
    AND review.last_decision_approved
    AND review.created <= @last_created
    AND review.id < @last_id
ORDER BY
    review.created DESC,
    review.id DESC
LIMIT
    ?;

-- name: ListLatestApprovedReviewsOfPlaceReverse :many
SELECT
    sqlc.embed(review), -- noqa
    sqlc.embed(user) -- noqa
FROM review
INNER JOIN user ON review.user_id = user.id
WHERE
    review.place_id = ?
    AND review.last_decision_approved
    AND review.created >= @first_created
    AND review.id > @first_id
ORDER BY
    review.created DESC,
    review.id DESC
LIMIT
    ?;

-- name: ListEarliestUnapprovedReviews :many
SELECT *
FROM review
WHERE
    last_decision_approved IS NULL
    AND created >= @last_created
    AND id > @last_id
ORDER BY
    created ASC,
    id ASC
LIMIT
    @limit;

-- name: LikeReview :exec
INSERT INTO review_like (
    review_id,
    user_id

) VALUES (
    ?,
    ?
)
ON CONFLICT DO NOTHING;

-- name: UnlikeReview :exec
DELETE FROM review_like
WHERE
    review_id = ?
    AND user_id = ?;

-- name: CreateUser :one
INSERT INTO user (
    id,
    display_name
) VALUES (
    ?,
    ?
)
ON CONFLICT (id) DO UPDATE SET display_name = excluded.display_name
RETURNING *;

-- name: CreateReviewDecision :one
INSERT INTO review_decision (
    review_id,
    moderator,
    approved,
    details
) VALUES (
    ?,
    ?,
    ?,
    ?
) RETURNING *;

-- name: LoadLatestDecisionOfReview :one
SELECT *
FROM review_decision
WHERE
    review_id = ?
ORDER BY
    created DESC
LIMIT 1;
