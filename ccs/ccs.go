// Package ccs provides GCM CCS (Cloud Connection Server) client implementation using XMPP.
// https://developer.android.com/google/gcm/ccs.html
package ccs

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/nbusy/go-xmpp"
)

const (
	gcmMessageStanza = `<message id=""><gcm xmlns="google:mobile:data">%v</gcm></message>`
	gcmDomain        = "gcm.googleapis.com"
)

// Conn is a GCM CCS connection.
type Conn struct {
	Host, SenderID string
	Debug          bool
	xmppConn       *xmpp.Client
}

// Connect connects to GCM CCS server denoted by host (production or staging CCS endpoint URI) along with relevant credentials.
// Debug mode dumps all CSS communications to stdout.
func Connect(host, senderID, apiKey string, debug bool) (*Conn, error) {
	if !strings.Contains(senderID, gcmDomain) {
		senderID += "@" + gcmDomain
	}

	c, err := xmpp.NewClient(host, senderID, apiKey, debug)
	if debug {
		if err == nil {
			log.Printf("New CCS connection established with XMPP parameters: %+v\n", c)
		} else {
			log.Printf("New CCS connection failed to establish with XMPP parameters: %+v and with error: %v\n", c, err)
		}
	}
	if err != nil {
		return nil, err
	}

	return &Conn{
		Host:     host,
		SenderID: senderID,
		Debug:    debug,
		xmppConn: c,
	}, nil
}

// Receive waits to receive the next incoming messages from the CCS connection.
func (c *Conn) Receive() (*InMsg, error) {
	stanza, err := c.xmppConn.Recv()
	if err != nil {
		return nil, err
	}

	if c.Debug {
		log.Printf("Incoming raw CCS stanza: %+v\n", stanza)
	}

	chat, ok := stanza.(xmpp.Chat)
	if !ok {
		return nil, nil
	}

	var m InMsg
	if err = json.Unmarshal([]byte(chat.Other[0]), &m); err != nil { // todo: handle other fields of chat (remote/type/text/other[1,2,..])
		return nil, errors.New("unknow message from CCS")
	}

	switch m.MessageType {
	case "ack":
		return nil, nil
	case "nack":
		errFormat := "From: %v, Message ID: %v, Error: %v, Error Description: %v"
		result := fmt.Sprintf(errFormat, m.From, m.ID, m.Err, m.ErrDesc)
		return nil, errors.New(result)
	case "receipt":
		return nil, nil
	case "control":
		return nil, nil
	case "":
		// acknowledge the incoming message as per specs
		if m.From != "" { // todo: what if From is empty? review specs
			ack := &OutMsg{MessageType: "ack", To: m.From, ID: m.ID}
			if _, err = c.Send(ack); err != nil {
				return nil, fmt.Errorf("Failed to send ack message to CCS. Error was: %v", err)
			}
			return &m, nil
		}
	default:
		// unknown message types are ignored as adviced by the specs
	}
	return &m, nil
}

// Send sends a message to GCM CCS server and returns the number of bytes written and any error encountered.
func (c *Conn) Send(m *OutMsg) (n int, err error) {
	if m.ID == "" {
		if m.ID, err = getMsgID(); err != nil {
			return 0, err
		}
	}

	mb, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}
	ms := string(mb)
	res := fmt.Sprintf(gcmMessageStanza, ms)
	return c.xmppConn.SendOrg(res)
}

// getID generates a unique message ID using crypto/rand in the form "m-96bitBase16"
func getMsgID() (string, error) {
	// todo: we can use sequential numbers optionally, just as the Android client does (1, 2, 3..) in upstream messages
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("m-%x", b), nil
}

// Close a CSS connection.
func (c *Conn) Close() error {
	return c.xmppConn.Close()
}
