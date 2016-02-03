# GCM CCS (XMPP)

[![Build Status](https://travis-ci.org/titan-x/gcm.svg?branch=master)](https://travis-ci.org/titan-x/gcm)
[![GoDoc](https://godoc.org/github.com/titan-x/gcm?status.svg)](https://godoc.org/github.com/titan-x/gcm)

GCM (Google Cloud Messaging) CCS (Cloud Connection Server) implementation for application servers as described in [Android developer docs](https://developer.android.com/google/gcm/ccs.html).

Uses the XMPP endpoint to have persistent and asynchronous connection with the Google's GCM servers.

The HTTP implementation is a work in progress but it is advisable to use the CCS (XMPP) implementation instead as it is asynchronous and hence utilizes server resources more efficiently.

## CCS Example

```go
package main

import (
	"fmt"
	"github.com/titan-x/gcm/ccs"
)

func main() {
	c, err := ccs.Connect("gcm-preprod.googleapis.com:5236", "gcm_sender_id", "gcm_api_key", true)
	if err != nil {
		return
	}

	// send a test message to a device
	_, err = c.Send(&ccs.OutMsg{To: "device_registration_id", Data: map[string]string{"test_message": "GCM CCS client testing message."}})

	// start receiving messages from CCS
	for {
		m, err := c.Receive()
		go func(m *ccs.InMsg) {
			log.Printf("message: %v\n error (if any): %v\n", m, err)
		}(m)
	}
}
```

To see a more comprehensive example, check the godocs.

## Testing

All the tests can be executed by `go test -race -cover ./...` command. Optionally you can add `-v` flag to observe all connection logs. Integration tests require the following environment variables to be defined. If they are missing, integration tests are skipped.

```bash
export GCM_CCS_HOST=gcm-preprod.googleapis.com:5236
export GCM_SENDER_ID=preprod_sender_id
export GOOGLE_API_KEY=preprod_api_key

export GCM_REG_ID=optional_reg_id_from_android_device
```

`GCM_REG_ID` is optional and it is a GCM registration ID taken from an Android device or simulator. If provided, it will be used in executing server-to-device messaging tests. Otherwise, these tests will be skipped.

## License

[MIT](LICENSE)
