package listener

import (
	"est-proxy/src/service"
	estbackapi "est_back_go"
	"est_proxy_go/models"
	"fmt"
	"time"

	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BoardListener struct {
	boardApi *estbackapi.BoardAPIService
	userRepository *service.UserRepository
}

func NewBoardListener(boardApi *estbackapi.BoardAPIService, userRepository *service.UserRepository) *BoardListener {
	return &BoardListener{
		boardApi: boardApi,
		userRepository: userRepository,
	}
}

func (b BoardListener) GetByUuid(ctx echo.Context) error {
	userId, sessionStatus := b.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	boardId := ctx.QueryParam("boardId")
	board, err := b.getBoard(ctx, boardId)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Доска не найдена")
	}
	if b.checkAccess(userId, board) != nil {
		return ctx.String(http.StatusForbidden, "Недостаточно прав")
	}

	return ctx.JSON(http.StatusOK, board)
}

func (b BoardListener) List(ctx echo.Context) error {
	userId, sessionStatus := b.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	list, _, err := b.boardApi.ListByUserId(ctx.Request().Context(), userId.String()).Execute()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Не получилось составить список досок")
	}

	return ctx.JSON(http.StatusOK, b.mapManyToProxy(list))
}

func (b BoardListener) Create(ctx echo.Context) error {
	userId, sessionStatus := b.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	var createRequest models.CreateRequest
	err := ctx.Bind(&createRequest)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Некорректный запрос")
	}

	linkSharedMode := (*estbackapi.LinkShareMode)(&createRequest.LinkSharedMode)
	if !linkSharedMode.IsValid() {
		return ctx.String(http.StatusBadRequest, "Некорректный link shared mode")
	}
	upsertBoardDto := estbackapi.UpsertBoardDto{
		Name: createRequest.Name,
		Description: &createRequest.Description,
		LinkSharedMode: linkSharedMode,
	}

	boardDto, _, err := b.boardApi.Create(
		ctx.Request().Context(), userId.String()).UpsertBoardDto(upsertBoardDto).Execute()
	if err != nil {
		return ctx.String(http.StatusConflict, "Не получилось создать доску")
	}

	return ctx.JSON(http.StatusOK, b.mapToProxy(boardDto))
}

func (b BoardListener) Update(ctx echo.Context) error {
	userId, sessionStatus := b.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	boardId := ctx.QueryParam("boardId")
	board, err := b.getBoard(ctx, boardId)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Доска не найдена")
	}
	if b.checkAccess(userId, board) != nil {
		return ctx.String(http.StatusForbidden, "Недостаточно прав")
	}

	var updateRequest models.CreateRequest
	err = ctx.Bind(&updateRequest)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Некорректный запрос")
	}

	linkSharedMode := (*estbackapi.LinkShareMode)(&updateRequest.LinkSharedMode)
	if !linkSharedMode.IsValid() {
		return ctx.String(http.StatusBadRequest, "Некорректный link shared mode")
	}
	upsertBoardDto := estbackapi.UpsertBoardDto{
		Name: updateRequest.Name,
		Description: &updateRequest.Description,
		LinkSharedMode: linkSharedMode,
	}

	boardDto, _, err := b.boardApi.Update(
		ctx.Request().Context(), boardId).UpsertBoardDto(upsertBoardDto).Execute()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Не получилось обновить доску")
	}

	return ctx.JSON(http.StatusOK, b.mapToProxy(boardDto))
}

func (b BoardListener) DeleteBoard(ctx echo.Context) error {
	userId, sessionStatus := b.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	boardId := ctx.QueryParam("boardId")
	board, err := b.getBoard(ctx, boardId)
	if err == nil && b.checkAccess(userId, board) != nil {
		return ctx.String(http.StatusForbidden, "Недостаточно прав")
	}
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Доска не найдена")
	}

	_, err = b.boardApi.DeleteBoard(ctx.Request().Context(), boardId).Execute()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Не получилось удалить доску")
	}

	return ctx.String(http.StatusOK, "Доска успешно удалена")
}

func (b BoardListener) Share(ctx echo.Context) error {
	ownerId, sessionStatus := b.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	boardId := ctx.QueryParam("boardId")
	board, err := b.getBoard(ctx, boardId)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Доска не найдена")
	}
	if board.OwnerInfo.Id != ownerId.String() {
		return ctx.String(http.StatusForbidden, "Недостаточно прав")
	}

	var shareBoardDto models.ShareBoardDto
	err = ctx.Bind(&shareBoardDto)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Некорректный запрос")
	}

	backSharingDto := estbackapi.BackSharingDto{
		UserId: shareBoardDto.UserId,
		Access: estbackapi.ShareMode(shareBoardDto.Access),
	}

	_, err = b.boardApi.Share(ctx.Request().Context(), 
							  boardId).BackSharingDto(backSharingDto).Execute()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Не получилось выдать права доступа на доску")
	}

	return ctx.String(http.StatusOK, "Права доступа на доску успешно выданы")
}

func (b BoardListener) ChangeAccess(ctx echo.Context) error {
	ownerId, sessionStatus := b.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	boardId := ctx.QueryParam("boardId")
	board, err := b.getBoard(ctx, boardId)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Доска не найдена")
	}
	if board.OwnerInfo.Id != ownerId.String() {
		return ctx.String(http.StatusForbidden, "Недостаточно прав")
	}

	var shareBoardDto models.ShareBoardDto
	err = ctx.Bind(&shareBoardDto)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Некорректный запрос")
	}

	backSharingDto := estbackapi.BackSharingDto{
		UserId: shareBoardDto.UserId,
		Access: estbackapi.ShareMode(shareBoardDto.Access),
	}

	_, err = b.boardApi.UpdateShare(ctx.Request().Context(), 
									boardId).BackSharingDto(backSharingDto).Execute()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Не получилось изменить права доступа на доску")
	}

	return ctx.String(http.StatusOK, "Права доступа на доску успешно изменены")
}

func (b BoardListener) Unshare(ctx echo.Context) error {
	ownerId, sessionStatus := b.getSession(ctx)
	if sessionStatus != nil {
		return sessionStatus
	}

	boardId := ctx.QueryParam("boardId")
	board, err := b.getBoard(ctx, boardId)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Доска не найдена")
	}
	if board.OwnerInfo.Id != ownerId.String() {
		return ctx.String(http.StatusForbidden, "Недостаточно прав")
	}

	var unshareRequest models.UnshareRequest
	err = ctx.Bind(&unshareRequest)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Некорректный запрос")
	}

	unshareBoardDto := estbackapi.UnshareBoardDto{
		UserId: unshareRequest.UserId,
	}

	_, err = b.boardApi.Unshare(ctx.Request().Context(), 
								boardId).UnshareBoardDto(unshareBoardDto).Execute()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Не получилось изменить права доступа на доску")
	}

	return ctx.String(http.StatusOK, "Права доступа на доску успешно изменены")
}

func (b BoardListener) Connect(ctx echo.Context) error {
	return ctx.String(http.StatusInternalServerError, "В разработке")
}

func (b BoardListener) mapManyToProxy(list *estbackapi.BackBoardListDto) models.BoardListDto {
	mine := make([]models.BoardDto, 0)
	shared := make([]models.BoardDto, 0)

	for _, dto := range list.Mine {
		mine = append(mine, b.mapToProxy(&dto))
	}

	for _, dto := range list.Shared {
		shared = append(shared, b.mapToProxy(&dto))
	}

	return models.BoardListDto{
		Mine:   mine,
		Shared: shared,
	}
}

func (b BoardListener) getMappedUserInfo(userIdStr string) models.UserDto {
	userInfo := models.UserDto{}
	userId, err := uuid.Parse(userIdStr)
	if err == nil {
		user, err := b.userRepository.GetUserByID(&userId)
		if err == nil {
			userInfo.Id = userIdStr
			userInfo.Username = user.Username
			userInfo.Avatar = user.Avatar
		}
	}
	return userInfo
}

func (b BoardListener) mapToProxy(dto *estbackapi.BackBoardDto) models.BoardDto {
	sharedWith := make([]models.SharingDto, 0)

	for _, shared_dto := range dto.SharedWith { 
		sharedWith = append(sharedWith, models.SharingDto{
			UserInfo: b.getMappedUserInfo(shared_dto.UserId),
			Access: string(shared_dto.Access),
		})
	}

	return models.BoardDto{
		Id: dto.Id,
		Name: dto.Name,
		Description: dto.Description,
		OwnerInfo: b.getMappedUserInfo(dto.OwnerId),
		SharedWith: sharedWith,
		LinkSharedMode: string(dto.LinkSharedMode),
		Preview: "TODO",
	}
}

func (b BoardListener) getSession(ctx echo.Context) (*uuid.UUID, error) {
	session, err := service.GetAndParseUserJWT(ctx)
	if err != nil {
		return nil, ctx.String(http.StatusUnauthorized, "Отсутствует или некорректный токен")
	}

	if session.ExpirationTime.After(time.Now()) {
		return nil, ctx.String(http.StatusUnauthorized, "Срок сессии истёк")
	}

	if _, err := b.userRepository.GetUserByID(&session.UserID); err != nil {
		return nil, ctx.String(http.StatusUnauthorized, "Пользователь не найден")
	}

	return &session.UserID, nil
}

func (b BoardListener) getBoard(ctx echo.Context, boardId string) (*models.BoardDto, error) {
	boardDto, _, err := b.boardApi.GetByUuid(ctx.Request().Context(), "").Execute()
	if err != nil {
		return nil, err
	}
	board := b.mapToProxy(boardDto)
	return &board, nil
}

func (b BoardListener) checkAccess(userId *uuid.UUID, boardDto *models.BoardDto) error {
	if userId == nil || boardDto == nil {
		return fmt.Errorf("no userId or boardDto given")
	}

	userIdStr := userId.String()
	accessed := userIdStr == boardDto.OwnerInfo.Id
	for _, sharedInfo := range boardDto.SharedWith {
		accessed = accessed || (userIdStr == sharedInfo.UserInfo.Id)
	}
	
	if !accessed {
		return fmt.Errorf("access denied")
	}

	return nil;
}