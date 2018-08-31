package httphelper

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	fullurl := "http://localhost:5001/api/1/address"
	status, body, err := Get(fullurl)
	fmt.Println("status", status)
	fmt.Println("buf", string(body))
	fmt.Println("err", err)
}

func TestPostJSON(t *testing.T) {
	fullurl := "http://localhost:5001/api/1/transfers/allowrevealsecret"
	payloadJSON := "{\"lock_secret_hash\":\"0x55c6593b6ab2d834d1e9b89cc9cdd3c866bf319b04ef6b08548dd85c031376a8\", \"token_address\":\"0xE1E7562D4b55D16073c972414C946bC1113EbeD1\"}"
	status, body, err := PostJSON(fullurl, payloadJSON)
	fmt.Println("status", status)
	fmt.Println("buf", string(body))
	fmt.Println("err", err)
}
