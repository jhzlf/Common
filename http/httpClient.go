package http

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jhzlf/Common/logs"
)

const (
	Http_req_get = iota
	Http_req_post
)

type HttpClient struct {
	*http.Client
}

func NewHttpClient(t time.Duration) *HttpClient {
	return &HttpClient{
		&http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					deadline := time.Now().Add(t)
					c, err := net.DialTimeout(netw, addr, t)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
			},
		},
	}
}

func (client *HttpClient) SendHttpReq(param []byte, funcName string, sendType int, headParam map[string][]string) ([]byte, error) {
	logs.Debug("send http req ", funcName, headParam, string(param))
	begin := time.Now()

	var req *http.Request
	var err error
	switch sendType {
	case Http_req_post:
		req, err = http.NewRequest("POST", funcName, strings.NewReader(string(param)))
	case Http_req_get:
		funcName = funcName + "?" + string(param)
		req, err = http.NewRequest("GET", funcName, nil)
	}

	if err != nil {
		logs.Errorf("SendHttpReq error", err, time.Since(begin), funcName)
		return []byte(""), err
	}

	if headParam != nil {
		for k, v := range headParam {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		logs.Errorf("SendHttpReq error", err, time.Since(begin), funcName)
		return []byte(""), err
	}

	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Errorf("SendHttpReq error", err, time.Since(begin), funcName)
		return []byte(""), err
	}

	logs.Debug("send http rsp ", time.Since(begin), string(resp_body))
	return resp_body, nil
}
