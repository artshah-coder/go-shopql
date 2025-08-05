package user

import (
	"errors"
	"net/http"

	"github.com/artshah-coder/go-shopql/internal/session"
	"github.com/artshah-coder/go-shopql/internal/utils/password"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Users    UsersStorage
	Sessions session.SessionsStorage
}

type RegisterResponse struct {
	Body struct {
		Token session.Token `json:"token"`
	} `json:"body"`
}

type RegisterRequest struct {
	User struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"user"`
}

func (uh *UserHandler) Register(c *gin.Context) {
	tmpUser := RegisterRequest{}
	if err := c.BindJSON(&tmpUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
		})
		return
	}
	salt := password.MakeSalt(8)
	pass := password.HashPass(tmpUser.User.Password, salt)

	user := User{
		Email:    tmpUser.User.Email,
		Username: tmpUser.User.Username,
		password: pass,
	}
	addedUser, err := uh.Users.Add(user)
	if err != nil {
		if errors.Is(err, ErrEmailExists) || errors.Is(err, ErrUsernameExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
		})
		return
	}

	sess := &session.Session{UserID: addedUser.ID}

	token, err := uh.Sessions.Add(sess)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}
	var resp RegisterResponse
	resp.Body.Token = token

	c.JSON(http.StatusOK, resp)
}
