package handler

import (
	"app/model"
	"app/repository"
	"fmt"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
)

// ログインユーザーを除いたユーザーのリスト
func UserIndex(c echo.Context) error {
	token := c.Get("token").(*auth.Token)
	user, _ := repository.UserGetByFirebaseUid(token.UID)

	users, err := repository.Users(user.ID)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func UserShow(c echo.Context) error {
	token := c.Get("token").(*auth.Token)

	user, err := repository.UserGetByFirebaseUid(token.UID)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
}

func UserCreate(c echo.Context) error {
	token := c.Get("token").(*auth.Token)

	user, _ := repository.UserGetByFirebaseUid(token.UID)

	if user == nil {
		user := &model.User{
			Email:       token.Claims["email"].(string),
			FirebaseUid: token.UID,
		}

		if err := c.Bind(&user); err != nil {
			c.Logger().Error(err.Error())

			return c.NoContent(http.StatusBadRequest)
		}

		res, err := repository.UserCreate(user)
		if err != nil {
			c.Logger().Error(err.Error())

			return c.NoContent(http.StatusInternalServerError)
		}

		id, _ := res.LastInsertId()
		fmt.Printf("ユーザーの作成に成功しました。 id: %v\n", id)
	}

	return c.NoContent(http.StatusOK)
}
