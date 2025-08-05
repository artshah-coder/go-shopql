package session

import (
	"context"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type Session struct {
	UserID uint32
}

type Token string

type uid int

const userID uid = 1

func UIDToContext(sessions SessionsStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		token = strings.TrimPrefix(token, "Token ")
		if token != "" {
			sess, err := sessions.Get(Token(token))
			if err != nil {
				log.Println(err)
				return
			}
			ctx := context.WithValue(c.Request.Context(), userID, sess.UserID)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}

func UIDFromContext(ctx context.Context) (uint32, bool) {
	uid, ok := ctx.Value(userID).(uint32)
	return uid, ok
}
