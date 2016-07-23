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

	"github.com/mattn/go-xmpp"
)

const (
	gcmMessageStanza = `<message id=""><gcm xmlns="google:mobile:data">%v</gcm></message>`
	gcmDomain        = "gcm.googleapis.com"
)

// Conn is a GCM CCS connection.
type Conn struct {
	Host, SenderID string
	debug          bool
	xmppConn       *xmpp.Client
}

// Connect connects to GCM CCS server denoted by host (production or staging CCS endpoint URI) along with relevant credentials.
// Debug mode dumps all CSS communications to stdout.
func Connect(host, senderID, apiKey string, debug bool) (*Conn, error) {
	if !strings.Contains(senderID, gcmDomain) {
		senderID += "@" + gcmDomain
	}

	c, err := xmpp.NewClient(host, senderID, apiKey, debug)
	if err != nil {
		return nil, err
	}

	if debug {
		log.Printf("New CCS connection established with XMPP parameters: %+v\n", c)
	}

	return &Conn{
		Host:     host,
		SenderID: senderID,
		debug:    debug,
		xmppConn: c,
	}, nil
}

// Receive waits to receive the next incoming messages from the CCS connection.
func (c *Conn) Receive() (*InMsg, error) {
	stanza, err := c.xmppConn.Recv()
	if err != nil {
		return nil, err
	}

	if c.debug {
		log.Printf("Incoming raw CCS stanza: %+v\n", stanza)
	}

	chat, ok := stanza.(xmpp.Chat)
	if !ok {
		return nil, nil
	}

	if chat.Type == "error" {
		// todo: once go-xmpp can parse XMPP error messages, return error with XMPP error message (issue: https://github.com/soygul/gcm/issues/14)
		return nil, errors.New("CCS returned an XMPP error (can be a stanza or JSON error or anything else)")
	}

	var m InMsg
	if err = json.Unmarshal([]byte(chat.Other[0]), &m); err != nil { // todo: handle other fields of chat (remote/type/text/other[1,2,..])
		return nil, errors.New("unknow message from CCS")
	}

	switch m.MessageType {
	case "ack":
		return &m, nil // todo: mark message as sent
	case "nack":
		return &m, nil // todo: try and resend the message (after reconnect if problem is about connection draining)
	case "receipt":
		return &m, nil // todo: mark message as delivered and remove from the queue
	case "control":
		return &m, nil // todo: handle connection draining (and any other control message type?)
	case "":
		// acknowledge the incoming ordinary messages as per spec
		ack := &OutMsg{MessageType: "ack", To: m.From, ID: m.ID}
		if _, err = c.Send(ack); err != nil {
			return nil, fmt.Errorf("failed to send ack message to CCS with error: %v", err)
		}
		return &m, nil
	default:
		// unknown message types can be ignored, as per GCM specs
	}
	return &m, nil
}

// Send sends a message to GCM CCS server and returns the number of bytes written and any error encountered.
// If empty message ID is given, it's auto-generated and message object is modified with the generated ID.
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
