package streams

import (
	"errors"
	"sync"
	notspb "wch/gen/nots/v1"
)

type Stream struct {
	userID string
	stream notspb.NotificationService_SubscribeServer
}

type StreamHub struct {
	streams map[string]Stream
	mu      sync.RWMutex
}

func NewStreamHub() *StreamHub {
	return &StreamHub{
		streams: make(map[string]Stream),
	}
}

func (h *StreamHub) Add(userID string, s notspb.NotificationService_SubscribeServer) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.streams[userID] = Stream{
		userID: userID,
		stream: s,
	}
}

func (h *StreamHub) Remove(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.streams, userID)
}

func (h *StreamHub) SendToUser(userID string, n *notspb.Notification) error {
	h.mu.RLock()
	s, ok := h.streams[userID]
	h.mu.RUnlock()

	if !ok {
		return errors.New("user not connected")
	}

	return s.stream.Send(n)
}
