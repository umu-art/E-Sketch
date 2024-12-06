package mapper

import (
	"context"
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

func MapBackBoardToProxy(ctx context.Context, dto estbackapi.BackBoardDto, getUsersList func(context.Context, []string) []models.PublicUser) *proxymodels.BoardDto {
	userDtos := getMappedBackBoardUserIds(ctx, []estbackapi.BackBoardDto{dto}, getUsersList)
	return mapBackBoardDtoToProxy(dto, userDtos)
}

func MapManyBoardsToProxy(ctx context.Context, list *estbackapi.BackBoardListDto, getUsersList func(context.Context, []string) []models.PublicUser) *proxymodels.BoardListDto {
	mine := make([]proxymodels.BoardDto, 0)
	shared := make([]proxymodels.BoardDto, 0)

	boardDtos := make([]estbackapi.BackBoardDto, 0)
	boardDtos = append(boardDtos, list.Mine...)
	boardDtos = append(boardDtos, list.Shared...)

	userDtos := getMappedBackBoardUserIds(ctx, boardDtos, getUsersList)

	for _, dto := range list.Mine {
		mine = append(mine, *mapBackBoardDtoToProxy(dto, userDtos))
	}

	for _, dto := range list.Shared {
		shared = append(shared, *mapBackBoardDtoToProxy(dto, userDtos))
	}

	return &proxymodels.BoardListDto{
		Mine:   mine,
		Shared: shared,
	}
}

func getMappedBackBoardUserIds(ctx context.Context, boards []estbackapi.BackBoardDto, getUsersList func(context.Context, []string) []models.PublicUser) map[string]models.PublicUser {
	userIds := make([]string, 0)
	for _, board := range boards {
		userIds = append(userIds, board.OwnerId)
		for _, sharedDto := range board.SharedWith {
			userIds = append(userIds, sharedDto.UserId)
		}
	}

	users := getUsersList(ctx, userIds)

	mappedUsers := make(map[string]models.PublicUser)
	for _, user := range users {
		mappedUsers[user.ID.String()] = user
	}

	return mappedUsers
}

func mapBackBoardDtoToProxy(dto estbackapi.BackBoardDto, userDtos map[string]models.PublicUser) *proxymodels.BoardDto {
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
		Preview:        "TODO",
	}
}
