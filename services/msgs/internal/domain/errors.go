package domain

import "errors"

var (
	ErrMessageID        = errors.New("message id is invalid")
	ErrMessageIDIsEmpty = errors.New("message id is empty")

	ErrChatID      = errors.New("chat id is invalid")
	ErrChatIsEmpty = errors.New("chat id is empty")

	ErrSenderID        = errors.New("sender id is invalid")
	ErrSenderIDIsEmpty = errors.New("sender id is empty")

	ErrNotAllowedToSend = errors.New("user is not allowed to send messages")
)
