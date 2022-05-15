package handler

import (
	"app/repository"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
)

// ログインユーザーを除いたユーザーのリスト
func RoomCreationUsers(c echo.Context) error {
	token := c.Get("token").(*auth.Token)
	user, _ := repository.UserGetByFirebaseUid(token.UID)

	users, err := repository.RoomCreationUsers(user.ID)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
