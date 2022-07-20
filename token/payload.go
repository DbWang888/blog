package token

import (
	"blog/e"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"username"`
	IssueAt   time.Time `json:"issue_at"`   //建立时间
	ExpiredAt time.Time `json:"expried_at"` //过期时间
}

//创建新的令牌负载，利用特定的uerid uername 以及持续时间
func NewPayload(Username string, UserID int, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    UserID,
		UserName:  Username,
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload,nil
}


func (payload *Payload) Vaild() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New(e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT))
	}
	return nil
}