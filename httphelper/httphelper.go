package httphelper

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// HTTPRequestTimeoutDefault : 默认http请求超时时间
var HTTPRequestTimeoutDefault = 180 * time.Second

// PostJSON :
func PostJSON(fullURL string, payloadJSON string) (status int, body []byte, err error) {
	client := newClient()
	var reqBody io.Reader
	if payloadJSON == "" {
		reqBody = nil
	} else {
		reqBody = strings.NewReader(payloadJSON)
	}
	req, err := http.NewRequest("POST", fullURL, reqBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)
	defer func() {
		if req.Body != nil {
			err = req.Body.Close()
		}
		if resp != nil && resp.Body != nil {
			err = resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	n := 0
	var buf [4096 * 1024]byte
	n, err = resp.Body.Read(buf[:])
	if err != nil && err.Error() == "EOF" {
		err = nil
	}
	return resp.StatusCode, buf[:n], err
}

// Get :
func Get(fullURL string) (status int, body []byte, err error) {
	client := newClient()
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)
	defer func() {
		if req.Body != nil {
			err = req.Body.Close()
		}
		if resp != nil && resp.Body != nil {
			err = resp.Body.Close()
		}
	}()
	if err != nil {
		return
	}
	n := 0
	var buf [4096 * 1024]byte
	n, err = resp.Body.Read(buf[:])
	if err != nil && err.Error() == "EOF" {
		err = nil
	}
	return resp.StatusCode, buf[:n], err
}

func newClient() *http.Client {
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, 5*time.Second) // default dial timeout 5 seconds
				if err != nil {
					return nil, err
				}
				err = c.SetDeadline(time.Now().Add(HTTPRequestTimeoutDefault))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}
	return &client
}
