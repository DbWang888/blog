package api

import (
	"blog/e"
	"blog/logger"
	"blog/token"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenmaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("授权未提供")
			c.AbortWithStatusJSON(http.StatusUnauthorized, e.GetErrResult(e.ERROR_AUTH_CHECK_TOKEN_FAIL, err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("无效的授权头格式")
			c.AbortWithStatusJSON(http.StatusUnauthorized, e.GetErrResult(e.ERROR_AUTH_CHECK_TOKEN_FAIL, err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("不支持该授权头格式")
			c.AbortWithStatusJSON(http.StatusUnauthorized, e.GetErrResult(e.ERROR_AUTH_CHECK_TOKEN_FAIL, err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenmaker.VerifyToken(accessToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, e.GetErrResult(e.ERROR_AUTH_TOKEN, err))
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}

	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWriter

		beginTime := time.Now().Unix()
		c.Next()

		endtime := time.Now().Unix()

		fields := logger.Fields{
			"request":  c.Request.PostForm.Encode(),
			"response": bodyWriter.body.String(),
		}
		e.Logger.WithFields(fields).Infof("access log: method: %s,status_code: %d, begin_time: %d, end_time: %d",
			c.Request.Method,
			bodyWriter.Status(),
			beginTime,
			endtime,
		)
	}
}
