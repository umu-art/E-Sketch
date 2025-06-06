package impl

import (
	"context"
	"est-proxy/src/api"
	"est-proxy/src/errors"
	"est-proxy/src/mapper"
	"est-proxy/src/models"
	"est-proxy/src/repository"
	"est-proxy/src/utils"
	estbackapi "est_back_go"
	proxymodels "est_proxy_go/models"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type BoardServiceImpl struct {
	boardApi       *estbackapi.BoardAPIService
	previewApi     api.PreviewApi
	userRepository repository.UserRepository
}

func NewBoardServiceImpl(boardApi *estbackapi.BoardAPIService, previewApi api.PreviewApi, userRepository repository.UserRepository) *BoardServiceImpl {
	return &BoardServiceImpl{
		boardApi:       boardApi,
		previewApi:     previewApi,
		userRepository: userRepository,
	}
}

func (bs *BoardServiceImpl) GetByUuid(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID) (*proxymodels.BoardDto, *errors.StatusError) {
	if userId == nil || boardId == nil {
		return nil, errors.NewStatusError(http.StatusNotFound, "Некорреткный запрос")
	}

	board, err := bs.getBoard(ctx, boardId.String())
	if err != nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Доска не найдена")
	}

	if utils.GetAccessLevel(userId, board) == utils.NONE {
		return nil, errors.NewStatusError(http.StatusForbidden, "Недостаточно прав")
	}

	bs.markRecent(ctx, userId.String(), boardId.String())

	return mapper.MapBackBoardToProxy(*board, bs.getUsersFunc(ctx), bs.getPreviewTokensFunc(ctx)), nil
}

func (bs *BoardServiceImpl) List(ctx context.Context, userId *uuid.UUID) (*proxymodels.BoardListDto, *errors.StatusError) {
	list, _, err := bs.boardApi.ListByUserId(ctx, userId.String()).Execute()
	if err != nil {
		log.Printf("\nwtf %v \n\n", err)
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Не получилось составить список досок")
	}

	return mapper.MapManyBoardsToProxy(list, bs.getUsersFunc(ctx), bs.getPreviewTokensFunc(ctx)), nil
}

func (bs *BoardServiceImpl) Create(ctx context.Context, userId *uuid.UUID, createRequest *proxymodels.CreateRequest) (*proxymodels.BoardDto, *errors.StatusError) {
	upsertBoardDto := mapper.MapCreateRequestToBack(createRequest)
	if upsertBoardDto == nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Неверный запрос")
	}

	boardDto, _, err := bs.boardApi.Create(
		ctx, userId.String()).UpsertBoardDto(*upsertBoardDto).Execute()

	if err != nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Не получилось создать доску")
	}

	bs.markRecent(ctx, userId.String(), boardDto.Id)

	return mapper.MapBackBoardToProxy(*boardDto, bs.getUsersFunc(ctx), bs.getPreviewTokensFunc(ctx)), nil
}

func (bs *BoardServiceImpl) Update(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID, updateRequest *proxymodels.CreateRequest) (*proxymodels.BoardDto, *errors.StatusError) {
	board, err := bs.getBoard(ctx, boardId.String())
	if err != nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Доска не найдена")
	}

	if utils.GetAccessLevel(userId, board) != utils.ADMIN {
		return nil, errors.NewStatusError(http.StatusForbidden, "Недостаточно прав")
	}

	upsertBoardDto := mapper.MapCreateRequestToBack(updateRequest)
	if upsertBoardDto == nil {
		return nil, errors.NewStatusError(http.StatusBadRequest, "Неверный запрос")
	}

	boardDto, _, err := bs.boardApi.Update(
		ctx, boardId.String()).UpsertBoardDto(*upsertBoardDto).Execute()

	if err != nil {
		return nil, errors.NewStatusError(http.StatusInternalServerError, "Не получилось обновить доску")
	}

	bs.markRecent(ctx, userId.String(), boardId.String())

	return mapper.MapBackBoardToProxy(*boardDto, bs.getUsersFunc(ctx), bs.getPreviewTokensFunc(ctx)), nil
}

func (bs *BoardServiceImpl) DeleteBoard(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID) *errors.StatusError {
	board, err := bs.getBoard(ctx, boardId.String())
	if err != nil {
		return errors.NewStatusError(http.StatusBadRequest, "Доска не найдена")
	}

	if utils.GetAccessLevel(userId, board) != utils.ADMIN {
		return errors.NewStatusError(http.StatusForbidden, "Недостаточно прав")
	}

	_, err = bs.boardApi.DeleteBoard(ctx, boardId.String()).Execute()
	if err != nil {
		return errors.NewStatusError(http.StatusInternalServerError, "Не получилось удалить доску")
	}

	return nil
}

func (bs *BoardServiceImpl) Share(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID, shareBoardDto *proxymodels.ShareBoardDto) *errors.StatusError {
	board, err := bs.getBoard(ctx, boardId.String())
	if err != nil {
		return errors.NewStatusError(http.StatusBadRequest, "Доска не найдена")
	}

	if board.OwnerId != userId.String() {
		return errors.NewStatusError(http.StatusForbidden, "Недостаточно прав")
	}

	backSharingDto := mapper.MapShareBoardDtoToBack(shareBoardDto)
	if backSharingDto == nil {
		return errors.NewStatusError(http.StatusBadRequest, "Неверный запрос")
	}

	_, err = bs.boardApi.
		Share(ctx, boardId.String()).
		BackSharingDto(*backSharingDto).
		Execute()

	if err != nil {
		return errors.NewStatusError(http.StatusInternalServerError, "Не получилось выдать права доступа на доску")
	}
	return nil
}

func (bs *BoardServiceImpl) ChangeAccess(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID, shareBoardDto *proxymodels.ShareBoardDto) *errors.StatusError {
	board, err := bs.getBoard(ctx, boardId.String())
	if err != nil {
		return errors.NewStatusError(http.StatusBadRequest, "Доска не найдена")
	}

	if board.OwnerId != userId.String() {
		return errors.NewStatusError(http.StatusForbidden, "Недостаточно прав")
	}

	backSharingDto := mapper.MapShareBoardDtoToBack(shareBoardDto)
	if backSharingDto == nil {
		return errors.NewStatusError(http.StatusBadRequest, "Неверный запрос")
	}

	_, err = bs.boardApi.UpdateShare(ctx,
		boardId.String()).BackSharingDto(*backSharingDto).Execute()
	if err != nil {
		return errors.NewStatusError(http.StatusInternalServerError, "Не получилось изменить права доступа на доску")
	}

	return nil
}

func (bs *BoardServiceImpl) Unshare(ctx context.Context, userId *uuid.UUID, boardId *uuid.UUID, unshareRequest *proxymodels.UnshareRequest) *errors.StatusError {
	board, err := bs.getBoard(ctx, boardId.String())
	if err != nil {
		return errors.NewStatusError(http.StatusBadRequest, "Доска не найдена")
	}
	if board.OwnerId != userId.String() {
		return errors.NewStatusError(http.StatusForbidden, "Недостаточно прав")
	}

	unshareBoardDto := mapper.MapUnshareRequestToBack(unshareRequest)
	if unshareBoardDto == nil {
		return errors.NewStatusError(http.StatusBadRequest, "Неверный запрос")
	}

	_, err = bs.boardApi.Unshare(ctx,
		boardId.String()).UnshareBoardDto(*unshareBoardDto).Execute()
	if err != nil {
		return errors.NewStatusError(http.StatusInternalServerError, "Не получилось изменить права доступа на доску")
	}

	return nil
}

func (bs *BoardServiceImpl) getBoard(ctx context.Context, boardId string) (*estbackapi.BackBoardDto, error) {
	boardDto, _, err := bs.boardApi.GetByUuid(ctx, boardId).Execute()
	if err != nil {
		return nil, err
	}
	return boardDto, nil
}

func (bs *BoardServiceImpl) getUsersFunc(ctx context.Context) func([]string) []models.PublicUser {
	return func(userIdStrs []string) []models.PublicUser {
		userIds := make([]uuid.UUID, len(userIdStrs))
		for i, userIdStr := range userIdStrs {
			userId, err := uuid.Parse(userIdStr)
			if err != nil {
				log.Printf("Failed to parse user id from userIdStr \"%s\": %v", userIdStr, err)
			}
			userIds[i] = userId
		}
		return *bs.userRepository.GetUserListByIds(ctx, userIds)
	}
}

func (bs *BoardServiceImpl) getPreviewTokensFunc(ctx context.Context) func(boardIds []string) map[string]string {
	return func(boardIds []string) map[string]string {
		res, err := bs.previewApi.GetTokens(boardIds, ctx)
		if err != nil {
			log.Printf("Failed to get token: %s", err.GetMessage())
			return nil
		}
		return res
	}
}

func (bs *BoardServiceImpl) markRecent(ctx context.Context, userId string, boardId string) {
	_, err := bs.boardApi.MarkAsRecent(ctx, userId).BoardIdDto(estbackapi.BoardIdDto{
		Id: boardId,
	}).Execute()
	if err != nil {
		log.Printf("Failed to mark board as recent: %v", err.Error())
	}
}
