package Common

import (
	"Common/logger"
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

// func SendHttpReq(param []byte, funcName string, sendType int, headParam map[string]string) ([]byte, error) {
// 	var req *http.Request
// 	var err error
// 	switch sendType {
// 	case Http_req_post:
// 		req, err = http.NewRequest("POST", funcName, strings.NewReader(string(param)))
// 	case Http_req_get:
// 		funcName = funcName + "?" + string(param)
// 		req, err = http.NewRequest("GET", funcName, nil)
// 	}

// 	if err != nil {
// 		return []byte(""), err
// 	}

// 	if headParam != nil {
// 		for k, v := range headParam {
// 			req.Header.Add(k, v)
// 		}
// 	}

// 	// 设置 TimeOut
// 	client_err := http.Client{
// 		Transport: &http.Transport{
// 			Dial: func(netw, addr string) (net.Conn, error) {
// 				deadline := time.Now().Add(30 * time.Second)
// 				c, err := net.DialTimeout(netw, addr, time.Second*30)
// 				if err != nil {
// 					return nil, err
// 				}
// 				c.SetDeadline(deadline)
// 				return c, nil
// 			},
// 		},
// 	}

// 	resp, err := client_err.Do(req)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return []byte(""), err
// 	}

// 	logger.Info(resp.Header)
// 	defer resp.Body.Close()

// 	resp_body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return []byte(""), err
// 	}
// 	return resp_body, nil
// }

// func SendHttpReqTime(param []byte, funcName string, sendType int, headParam map[string]string, t time.Duration) ([]byte, error) {
// 	var req *http.Request
// 	var err error
// 	switch sendType {
// 	case Http_req_post:
// 		req, err = http.NewRequest("POST", funcName, strings.NewReader(string(param)))
// 	case Http_req_get:
// 		funcName = funcName + "?" + string(param)
// 		req, err = http.NewRequest("GET", funcName, nil)
// 	}

// 	if err != nil {
// 		return []byte(""), err
// 	}

// 	if headParam != nil {
// 		for k, v := range headParam {
// 			req.Header.Add(k, v)
// 		}
// 	}

// 	// 设置 TimeOut
// 	client_err := http.Client{
// 		Transport: &http.Transport{
// 			Dial: func(netw, addr string) (net.Conn, error) {
// 				deadline := time.Now().Add(t)
// 				c, err := net.DialTimeout(netw, addr, time.Second*t)
// 				if err != nil {
// 					return nil, err
// 				}
// 				c.SetDeadline(deadline)
// 				return c, nil
// 			},
// 		},
// 	}

// 	resp, err := client_err.Do(req)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return []byte(""), err
// 	}

// 	defer resp.Body.Close()

// 	resp_body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return []byte(""), err
// 	}
// 	return resp_body, nil
// }

func GetHttpClient(t time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(t * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*t)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
}

func SendHttpReq(param []byte, funcName string, sendType int, headParam map[string]string, client *http.Client) ([]byte, error) {
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

	resp, err := client.Do(req)
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

func BaseParseReq(w http.ResponseWriter, r *http.Request) bool {
	if origin := r.Header.Get("Origin"); origin != "" {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token,Accept,X-Requested-With")
	}

	if r.Method == "OPTIONS" {
		return false
	}

	err := r.ParseForm()
	if err != nil {
		logger.Warn("parse param error ", err)
		return false
	}
	return true
}
