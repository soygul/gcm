package ccs_test

import (
	"log"

	"github.com/nbusy/gcm/ccs"
)

// Package level example demonstrating the use of CCS implementation in an application server.
func Example() {
	c, err := ccs.Connect("gcm-preprod.googleapis.com:5236", "gcm_sender_id", "gcm_api_key", true)
	if err != nil {
		log.Fatalf("GCM CCS connection cannot be established.")
	}

	err = c.Send(&ccs.OutMsg{To: "device_registration_id", Data: map[string]string{"test_message": "GCM CCS client testing message."}})
	if err != nil {
		log.Fatalf("Failed to send message to CCS server with error: %v\n", err)
	}

	for {
		log.Printf("Waiting for incoming CCS messages")
		m, err := c.Receive()
		if err != nil {
			log.Printf("Incoming CCS error: %v\n", err)
		}

		go readHandler(m)
	}
}

func readHandler(m *ccs.InMsg) {
	log.Printf("Incoming CCS message: %v\n", m)
}
