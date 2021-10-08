package HttpServer

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jhzlf/Common/logs"
)

type HttpServer struct {
	Name    string
	Server  *http.Server
	Handler *http.ServeMux
}

func (this *HttpServer) ListenAndServe(port int, timeout int) {
	//	endRunning := make(chan bool, 1)
	go func() {
		this.Server = &http.Server{
			Addr:           ":" + strconv.Itoa(port),
			Handler:        this.Handler,
			ReadTimeout:    time.Duration(timeout) * time.Second,
			WriteTimeout:   time.Duration(timeout) * time.Second,
			MaxHeaderBytes: 1 << 20}
		err := this.Server.ListenAndServe()
		if err != nil {
			panic(fmt.Sprintf("%s ListenAndServe: %v", this.Name, err))
			time.Sleep(100 * time.Microsecond)
			//			endRunning <- true
		}
	}()
	//	<-endRunning
}

func (this *HttpServer) ListenAndServeTLS(port int, timeout int, certFile, keyFile string) {
	//	endRunning := make(chan bool, 1)
	go func() {
		this.Server = &http.Server{
			Addr:           ":" + strconv.Itoa(port),
			Handler:        this.Handler,
			ReadTimeout:    time.Duration(timeout) * time.Second,
			WriteTimeout:   time.Duration(timeout) * time.Second,
			MaxHeaderBytes: 1 << 20}
		err := this.Server.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			panic(fmt.Sprintf("%s ListenAndServeTLS: %v", this.Name, err))
			time.Sleep(100 * time.Microsecond)
			//			endRunning <- true
		}
	}()
	//	<-endRunning
}

func CrossDomain(w http.ResponseWriter, r *http.Request) bool {
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
	return true
}

func BaseParseReq(w http.ResponseWriter, r *http.Request) bool {
	if !CrossDomain(w, r) {
		return false
	}
	err := r.ParseForm()
	if err != nil {
		logs.Errorf("parse param error ", err)
		return false
	}
	return true
}
