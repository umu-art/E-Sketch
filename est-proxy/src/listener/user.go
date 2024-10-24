package listener

import (
	"github.com/labstack/echo/v4"
)

type UserListener struct {
}

func NewUserListener() *UserListener {
	return &UserListener{}
}

func (u UserListener) CheckSession(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u UserListener) GetSelf(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u UserListener) GetUserById(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u UserListener) Login(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u UserListener) Register(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u UserListener) Search(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}
