package utils

import (
	estbackapi "est_back_go"
	"github.com/google/uuid"
)

type AccessLevel string

const (
	NONE  AccessLevel = "none"
	READ  AccessLevel = "read"
	WRITE AccessLevel = "write"
	ADMIN AccessLevel = "admin"
)

func GetAccessLevel(userId *uuid.UUID, boardDto *estbackapi.BackBoardDto) AccessLevel {
	if userId == nil || boardDto == nil {
		return NONE
	}

	userIdStr := userId.String()

	if userIdStr == boardDto.OwnerId {
		return ADMIN
	}

	for _, sharedInfo := range boardDto.SharedWith {
		if userIdStr == sharedInfo.UserId {
			switch sharedInfo.Access {
			case estbackapi.READ:
				return READ
			case estbackapi.WRITE:
				return WRITE
			case estbackapi.ADMIN:
				return ADMIN
			}
		}
	}

	return NONE
}
