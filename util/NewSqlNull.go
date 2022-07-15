package util

import "database/sql"

func NewSqlNullString(str string) sql.NullString {
	return sql.NullString{
		String: str,
		Valid:  true,
	}
}

func NewSqlNullInt32(n int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: n,
		Valid: true,
	}
}
