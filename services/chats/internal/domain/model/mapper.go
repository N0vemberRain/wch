package model

import (
	chatspb "wch/gen/chats/v1"
	chats "wch/services/chats/internal/domain"

	"github.com/google/uuid"
	timeconv "google.golang.org/protobuf/types/known/timestamppb"
)

func ChatTypeToString(t ChatType) string {
	switch t {
	case ChatTypeDirect:
		return "direct"
	case ChatTypeGroup:
		return "group"
	default:
		return "unknown"
	}
}

func ChatTypeFromString(t string) (ChatType, error) {
	switch t {
	case "direct":
		return ChatTypeDirect, nil
	case "group":
		return ChatTypeGroup, nil
	case "unknown":
		return ChatTypeUnknown, nil
	default:
		return ChatTypeUnknown, chats.ErrChatType
	}
}

func ChatToProto(c *Chat) *chatspb.Chat {
	return &chatspb.Chat{
		Id:        c.ID.String(),
		Type:      ChatTypeToString(c.Type),
		Name:      c.Name,
		CreatedAt: timeconv.New(c.CreatedAt),
		UpdatedAt: timeconv.New(c.UpdatedAt),
	}
}

func ChatFromProto(c *chatspb.Chat) (*Chat, error) {
	t, err := ChatTypeFromString(c.Type)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(c.Id)
	if err != nil {
		return nil, err
	}
	return &Chat{
		ID:        id,
		Type:      t,
		Name:      c.Name,
		CreatedAt: c.CreatedAt.AsTime(),
		UpdatedAt: c.UpdatedAt.AsTime(),
	}, nil
}

func IDsToString(ids []uuid.UUID) []string {
	ids_str := make([]string, len(ids))
	for _, id := range ids {
		ids_str = append(ids_str, id.String())
	}

	return ids_str
}

func UserToProto(u *User) *chatspb.ChatUser {
	return &chatspb.ChatUser{
		Id:          u.ID.String(),
		DisplayName: u.Username,
		AvatarUrl:   u.AvatarURL,
	}
}

func ChatParticipantRoleToString(role ChatParticipantRole) string {
	switch role {
	case ChatParticipantAdmin:
		return "Admin"
	case ChatParticipantMember:
		return "Member"
	case ChatParticipantOwner:
		return "Admin"

	default:
		return "Unknown"
	}
}

func ChatParticipantRoleFromString(role string) ChatParticipantRole {
	switch role {
	case "Admin":
		return ChatParticipantAdmin
	case "Member":
		return ChatParticipantMember
	case "Owner":
		return ChatParticipantOwner
	default:
		return ChatParticipantUnknown
	}
}
