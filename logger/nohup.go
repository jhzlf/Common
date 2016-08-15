package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const BACKUP_SIZE = int64(100 * 1024)

type Notify interface {
	Notify(info string)
}

func BackupNohup(notify Notify) {
	appname := GetAppName()

	fd, err := os.OpenFile("nohup.out", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("backupNohup : open file 'nohup.out' failed.")
		return
	}
	defer fd.Close()

	stat, err := fd.Stat()
	if err != nil {
		fmt.Println("backupNohup : open file 'nohup.out' failed.")
		return
	}

	bufLen := BACKUP_SIZE
	if bufLen > stat.Size() {
		bufLen = stat.Size()
	}

	fmt.Printf("backupNohup : nohup.out size=%d, get len=%d\n", stat.Size(), bufLen)
	fd.Seek(0-bufLen, os.SEEK_END)

	buf := make([]byte, bufLen)
	n, err := fd.Read(buf)
	if err != nil {
		fmt.Println("backupNohup : nohup.out read err= ", err.Error())
		return
	}

	strBuf := string(buf[0:n])

	if strings.Index(strBuf, appname) == -1 {
		fmt.Println("backupNohup : nohup.out not found ", appname)
		return
	}

	if strings.Index(strBuf, "panic:") == -1 {
		fmt.Println("backupNohup : nohup.out not found 'panic:'")
		return
	}

	bakFd, err := os.OpenFile(appname+".bck", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("backupNohup open %s.bak failed, err=%s.", appname, err.Error())
		return
	}
	defer bakFd.Close()

	strHead := fmt.Sprintf("\n----------------------------------%s----------------------------------\n\n", time.Now().String())
	bakFd.Write([]byte(strHead))
	bakFd.Write(buf[0:n])

	if notify != nil {
		strNotify := fmt.Sprintf("%s %s system break", GetIntranetIP(), appname)
		fmt.Println("backupNohup : ", strNotify)
		notify.Notify(strNotify)
	}
}
