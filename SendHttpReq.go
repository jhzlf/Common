package Common

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	Http_req_get = iota
	Http_req_post
)

func SendHttpReq(param []byte, funcName string, sendType int, headParam map[string]string) ([]byte, error) {
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
		return []byte(""), err
	}

	if headParam != nil {
		for k, v := range headParam {
			req.Header.Add(k, v)
		}
	}

	// 设置 TimeOut
	client_err := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(30 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*30)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	resp, err := client_err.Do(req)
	if err != nil {
		log.Println(err.Error())
		return []byte(""), err
	}

	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return []byte(""), err
	}
	return resp_body, nil
}

func SendHttpReqTime(param []byte, funcName string, sendType int, headParam map[string]string, t time.Duration) ([]byte, error) {
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
		return []byte(""), err
	}

	if headParam != nil {
		for k, v := range headParam {
			req.Header.Add(k, v)
		}
	}

	// 设置 TimeOut
	client_err := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(t)
				c, err := net.DialTimeout(netw, addr, time.Second*30)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	resp, err := client_err.Do(req)
	if err != nil {
		log.Println(err.Error())
		return []byte(""), err
	}

	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return []byte(""), err
	}
	return resp_body, nil
}
