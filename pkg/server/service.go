package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// clients who communicate with the sever
type ChatService struct {
	clients  map[string]*Client
	maxUsers int
	mu       sync.RWMutex
}

func NewChatService(maxUsers int) *ChatService {
	return &ChatService{
		clients:  make(map[string]*Client),
		maxUsers: maxUsers,
	}
}

// InactiveClientDeleter searches for inactive clients and deletes them as well as closes their message-channel
func (s *ChatService) InactiveClientDeleter(timeLimit time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for clientId, client := range s.clients {
		if client.Idle(timeLimit) {
			delete(s.clients, clientId)
		}
	}
}

// echo sends a response to the request submitter
func (s *ChatService) echo(clientId string, rsp *Response) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, ok := s.clients[clientId]
	if !ok {
		return fmt.Errorf("%w: message couldn't be echoed", ErrClientNotAvailable)
	}

	return client.Send(rsp)
}

// GetClientChannel tests if there is a registered client to the given clientId and returns
// it's channel and name
func (s *ChatService) GetClient(clientId string) (*Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	client, exists := s.clients[clientId]
	if !exists {
		return client, fmt.Errorf("%w: there is no client with id: %s registered", ErrClientNotAvailable, clientId)
	}

	return client, nil
}

// DecodeToMessage decodes a responseBody to a Message struct
func DecodeToMessage(body []byte) (Message, error) {
	message := Message{}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	err := dec.Decode(&message)

	if err != nil {
		return message, err
	}

	return message, nil
}
