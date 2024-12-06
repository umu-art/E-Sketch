package listener

import (
	"est-proxy/src/service"
	proxymodels "est_proxy_go/models"
	"fmt"

	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BoardListener struct {
	boardService service.BoardService
	userService  service.UserService
}

func NewBoardListener(boardService service.BoardService, userService service.UserService) *BoardListener {
	return &BoardListener{
		boardService: boardService,
		userService:  userService,
	}
}

func (b BoardListener) GetByUuid(ctx echo.Context) error {
	sessionUserId := service.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	boardId, err := uuid.Parse(ctx.Param("boardId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	board, statusError := b.boardService.GetByUuid(ctx.Request().Context(), sessionUserId, &boardId)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.JSON(http.StatusOK, board)
}

func (b BoardListener) List(ctx echo.Context) error {
	sessionUserId := service.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	list, statusError := b.boardService.List(ctx.Request().Context(), sessionUserId)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.JSON(http.StatusOK, list)
}

func (b BoardListener) Create(ctx echo.Context) error {
	sessionUserId := service.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	var createRequest proxymodels.CreateRequest
	err := ctx.Bind(&createRequest)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорректный запрос: %v", err))
	}

	board, statusError := b.boardService.Create(ctx.Request().Context(), sessionUserId, &createRequest)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.JSON(http.StatusOK, board)
}

func (b BoardListener) Update(ctx echo.Context) error {
	sessionUserId := service.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	boardId, err := uuid.Parse(ctx.Param("boardId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	var updateRequest proxymodels.CreateRequest
	if ctx.Bind(&updateRequest) != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорректный запрос: %v", err))
	}

	board, statusError := b.boardService.Update(ctx.Request().Context(), sessionUserId, &boardId, &updateRequest)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.JSON(http.StatusOK, board)
}

func (b BoardListener) DeleteBoard(ctx echo.Context) error {
	sessionUserId := service.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	boardId, err := uuid.Parse(ctx.Param("boardId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	statusError := b.boardService.DeleteBoard(ctx.Request().Context(), sessionUserId, &boardId)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.String(http.StatusOK, "Доска успешно удалена")
}

func (b BoardListener) Share(ctx echo.Context) error {
	sessionUserId := service.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	boardId, err := uuid.Parse(ctx.Param("boardId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	var shareBoardDto proxymodels.ShareBoardDto
	err = ctx.Bind(&shareBoardDto)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорректный запрос: %v", err))
	}

	statusError := b.boardService.Share(ctx.Request().Context(), sessionUserId, &boardId, &shareBoardDto)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.String(http.StatusOK, "Права доступа на доску успешно выданы")
}

func (b BoardListener) ChangeAccess(ctx echo.Context) error {
	sessionUserId := service.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	boardId, err := uuid.Parse(ctx.Param("boardId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	var shareBoardDto proxymodels.ShareBoardDto
	err = ctx.Bind(&shareBoardDto)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорректный запрос: %v", err))
	}

	statusError := b.boardService.ChangeAccess(ctx.Request().Context(), sessionUserId, &boardId, &shareBoardDto)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.String(http.StatusOK, "Права доступа на доску успешно изменены")
}

func (b BoardListener) Unshare(ctx echo.Context) error {
	sessionUserId := service.GetSessionUserId(ctx)
	if sessionUserId == nil {
		return ctx.String(http.StatusUnauthorized, "Отсутствует или некорректная сессия")
	}

	boardId, err := uuid.Parse(ctx.Param("boardId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорретный запрос: %v", err))
	}

	var unshareRequest proxymodels.UnshareRequest
	err = ctx.Bind(&unshareRequest)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Sprintf("Некорреткный запрос: %v", err))
	}

	statusError := b.boardService.Unshare(ctx.Request().Context(), sessionUserId, &boardId, &unshareRequest)
	if statusError != nil {
		return statusError.Send(ctx)
	}

	return ctx.String(http.StatusOK, "Права доступа на доску успешно изменены")
}

func (b BoardListener) Connect(ctx echo.Context) error {
	return ctx.String(http.StatusInternalServerError, "В разработке")
}
