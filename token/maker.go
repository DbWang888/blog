package token

import "time"

type Maker interface {
	//通过用户信息和duration创建新的token
	CreateToken(username string, userid int, duration time.Duration) (string, error)

	//验证token是否有效，Payload
	VerifyToken(token string) (*Payload, error)
}
