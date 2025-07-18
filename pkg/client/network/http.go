package network

import (
	"bytes"
	"fmt"
	"net/http"
)

// GetRequest sends a GET Request to the server including the authorization token
func (c *ChatClient) GetRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: Fehler beim erstellen der GET request", err)
	}

	authToken, ok := c.GetAuthToken()
	if !ok && c.Registered {
		return nil, fmt.Errorf("%w: client not registered anymore", err)
	}

	req.Header.Add("Authorization", authToken)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: Fehler beim senden der GET request", err)
	}

	return res, nil
}

// DeleteRequest sends a DELETE Request to delete the client out of the server
// including the authorization token
func (c *ChatClient) DeleteRequest(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("%w: Fehler beim Erstellen der DELETE req", err)
	}

	authToken, ok := c.GetAuthToken()
	if !ok {
		return nil, fmt.Errorf("%w: client not registered anymore", err)
	}

	req.Header.Add("Authorization", authToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: Fehler beim Absenden des Deletes", err)
	}

	return res, nil
}

// PostReqeust sends a Post Request to send a message to the server
// including the authorization token
func (c *ChatClient) PostRequest(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("%w: Fehler beim Erstellen der POST req", err)
	}

	authToken, ok := c.GetAuthToken()
	if !ok && c.Registered {
		return nil, fmt.Errorf("%w: client not registered anymore", err)
	}

	req.Header.Add("Authorization", authToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: Fehler beim Absenden der Nachricht", err)
	}

	return res, nil
}
