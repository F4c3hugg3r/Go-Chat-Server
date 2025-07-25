package plugins

import (
	"encoding/json"
	"fmt"
	"strings"

	n "github.com/F4c3hugg3r/Go-Chat-Server/pkg/client/network"
	t "github.com/F4c3hugg3r/Go-Chat-Server/pkg/shared"
)

// CallPlugin lets you participate in a voice call
type CallPlugin struct {
	c       *n.ChatClient
	logChan chan t.Logg
}

func NewCallPlugin(chatClient *n.ChatClient, logChan chan t.Logg) *CallPlugin {
	return &CallPlugin{c: chatClient}
}

func (cp *CallPlugin) CheckScope() int {
	return RegisteredOnly
}

func (cp *CallPlugin) Execute(message *t.Message) (error, string) {
	rsp, err := cp.c.PostMessage(message, t.PostPlugin)

	// rsp in peers umwandeln
	var Ids []string
	dec := json.NewDecoder(strings.NewReader(rsp.Content))
	dec.Decode(&Ids)
	cp.logChan <- t.Logg{Text: fmt.Sprintf("Ids aus Response: %s", strings.Join(Ids, ", ")), Method: "NewCallPlugin.Execute()"}

	// für jede Peer Handle Signal aufrufen
	for _, clientId := range Ids {
		cp.c.HandleSignal(&t.Response{ClientId: clientId}, cp.logChan)
		cp.logChan <- t.Logg{Text: fmt.Sprintf("Handle Signal started für client %s", clientId), Method: "NewCallPlugin.Execute()"}
	}

	return err, ""
}

// GroupPlugin lets you participate in a group chat
type GroupPlugin struct {
	c *n.ChatClient
}

func NewGroupPlugin(chatClient *n.ChatClient) *GroupPlugin {
	return &GroupPlugin{c: chatClient}
}

func (gp *GroupPlugin) CheckScope() int {
	return RegisteredOnly
}

func (gp *GroupPlugin) Execute(message *t.Message) (error, string) {
	_, err := gp.c.PostMessage(message, t.PostPlugin)
	return err, ""
}

// PrivateMessage Plugin lets a client send a private message to another client identified by it's clientId
type PrivateMessagePlugin struct {
	c *n.ChatClient
}

func NewPrivateMessagePlugin(chatClient *n.ChatClient) *PrivateMessagePlugin {
	return &PrivateMessagePlugin{c: chatClient}
}

func (pp *PrivateMessagePlugin) CheckScope() int {
	return RegisteredOnly
}

func (pp *PrivateMessagePlugin) Execute(message *t.Message) (error, string) {
	if message.Content == "" {
		return fmt.Errorf("%w: prefix shouldn't be empty", t.ErrParsing), ""
	}
	opposingClientId := strings.Fields(message.Content)[0]

	content, ok := strings.CutPrefix(message.Content, fmt.Sprintf("%s ", opposingClientId))
	if !ok {
		return fmt.Errorf("%w: prefix '%s ' not found", t.ErrParsing, opposingClientId), ""
	}

	_, err := pp.c.PostMessage(pp.c.CreateMessage(message.Name, message.Plugin, content, opposingClientId), t.PostPlugin)

	return err, ""
}

// LogOutPlugin logs out a client by deleting it out of the clients map
type LogOutPlugin struct {
	c *n.ChatClient
}

func NewLogOutPlugin(chatClient *n.ChatClient) *LogOutPlugin {
	return &LogOutPlugin{c: chatClient}
}

func (lp *LogOutPlugin) CheckScope() int {
	return RegisteredOnly
}

func (lp *LogOutPlugin) Execute(message *t.Message) (error, string) {
	return lp.c.PostDelete(message), t.UnregisterFlag
}

// RegisterClientPlugin safely registeres a client by creating a Client with the received values
// and putting it into the global clients map
type RegisterClientPlugin struct {
	c *n.ChatClient
}

func NewRegisterClientPlugin(chatClient *n.ChatClient) *RegisterClientPlugin {
	return &RegisterClientPlugin{c: chatClient}
}

func (rp *RegisterClientPlugin) CheckScope() int {
	return UnregisteredOnly
}

func (rp *RegisterClientPlugin) Execute(message *t.Message) (error, string) {
	clientName := message.Content
	if len(clientName) > 50 || len(clientName) < 3 {
		return fmt.Errorf("%w: your name has to be between 3 and 50 chars long", t.ErrParsing), ""
	}

	rsp, err := rp.c.PostMessage(rp.c.CreateMessage(clientName, message.Plugin, message.Content, message.ClientId), t.PostRegister)
	if err != nil {
		return fmt.Errorf("%w: error sending message", err), ""
	}

	err = rp.c.Register(rsp)
	if err != nil {
		return fmt.Errorf("%w: error registering client", err), ""
	}

	return err, t.RegisterFlag
}

// BroadcaastPlugin distributes an incomming message abroad all client channels if
// a client can't receive, i'ts active status is set to false
type BroadcastPlugin struct {
	c *n.ChatClient
}

func NewBroadcastPlugin(chatClient *n.ChatClient) *BroadcastPlugin {
	return &BroadcastPlugin{c: chatClient}
}

func (bp *BroadcastPlugin) CheckScope() int {
	return RegisteredOnly
}

func (bp *BroadcastPlugin) Execute(message *t.Message) (error, string) {
	_, err := bp.c.PostMessage(message, t.PostPlugin)
	return err, ""
}

// HelpPlugin tells you information about available plugins
type HelpPlugin struct {
	c *n.ChatClient
}

func NewHelpPlugin(chatClient *n.ChatClient) *HelpPlugin {
	return &HelpPlugin{c: chatClient}
}

func (hp *HelpPlugin) CheckScope() int {
	return RegisteredOnly
}

func (h *HelpPlugin) Execute(message *t.Message) (error, string) {
	_, err := h.c.PostMessage(message, t.PostPlugin)
	return err, ""
}

// UserPlugin tells you information about all the current users
type UserPlugin struct {
	c *n.ChatClient
}

func NewUserPlugin(chatClient *n.ChatClient) *UserPlugin {
	return &UserPlugin{c: chatClient}
}

func (up *UserPlugin) CheckScope() int {
	return RegisteredOnly
}

func (u *UserPlugin) Execute(message *t.Message) (error, string) {
	_, err := u.c.PostMessage(message, t.PostPlugin)
	return err, ""
}

// TimePlugin tells you the current time
type TimePlugin struct {
	c *n.ChatClient
}

func NewTimePlugin(chatClient *n.ChatClient) *TimePlugin {
	return &TimePlugin{c: chatClient}
}

func (tp *TimePlugin) CheckScope() int {
	return RegisteredOnly
}

func (tp *TimePlugin) Execute(message *t.Message) (error, string) {
	_, err := tp.c.PostMessage(message, t.PostPlugin)
	return err, ""
}
