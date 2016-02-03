package ccs_test

import (
	"log"

	"github.com/titan-x/gcm/ccs"
)

// Example demonstrating the use of CCS implementation in an application server.
func Example() {
	c, err := ccs.Connect("gcm-preprod.googleapis.com:5236", "gcm_sender_id", "gcm_api_key", true)
	if err != nil {
		log.Fatalf("GCM CCS connection cannot be established")
	}

	// Send a test message. Replace "device_registration_id" with an actual GCM registration ID from a device.
	n, err := c.Send(&ccs.OutMsg{To: "device_registration_id", Data: map[string]string{"test_message": "GCM CCS client testing message."}})
	if err != nil {
		log.Printf("Failed to send message to CCS server with error: %v\n", err)
	}
	log.Printf("Message sent with %v bytes written to the connection\n", n)

	// Start receiving messages from the CCS server.
	for {
		log.Println("Waiting for incoming CCS messages")
		m, err := c.Receive()
		if err != nil {
			log.Printf("Incoming CCS error: %v\n", err)
		}

		go handleMessage(m)
	}
}

func handleMessage(m *ccs.InMsg) {
	log.Printf("Incoming CCS message: %v\n", m)
}
