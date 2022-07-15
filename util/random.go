package util

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqretuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

//生成随机数,大小在min和max之间
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

//生成长度为N的随机字符串
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()

}

func RandomName(n int) sql.NullString {
	return NewSqlNullString(RandomString(n))
}
