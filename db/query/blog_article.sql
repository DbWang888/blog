-- name: CreateBlogArticle :execresult
INSERT INTO blog_article (
    tag_id,title,`desc`,content,created_by
) VALUES(
    ?,?,?,?,?
);

-- name: UpdateBlogArticle :execresult
UPDATE blog_article
SET tag_id=?,title=?,`desc`=?,content=?,modified_by=?,modified_on=?
WHERE id = ?;

-- name: GetBlogArticles :one
SELECT * FROM blog_article
WHERE id = ? LIMIT 1;

-- name: ListBlogAtricles :many
SELECT * FROM blog_article
WHERE created_by = ?
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: DeleteArticle :execresult
UPDATE blog_article
SET state = 0,deleted_on=?
WHERE id = ?;