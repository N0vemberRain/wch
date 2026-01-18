package model

type ChatParticipantRole int32

const (
	ChatParticipantAdmin   = 1
	ChatParticipantMember  = 2
	ChatParticipantOwner   = 3
	ChatParticipantUnknown = 0
)
