package handler

import (
	"app/model"
	"app/repository"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
)

func RoomIndex(c echo.Context) error {
	token := c.Get("token").(*auth.Token)
	user, _ := repository.UserGetByFirebaseUid(token.UID)

	rooms, err := repository.RoomList(user.ID)
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

func RoomCreate(c echo.Context) error {
	token := c.Get("token").(*auth.Token)
	user, _ := repository.UserGetByFirebaseUid(token.UID)

	var roomCreate model.RoomCreate
	var out RoomCreateOutput

	if err := c.Bind(&roomCreate); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, out)
	}

	// バリデーション
	if err := c.Validate(&roomCreate); err != nil {
		c.Logger().Error(err.Error())

		out.ValidationErrors = roomCreate.ValidationErrors(err)

		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	_, err := repository.RoomCreate(roomCreate.Name, user, roomCreate.UserIds)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, out)
	}

	return c.NoContent(http.StatusOK)
}
