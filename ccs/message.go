package ccs

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

// OutMsg is an XMPP <message> stanzas used in sending messages to the GCM CCS server.
// https://developer.android.com/google/gcm/ccs.html#format
type OutMsg struct {
	To                       string            `json:"to"`
	ID                       string            `json:"message_id"`
	Data                     map[string]string `json:"data,omitempty"`
	MessageType              string            `json:"message_type,omitempty"`
	CollapseKey              string            `json:"collapse_key,omitempty"`
	TimeToLive               int               `json:"time_to_live,omitempty"`               //default:2419200 (in seconds = 4 weeks)
	DelayWhileIdle           bool              `json:"delay_while_idle,omitempty"`           //default:false
	DeliveryReceiptRequested bool              `json:"delivery_receipt_requested,omitempty"` //default:false
}

// InMsg is an XMPP <message> stanzas coming from the CCS server.
type InMsg struct {
	From        string            `json:"from"`
	ID          string            `json:"message_id"`
	Data        map[string]string `json:"data"`
	MessageType string            `json:"message_type"`
	ControlType string            `json:"control_type"`
	Err         string            `json:"error"`
	ErrDesc     string            `json:"error_description"`
}

// NewMsg creates a outgoing CCS message.
func NewMsg(id string) OutMsg {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return OutMsg{
		To:   id,
		ID:   "m-" + strconv.Itoa(r.Intn(100000)),
		Data: make(map[string]string),
	}
}

// SetData adds a key/value pair to the message payload data. Google recommends key/value pairs to be strings and
// keys cannot be reserved words described in GCM server documentation.
func (m *OutMsg) SetData(key string, value string) {
	if m.Data == nil {
		m.Data = make(map[string]string)
	}
	m.Data[key] = value
}

func (m *OutMsg) String() string {
	result, _ := json.Marshal(m)
	return string(result)
}
