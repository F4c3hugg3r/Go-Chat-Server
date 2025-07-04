package client

import (
	"io"
	"net/http"
)

type Reader interface {
	ReadString(delim byte) (string, error)
}

type Client struct {
	clientName string
	clientId   string
	reader     Reader
	writer     io.Writer
	authToken  string
	HttpClient *http.Client
}

// Message contains the name of the requester and the message (content) itsself
type Message struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	Plugin   string `json:"plugin"`
	ClientId string `json:"clientId"`
}

// Response contains the name of the sender and the response (content) itsself
type Response struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
