package db

import (
	"blog/util"
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	RegisterTX(ctx context.Context, arg RegisterParams) (RegisterTxResult, error)
}

type SQLstore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLstore{
		db:      db,
		Queries: New(db),
	}
}

//建立事务
func (store *SQLstore) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err :%v, rb err :%v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

//注册结果返回参数
type RegisterTxResult struct {
	Tag  BlogTag  `json:"tag"`
	Auth BlogAuth `json:"auth"`
}

//注册所需参数
type RegisterParams struct {
	Auth BlogAuth `json:"auth"`
}

//注册事务，注册成功的同时在tag表中生成默认标签
func (store *SQLstore) RegisterTX(ctx context.Context, arg RegisterParams) (RegisterTxResult, error) {
	var result RegisterTxResult

	err := store.execTX(ctx, func(q *Queries) error {
		var err error

		//建立auth信息
		authResult, err := q.CreateAuth(ctx, CreateAuthParams{
			Username: arg.Auth.Username,
			Password: arg.Auth.Password,
		})

		if err != nil {
			return err
		}

		authResultID, err := authResult.LastInsertId()
		if err != nil {
			return err
		}

		result.Auth, err = q.GetAuthByID(ctx, int32(authResultID))
		if err != nil {
			return err
		}

		//建立tag信息
		tagResult, err := q.CreateBlogTag(ctx, CreateBlogTagParams{
			Name:      util.NewSqlNullString("默认"),
			CreatedBy: result.Auth.Username,
		})

		if err != nil {
			return err
		}

		tagResultID, err := tagResult.LastInsertId()
		if err != nil {
			return err
		}

		result.Tag, err = q.GetBlogTag(ctx, int32(tagResultID))
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}
