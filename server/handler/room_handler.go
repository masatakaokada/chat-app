package handler

import (
	"app/model"
	"app/repository"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
)

func RoomIndex(c echo.Context) error {
	rooms, err := repository.RoomList()
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, rooms)
}

type RoomCreateOutput struct {
	Room             *model.Room
	Message          string
	ValidationErrors []string
}

type RoomCreateStruct struct {
	Name    string `json:"name"`
	UserIds []int  `json:"userIds"`
}

func RoomCreate(c echo.Context) error {
	token := c.Get("token").(*auth.Token)
	user, _ := repository.UserGetByFirebaseUid(token.UID)

	var room RoomCreateStruct
	var out RoomCreateOutput

	if err := c.Bind(&room); err != nil {
		// エラーの内容をサーバーのログに出力します。
		c.Logger().Error(err.Error())

		// リクエストの解釈に失敗した場合は 400 エラーを返却します。
		return c.JSON(http.StatusBadRequest, out)
	}

	_, err := repository.RoomCreate(room.Name, user, room.UserIds)
	if err != nil {
		// エラーの内容をサーバーのログに出力します。
		c.Logger().Error(err.Error())

		// サーバー内の処理でエラーが発生した場合は 500 エラーを返却します。
		return c.JSON(http.StatusInternalServerError, out)
	}

	return c.NoContent(http.StatusOK)
}
