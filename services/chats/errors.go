package chats

import "errors"

var (
	ErrChatType = errors.New("chat type is invalid")
)

var (
	ErrUserID = errors.New("user id is invalid")
	ErrChatID = errors.New("chat id is invalid")

	ErrUserIDNil = errors.New("user id can't be nil")
	ErrChatIDNil = errors.New("chat id can't be nil")
)
