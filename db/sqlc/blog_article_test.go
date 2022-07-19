package db

import (
	"blog/util"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomArticle(t *testing.T) int32 {

	arg := CreateBlogArticleParams{
		TagID:     util.NewSqlNullInt32(int32(util.RandomInt(1, 8))),
		Title:     util.NewSqlNullString(util.RandomString(30)),
		Desc:      util.NewSqlNullString(util.RandomString(30)),
		Content:   util.NewSqlNullString(util.RandomString(30)),
		CreatedBy: util.NewSqlNullString(util.RandomString(30)),
	}

	result, err := testQueries.CreateBlogArticle(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	resultID, err := result.LastInsertId()
	require.NoError(t, err)

	blogArticle, err := testQueries.GetBlogArticles(context.Background(), int32(resultID))
	require.NoError(t, err)
	require.NotEmpty(t, blogArticle)

	require.Equal(t, arg.TagID, blogArticle.TagID)
	require.Equal(t, arg.Title, blogArticle.Title)
	require.Equal(t, arg.Desc, blogArticle.Desc)
	require.Equal(t, arg.Content, blogArticle.Content)
	require.Equal(t, arg.CreatedBy, blogArticle.CreatedBy)

	require.NotZero(t, blogArticle.CreatedOn)

	return int32(resultID)

}

func TestCreateBlogArticle(t *testing.T) {

	createRandomArticle(t)

}

func TestGetBlogArticles(t *testing.T) {

	blogArticle, err := testQueries.GetBlogArticles(context.Background(), int32(1))
	require.NoError(t, err)
	require.NotEmpty(t, blogArticle)

	require.Equal(t, blogArticle.ID, int32(1))
	require.Equal(t, blogArticle.Desc.String, "asjdfasdf")

}

func TestListBlogAtricles(t *testing.T) {

	// articleID := createRandomArticle(t)

	// article,err := testQueries.GetBlogTag(context.Background(),articleID)

	// require.NoError(t,err)
	// require.NotEmpty(t,article)

	arg := ListBlogAtriclesParams{
		CreatedBy: util.NewSqlNullString("jack"),
		Limit:     5,
		Offset:    0,
	}

	articles, err := testQueries.ListBlogAtricles(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, articles)

	for _, article := range articles {
		require.NotEmpty(t, article)
		require.Equal(t, "jack", article.CreatedBy.String)
	}

}

func TestUpdateBlogArticle(t *testing.T) {

	articleID1 := createRandomArticle(t)
	article1, err := testQueries.GetBlogArticles(context.Background(), articleID1)
	require.NoError(t, err)
	require.NotEmpty(t, article1)

	arg := UpdateBlogArticleParams{
		ID:         article1.ID,
		Title:      util.NewSqlNullString("new title"),
		Desc:       util.NewSqlNullString("new DESC"),
		Content:    util.NewSqlNullString("new content"),
		ModifiedBy: util.NewSqlNullString("manager"),
		ModifiedOn: time.Now(),
		TagID:      article1.TagID,
	}

	updateResult, err := testQueries.UpdateBlogArticle(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updateResult)

	// affectRows,err := updateResult.RowsAffected()
	// require.NoError(t,err)
	// require.Equal(t,affectRows,int64(4))

	newArticle, err := testQueries.GetBlogArticles(context.Background(), article1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, newArticle)

	require.Equal(t, arg.Content, newArticle.Content)
	require.Equal(t, arg.Desc, newArticle.Desc)
	require.Equal(t, arg.Title, newArticle.Title)
	require.Equal(t, arg.ModifiedBy, newArticle.ModifiedBy)

}

func TestDeleteArticle(t *testing.T) {

	articleID1 := createRandomArticle(t)
	article1, err := testQueries.GetBlogArticles(context.Background(), articleID1)
	require.NoError(t, err)
	require.NotEmpty(t, article1)
	require.Equal(t, util.NewSqlNullInt32(1), article1.State)

	arg := DeleteArticleParams{
		ID:        articleID1,
		DeletedOn: time.Now(),
	}

	delResult, err := testQueries.DeleteArticle(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, delResult)

	deleteArticle, err := testQueries.GetBlogArticles(context.Background(), articleID1)
	require.NoError(t, err)
	require.NotEmpty(t, deleteArticle)
	require.Equal(t, util.NewSqlNullInt32(0), deleteArticle.State)

}
