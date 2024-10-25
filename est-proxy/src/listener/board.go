package listener

import (
	estbackapi "est_back_go"
	"est_proxy_go/models"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BoardListener struct {
	boardApi *estbackapi.BoardAPIService
}

func NewBoardListener(boardApi *estbackapi.BoardAPIService) *BoardListener {
	return &BoardListener{
		boardApi: boardApi,
	}
}

func (b BoardListener) GetByUuid(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b BoardListener) List(ctx echo.Context) error {
	// TODO: contexts and user
	userId := uuid.New()

	list, _, err := b.boardApi.ListByUserId(ctx.Request().Context(), userId.String()).Execute()
	if err != nil {
		return fmt.Errorf("failed to list boards: %w", err)
	}

	return ctx.JSON(200, mapManyToProxy(list))
}

func (b BoardListener) Create(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b BoardListener) Update(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b BoardListener) DeleteBoard(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b BoardListener) Share(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b BoardListener) ChangeAccess(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b BoardListener) Unshare(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b BoardListener) Connect(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func mapManyToProxy(list *estbackapi.BackBoardListDto) models.BoardListDto {
	mine := make([]models.BoardDto, 0)
	shared := make([]models.BoardDto, 0)

	for _, dto := range list.Mine {
		mine = append(mine, mapToProxy(&dto))
	}

	for _, dto := range list.Shared {
		shared = append(shared, mapToProxy(&dto))
	}

	return models.BoardListDto{
		Mine:   mine,
		Shared: shared,
	}
}

func mapToProxy(dto *estbackapi.BackBoardDto) models.BoardDto {
	return models.BoardDto{
		Id: dto.Id,
		// TODO: map other fields
	}
}
