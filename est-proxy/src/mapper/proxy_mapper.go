package mapper

import (
	"est-proxy/src/models"
	estbackapi "est_back_go"
	proxymodels "est_proxy_go/models"
	"log"
)

func MapUserToProxy(user models.PublicUser) *proxymodels.UserDto {
	return &proxymodels.UserDto{
		Id:       user.ID.String(),
		Username: user.Username,
		Avatar:   user.Avatar,
	}
}

func MapUserListToProxy(users []models.PublicUser) *[]proxymodels.UserDto {
	dtos := make([]proxymodels.UserDto, 0)
	for _, user := range users {
		dtos = append(dtos, *MapUserToProxy(user))
	}
	return &dtos
}

func MapBackBoardToProxy(
	dto estbackapi.BackBoardDto,
	getUsersList GetUsersFunc,
	getPreviewTokens GetPreviewTokensFunc) *proxymodels.BoardDto {

	userDtos := getMappedBackBoardUserIds([]estbackapi.BackBoardDto{dto}, getUsersList)
	return mapBackBoardDtoToProxy(dto, userDtos, getPreviewTokens([]string{dto.Id})[dto.Id])
}

func MapManyBoardsToProxy(
	list *estbackapi.BackBoardListDto,
	getUsersList GetUsersFunc,
	getPreviewTokens GetPreviewTokensFunc) *proxymodels.BoardListDto {

	mine := make([]proxymodels.BoardDto, 0)
	shared := make([]proxymodels.BoardDto, 0)
	recent := make([]proxymodels.BoardDto, 0)

	boardDtos := make([]estbackapi.BackBoardDto, 0)
	boardDtos = append(boardDtos, list.Mine...)
	boardDtos = append(boardDtos, list.Shared...)
	boardDtos = append(boardDtos, list.Recent...)

	userDtos := getMappedBackBoardUserIds(boardDtos, getUsersList)
	previewTokens := getMappedPreviewTokens(boardDtos, getPreviewTokens)

	for _, dto := range list.Mine {
		mine = append(mine, *mapBackBoardDtoToProxy(dto, userDtos, previewTokens[dto.Id]))
	}

	for _, dto := range list.Shared {
		shared = append(shared, *mapBackBoardDtoToProxy(dto, userDtos, previewTokens[dto.Id]))
	}

	for _, dto := range list.Recent {
		recent = append(recent, *mapBackBoardDtoToProxy(dto, userDtos, previewTokens[dto.Id]))
	}

	return &proxymodels.BoardListDto{
		Mine:   mine,
		Shared: shared,
		Recent: recent,
	}
}

func mapBackBoardDtoToProxy(
	dto estbackapi.BackBoardDto,
	userDtos map[string]models.PublicUser,
	previewToken string) *proxymodels.BoardDto {

	sharedWith := make([]proxymodels.SharingDto, 0)

	for _, sharedDto := range dto.SharedWith {
		_, exists := userDtos[sharedDto.UserId]
		if !exists {
			log.Printf("User %v does not exist", sharedDto.UserId)
			continue
		}
		sharedWith = append(sharedWith, proxymodels.SharingDto{
			UserInfo: *MapUserToProxy(userDtos[sharedDto.UserId]),
			Access:   string(sharedDto.Access),
		})
	}

	return &proxymodels.BoardDto{
		Id:             dto.Id,
		Name:           dto.Name,
		Description:    dto.Description,
		OwnerInfo:      *MapUserToProxy(userDtos[dto.OwnerId]),
		SharedWith:     sharedWith,
		LinkSharedMode: string(dto.LinkSharedMode),
		Preview:        previewToken,
	}
}

func getMappedBackBoardUserIds(
	boards []estbackapi.BackBoardDto,
	getUsersList GetUsersFunc) map[string]models.PublicUser {

	userIds := make([]string, 0)
	for _, board := range boards {
		userIds = append(userIds, board.OwnerId)
		for _, sharedDto := range board.SharedWith {
			userIds = append(userIds, sharedDto.UserId)
		}
	}

	users := getUsersList(userIds)

	mappedUsers := make(map[string]models.PublicUser)
	for _, user := range users {
		mappedUsers[user.ID.String()] = user
	}

	return mappedUsers
}

func getMappedPreviewTokens(boards []estbackapi.BackBoardDto, getPreviewTokens GetPreviewTokensFunc) map[string]string {
	boardIds := make([]string, 0)
	for _, board := range boards {
		boardIds = append(boardIds, board.Id)
	}
	return getPreviewTokens(boardIds)
}

type GetUsersFunc func([]string) []models.PublicUser
type GetPreviewTokensFunc func([]string) map[string]string
