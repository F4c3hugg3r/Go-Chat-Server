package server

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/F4c3hugg3r/Go-Chat-Server/pkg/shared"
)

// PrivateMessage Plugin lets a client send a private message to another client identified by it's clientId
type PrivateMessagePlugin struct {
	chatService *ChatService
}

func NewPrivateMessagePlugin(s *ChatService) *PrivateMessagePlugin {
	return &PrivateMessagePlugin{chatService: s}
}

func (pp *PrivateMessagePlugin) Description() string {
	return "'/private {Id} {message}'"
}

func (pp *PrivateMessagePlugin) Execute(message *Message) (*Response, error) {
	pp.chatService.mu.RLock()
	defer pp.chatService.mu.RUnlock()

	client, ok := pp.chatService.clients[message.ClientId]
	if !ok {
		return nil, fmt.Errorf("%w: client with id: %s not found", ErrClientNotAvailable, message.ClientId)
	}

	rsp := &Response{fmt.Sprintf("[Private] - %s", message.Name), message.Content}

	err := client.Send(rsp)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

// LogOutPlugin logs out a client by deleting it out of the clients map
type LogOutPlugin struct {
	chatService *ChatService
	pr          *PluginRegistry
}

func NewLogOutPlugin(s *ChatService, pr *PluginRegistry) *LogOutPlugin {
	return &LogOutPlugin{
		chatService: s,
		pr:          pr,
	}
}

func (lp *LogOutPlugin) Description() string {
	return "'/quit'"
}

func (lp *LogOutPlugin) Execute(message *Message) (*Response, error) {
	lp.chatService.mu.Lock()
	defer lp.chatService.mu.Unlock()

	client, ok := lp.chatService.clients[message.ClientId]
	if !ok {
		return nil, fmt.Errorf("%w: client (probably) already deleted", ErrClientNotAvailable)
	}

	fmt.Println("\nlogged out ", client.Name)
	client.Close()
	delete(lp.chatService.clients, message.ClientId)

	go client.Execute(lp.pr, &Message{Name: "", Plugin: "/broadcast", Content: fmt.Sprintf("%s hat den Chat verlassen", message.Name), ClientId: message.ClientId})

	return &Response{Name: message.Name, Content: "logged out"}, nil
}

// RegisterClientPlugin safely registeres a client by creating a Client with the received values
// and putting it into the global clients map
type RegisterClientPlugin struct {
	chatService *ChatService
	pr          *PluginRegistry
}

func NewRegisterClientPlugin(s *ChatService, pr *PluginRegistry) *RegisterClientPlugin {
	return &RegisterClientPlugin{
		chatService: s,
		pr:          pr,
	}
}

func (rp *RegisterClientPlugin) Description() string {
	return "'/register {name}'"
}

func (rp *RegisterClientPlugin) Execute(message *Message) (*Response, error) {
	rp.chatService.mu.Lock()
	defer rp.chatService.mu.Unlock()

	if len(rp.chatService.clients) >= rp.chatService.maxUsers {
		return nil,
			fmt.Errorf("%w: usercap %d reached, try again later. users:%d", ErrNoPermission, rp.chatService.maxUsers, len(rp.chatService.clients))
	}

	if _, exists := rp.chatService.clients[message.ClientId]; exists {
		return nil, fmt.Errorf("%w: client already defined", ErrNoPermission)
	}

	clientCh := make(chan *Response, 100)
	token := shared.GenerateSecureToken(64)
	client := &Client{
		Name:      message.Name,
		ClientId:  message.ClientId,
		clientCh:  clientCh,
		Active:    true,
		authToken: token,
		lastSign:  time.Now().UTC(),
		chClosed:  false,
	}
	rp.chatService.clients[message.ClientId] = client

	fmt.Printf("\nnew client '%s' registered.", message.Content)

	go client.Execute(rp.pr, &Message{Name: "", Plugin: "/broadcast", Content: fmt.Sprintf("%s ist dem Chat beigetreten", client.Name), ClientId: client.ClientId})

	return &Response{Name: message.Name, Content: token}, nil
}

// BroadcaastPlugin distributes an incomming message abroad all client channels if
// a client can't receive, i'ts active status is set to false
type BroadcastPlugin struct {
	chatService *ChatService
}

func NewBroadcastPlugin(s *ChatService) *BroadcastPlugin {
	return &BroadcastPlugin{chatService: s}
}

func (bp *BroadcastPlugin) Description() string {
	return "'{message}' or '/broadcast {message}"
}

func (bp *BroadcastPlugin) Execute(message *Message) (*Response, error) {
	bp.chatService.mu.RLock()
	defer bp.chatService.mu.RUnlock()

	rsp := &Response{Name: message.Name, Content: message.Content}

	if strings.TrimSpace(message.Content) == "" {
		return rsp, nil
	}

	if len(bp.chatService.clients) <= 0 {
		return nil, fmt.Errorf("%w: There are no clients registered", ErrClientNotAvailable)
	}

	for _, client := range bp.chatService.clients {
		if client.ClientId != message.ClientId {
			err := client.Send(rsp)
			if err != nil {
				log.Printf("\n%v: %s -> %s", err, message.Name, client.Name)
			}
		}
	}

	return rsp, nil
}

// HelpPlugin tells you information about available plugins
type HelpPlugin struct {
	pr *PluginRegistry
}

func NewHelpPlugin(pr *PluginRegistry) *HelpPlugin {
	return &HelpPlugin{pr: pr}
}

func (h *HelpPlugin) Description() string {
	return "'/help'"
}

func (h *HelpPlugin) Execute(message *Message) (*Response, error) {
	jsonList, err := json.Marshal(h.pr.ListPlugins())
	if err != nil {
		return nil, fmt.Errorf("%w: error parsing plugins to json", err)
	}

	return &Response{"Help", string(jsonList)}, nil
}

// UserPlugin tells you information about all the current users
type UserPlugin struct {
	chatService *ChatService
}

func NewUserPlugin(s *ChatService) *UserPlugin {
	return &UserPlugin{chatService: s}
}

func (u *UserPlugin) Description() string {
	return "'/users'"
}

func (u *UserPlugin) Execute(message *Message) (*Response, error) {
	u.chatService.mu.RLock()
	defer u.chatService.mu.RUnlock()

	clientsSlice := []json.RawMessage{}

	for _, client := range u.chatService.clients {
		jsonString, err := json.Marshal(client)
		if err != nil {
			log.Printf("error parsing client %s to json", client.Name)
		}

		clientsSlice = append(clientsSlice, jsonString)
	}

	jsonList, err := json.Marshal(clientsSlice)
	if err != nil {
		return nil, fmt.Errorf("%w: error parsing clients to json", err)
	}

	return &Response{"Users", string(jsonList)}, nil
}

// TimePlugin tells you the current time
type TimePlugin struct{}

func NewTimePlugin() *TimePlugin {
	return &TimePlugin{}
}

func (t *TimePlugin) Description() string {
	return "'/time'"
}

func (t *TimePlugin) Execute(message *Message) (*Response, error) {
	return &Response{"Time", time.Now().UTC().String()}, nil
}
