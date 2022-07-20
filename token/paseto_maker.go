package token

import (
	"blog/e"
	"errors"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PastoMaker struct {
	paseto       *paseto.V2
	symmetrickey []byte
}

func NewPasetoMaker(symmetrickey string) (*PastoMaker, error) {
	if len(symmetrickey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("对称密钥长度错误,应该为：%v", chacha20poly1305.KeySize)
	}

	maker := &PastoMaker{
		paseto:       paseto.NewV2(),
		symmetrickey: []byte(symmetrickey),
	}

	return maker, nil
}

//创建token
func (maker *PastoMaker) CreateToken(username string, userid int, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, userid, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetrickey, payload, nil)
}

//验证token
func (maker *PastoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetrickey, payload, nil)
	if err != nil {
		return nil, errors.New(e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL))
	}

	err = payload.Vaild()

	if err != nil {
		return nil, err
	}
	return payload, nil
}
