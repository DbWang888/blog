package db

import (
	"blog/util"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomBlogTag(t *testing.T) (tagID int64) {

	arg := CreateBlogTagParams{
		Name:      util.RandomName(6),
		CreatedBy: util.RandomName(6),
	}

	tag_result, err := testQueries.CreateBlogTag(context.Background(), arg)

	require.NoError(t, err)

	tagID, err = tag_result.LastInsertId()

	require.NoError(t, err)
	require.NotEmpty(t, tagID)

	tag, err := testQueries.GetBlogTag(context.Background(), int32(tagID))
	require.NoError(t, err)
	require.NotEmpty(t, tag)

	require.Equal(t, arg.Name.String, tag.Name.String)
	require.Equal(t, arg.CreatedBy, tag.CreatedBy)
	require.Equal(t, tag.State.Int32, int32(1))
	require.NotZero(t, tag.ID)
	require.NotZero(t, tag.ModifiedBy)
	require.NotZero(t, tag.CreatedOn)

	return
}

func TestCreateBlogTag(t *testing.T) {
	createRandomBlogTag(t)
}

func TestGetBlogTag(t *testing.T) {
	blogTag, err := testQueries.GetBlogTag(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, blogTag)
}

func TestUpdateBlogTag(t *testing.T) {

	tagID := createRandomBlogTag(t)

	arg := UpdateBlogTagParams{
		Name:       util.RandomName(6),
		ModifiedBy: util.RandomName(6),
		ModifiedOn: time.Now(),
		ID:         int32(tagID),
	}

	result, err := testQueries.UpdateBlogTag(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// resultID,err := result.LastInsertId()
	// require.NoError(t,err)
	// require.Equal(t,resultID,tagID)

	affectRows, err := result.RowsAffected()
	require.NoError(t, err)
	require.NotEmpty(t, affectRows)
	require.Equal(t, affectRows, int64(1))

	blogTag, err := testQueries.GetBlogTag(context.Background(), int32(tagID))
	require.NoError(t, err)
	require.NotEmpty(t, blogTag)

	require.Equal(t, arg.Name.String, blogTag.Name.String)
	require.Equal(t, arg.ModifiedBy.String, blogTag.ModifiedBy.String)
	require.NotZero(t, blogTag.ModifiedOn)

}

func TestDeleteBlogTag(t *testing.T) {

	tagID := createRandomBlogTag(t)

	arg := DeleteBlogTagParams{
		DeletedOn: time.Now(),
		ID:        int32(tagID),
	}

	result, err := testQueries.DeleteBlogTag(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	affectRows, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), affectRows)

	blogTag, err := testQueries.GetBlogTag(context.Background(), int32(tagID))
	require.NoError(t, err)
	require.NotEmpty(t, blogTag)

	require.Equal(t, blogTag.State.Int32, int32(0))
	require.NotZero(t, blogTag.DeletedOn)

}

func TestListBlogTag(t *testing.T) {

	var lastTagID int64

	for i := 0; i < 10; i++ {
		lastTagID = createRandomBlogTag(t)
	}

	lastTag, err := testQueries.GetBlogTag(context.Background(), int32(lastTagID))
	require.NoError(t, err)
	require.NotEmpty(t, lastTag)

	arg := ListBlogTagParams{
		CreatedBy: lastTag.CreatedBy,
		Limit:     5,
		Offset:    0,
	}

	blogTags, err := testQueries.ListBlogTag(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, blogTags)

	for _, blogTag := range blogTags {
		require.NotEmpty(t, blogTag)
		require.Equal(t, blogTag.CreatedBy, arg.CreatedBy)
	}

}
