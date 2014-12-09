package ccs

import "log"

func Example() {
	c, err := Connect("gcm-preprod.googleapis.com:5236", "gcm_sender_id", "gcm_api_key", true)
	if err != nil {
		log.Fatalf("GCM CCS connection cannot be established.")
	}

	for {
		log.Printf("Waiting for incoming CCS messages")
		m, err := c.Receive()
		if err != nil {
			log.Printf("Incoming CCS error: %v\n", err)
		}

		go func(m *InMsg) {
			log.Printf("Incoming CCS message: %v\n", m)
		}(m)
	}
}
