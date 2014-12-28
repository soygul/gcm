package ccs

// OutMsg is a message to be sent to GCM CCS.
// If ID field is not set, it will be generated automatically using crypto/rand.
// Google recommends Data field to be strings key/value pairs and keys cannot be
// reserved words described in GCM server documentation.
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

// InMsg is an incoming GCM CCS message.
type InMsg struct {
	From        string            `json:"from"`
	ID          string            `json:"message_id"`
	Category    string            `json:"category"`
	Data        map[string]string `json:"data"`
	MessageType string            `json:"message_type"`
	ControlType string            `json:"control_type"`
	Err         string            `json:"error"`
	ErrDesc     string            `json:"error_description"`
}
