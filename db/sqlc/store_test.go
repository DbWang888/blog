package db

import (
	"blog/util"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRegisterTX(t *testing.T) {

	store := NewStore(testDB)

	n := 10

	errs := make(chan error)
	results := make(chan RegisterTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.RegisterTX(ctx, RegisterParams{
				Auth: BlogAuth{
					Username: util.RandomName(8),
					Password: util.RandomName(6),
				},
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		require.Equal(t, result.Auth.Username.String, result.Tag.CreatedBy.String)
		require.WithinDuration(t, result.Auth.CreatedOn, result.Tag.CreatedOn, time.Second)
	}

}
