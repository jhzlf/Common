package Manager

///*
//const char* build_time(void)
//{
//static const char* psz_build_time = "["__DATE__ " " __TIME__ "]";
//    return psz_build_time;
//}
//*/
//import "C"

import (
	"Common/HttpServer"
	"Common/logger"
	"fmt"
	"io"
	"net/http"
)

var (
	http_manage     *HttpServer.HttpServer
	http_manage_tls *HttpServer.HttpServer
	buildtime       = ""
	buildversion    = ""
)

func Init(port int) {
	manage_port := port
	if manage_port == 0 {
		manage_port = 24438
	}
	http_manage = &HttpServer.HttpServer{Name: "manage", Handler: http.NewServeMux()}

	logger.Info("manage listen at localhost:", manage_port)

	http_manage.Handler.HandleFunc("/SetLogLevel", setLogLevel)
	http_manage.Handler.HandleFunc("/Build", getBuild)
	http_manage.ListenAndServe(manage_port, 360)
}

func Create(port int, build, version string) {
	manage_port := port
	if manage_port == 0 {
		manage_port = 24438
	}
	buildtime = build
	buildversion = version
	http_manage = &HttpServer.HttpServer{Name: "manage", Handler: http.NewServeMux()}

	logger.Info("manage listen at http localhost:", manage_port)

	http_manage.Handler.HandleFunc("/SetLogLevel", setLogLevel)
	http_manage.Handler.HandleFunc("/Build", getBuild)
	http_manage.ListenAndServe(manage_port, 360)
}

func CreateTLS(port int, build, version string, certFile, keyFile string) {
	manage_port := port
	if manage_port == 0 {
		manage_port = 24438
	}
	buildtime = build
	buildversion = version
	http_manage_tls = &HttpServer.HttpServer{Name: "manage", Handler: http.NewServeMux()}

	logger.Info("manage listen at https localhost:", manage_port)

	http_manage_tls.Handler.HandleFunc("/SetLogLevel", setLogLevel)
	http_manage_tls.Handler.HandleFunc("/Build", getBuild)
	http_manage_tls.ListenAndServeTLS(manage_port, 360, certFile, keyFile)
}

func CreateAll(port int, ports int, build, version string, certFile, keyFile string) {
	buildtime = build
	buildversion = version
	http_manage = &HttpServer.HttpServer{Name: "manage", Handler: http.NewServeMux()}
	http_manage_tls = &HttpServer.HttpServer{Name: "manage", Handler: http.NewServeMux()}

	logger.Info("manage listen at http localhost:", port)
	logger.Info("manage listen at https localhost:", ports)

	http_manage.Handler.HandleFunc("/SetLogLevel", setLogLevel)
	http_manage.Handler.HandleFunc("/Build", getBuild)
	http_manage.ListenAndServe(port, 360)

	http_manage_tls.Handler.HandleFunc("/SetLogLevel", setLogLevel)
	http_manage_tls.Handler.HandleFunc("/Build", getBuild)
	http_manage_tls.ListenAndServeTLS(ports, 360, certFile, keyFile)
}

func AddManagerFunc(name string, f func(w http.ResponseWriter, r *http.Request)) {
	if http_manage != nil {
		http_manage.Handler.HandleFunc(name, f)
	}
	if http_manage_tls != nil {
		http_manage_tls.Handler.HandleFunc(name, f)
	}
}

func setLogLevel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var level int
	var pro string
	switch r.FormValue("level") {
	case "debug":
		level = 0
	case "info":
		level = 1
	case "warn":
		level = 2
	case "error":
		level = 3
	case "panic":
		level = 4
	default:
		level = 5
	}

	switch r.FormValue("type") {
	case "file":
		pro = logger.FILE_PROTOCOL
	case "console":
		pro = logger.CONSOLE_PROTOCOL
	case "":
		pro = logger.ALL_PROTOCOL
	}
	if level > 4 {
		io.WriteString(w, string("Set Log Level Faild"))
	}
	logger.Debug(pro, " log set ", level)
	logger.SetLogLevel(pro, level)
	io.WriteString(w, string("Set Log Level Success"))
}

func getBuild(w http.ResponseWriter, r *http.Request) {
	//	buildTime := C.GoString(C.build_time())
	s := fmt.Sprintf("UTC Build Time:%s\n", buildtime)
	s += fmt.Sprintf("Git Commot Hash:%s\n", buildversion)
	io.WriteString(w, s)
}
