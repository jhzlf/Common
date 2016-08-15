package logger

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

var (
	pid      int
	progname string
)

func init() {
	pid = os.Getpid()
	paths := strings.Split(os.Args[0], "/")
	paths = strings.Split(paths[len(paths)-1], string(os.PathSeparator))
	progname = paths[len(paths)-1]
	runtime.MemProfileRate = 1
}

func SaveHeapProfile() {
	runtime.GC()
	filename := GetCurrentPath() + "/../prof/" + fmt.Sprintf("heap_%s_%d_%s.prof", progname, pid, time.Now().Format("2006_01_02_03_04_05"))
	offset := strings.LastIndex(filename, "/")
	if offset > 0 {
		pathname := filename[0:offset]
		os.MkdirAll(pathname, 0666)
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	pprof.Lookup("heap").WriteTo(f, 1)
}

func StartHttpProf(port int) {
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
