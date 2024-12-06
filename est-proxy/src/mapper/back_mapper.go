package mapper

import (
	estbackapi "est_back_go"
	proxymodels "est_proxy_go/models"
	"log"
)

func MapCreateRequestToBack(createRequest *proxymodels.CreateRequest) *estbackapi.UpsertBoardDto {
	linkSharedMode := (*estbackapi.LinkShareMode)(&createRequest.LinkSharedMode)
	if !linkSharedMode.IsValid() {
		log.Println("Invalid link shared mode")
		return nil
	}

	return &estbackapi.UpsertBoardDto{
		Name:           createRequest.Name,
		Description:    &createRequest.Description,
		LinkSharedMode: linkSharedMode,
	}
}

func MapShareBoardDtoToBack(shareBoardDto *proxymodels.ShareBoardDto) *estbackapi.BackSharingDto {
	return &estbackapi.BackSharingDto{
		UserId: shareBoardDto.UserId,
		Access: estbackapi.ShareMode(shareBoardDto.Access),
	}
}

func MapUnshareRequestToBack(unshareRequest *proxymodels.UnshareRequest) *estbackapi.UnshareBoardDto {
	return &estbackapi.UnshareBoardDto{
		UserId: unshareRequest.UserId,
	}
}
