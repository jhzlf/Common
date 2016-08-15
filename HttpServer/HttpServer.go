package HttpServer

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
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
