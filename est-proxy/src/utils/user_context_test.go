package utils

import (
	estbackapi "est_back_go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAccessLevel(t *testing.T) {
	ownerID := uuid.New()
	readID := uuid.New()
	writeID := uuid.New()
	adminID := uuid.New()
	noneID := uuid.New()

	board := &estbackapi.BackBoardDto{
		OwnerId: ownerID.String(),
		SharedWith: []estbackapi.BackSharingDto{
			{UserId: readID.String(), Access: estbackapi.READ},
			{UserId: writeID.String(), Access: estbackapi.WRITE},
			{UserId: adminID.String(), Access: estbackapi.ADMIN},
		},
	}

	tests := []struct {
		name  string
		user  *uuid.UUID
		board *estbackapi.BackBoardDto
		want  AccessLevel
	}{
		{"Нет пользователя", nil, board, NONE},
		{"Нет доски", &ownerID, nil, NONE}, // если передать board=nil
		{"owner", &ownerID, board, ADMIN},
		{"read", &readID, board, READ},
		{"write", &writeID, board, WRITE},
		{"admin", &adminID, board, ADMIN},
		{"no access", &noneID, board, NONE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAccessLevel(tt.user, tt.board)
			require.Equal(t, tt.want, got)
		})
	}
}
