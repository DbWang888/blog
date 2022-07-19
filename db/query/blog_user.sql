-- name: CreateAuth :execresult
INSERT INTO blog_auth (
    username,`password`,created_on
) VALUES(
    ?,?,?
);

-- name: GetAuthByID :one
SELECT * FROM blog_auth
WHERE id = ?;


-- name: GetAuthByUserName :one
SELECT * FROM blog_auth
WHERE username = ?;