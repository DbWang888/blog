-- name: CreateBlogTag :execresult
INSERT INTO blog_tag (
    `name`,created_by
) VALUES(
    ?,?
);

-- name: UpdateBlogTag :execresult
UPDATE blog_tag
SET `name` = ?, modified_by = ?, modified_on = ?
WHERE id = ?;

-- name: DeleteBlogTag :execresult
UPDATE blog_tag
SET `state` = 0, deleted_on = ?
WHERE id = ?;

-- name: GetBlogTag :one
SELECT * FROM blog_tag
WHERE id=? LIMIT 1;

-- name: ListBlogTag :many
SELECT * FROM blog_tag
WHERE created_by = ?
ORDER BY id
LIMIT ?,?;