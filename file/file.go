package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// IsPathExist : 路径是否存在
func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

// SaveFile : 将数据保存到文件
func SaveFile(path string, value string) {
	fileHandler, err := os.Create(path)
	defer fileHandler.Close()
	if err != nil {
		return
	}
	buf := bufio.NewWriter(fileHandler)

	fmt.Fprintln(buf, value)

	buf.Flush()
}

// ReadFile : 一次性读取文件 string
func ReadFile(filePth string) string {
	f, err := os.Open(filePth)
	defer f.Close()
	if err != nil {
		return ""
	}
	r, err := ioutil.ReadAll(f)
	return string(r)
}

// ReadFileByte : 一次性读取文件 []byte
func ReadFileByte(filePth string) []byte {
	f, err := os.Open(filePth)
	defer f.Close()
	if err != nil {
		return nil
	}
	r, err := ioutil.ReadAll(f)
	return r
}

// WalkDir : 递归读取目录下所有文件名,存在 *fileList
func WalkDir(dirpath string, fileList *[]string) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			WalkDir(dirpath+"/"+file.Name(), fileList)
			continue
		} else {
			*fileList = append(*fileList, dirpath+"/"+file.Name())
		}
	}
}

// ReadListFile : 按行读文件
func ReadListFile(filename string) []string {
	result := make([]string, 0)
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n')
		line = strings.TrimSpace(line)

		if line != "" {
			result = append(result, line)
		}

		if err == io.EOF {
			break
		}
	}
	return result
}

func ReadListFileDul(filename string) []string {
	strList := ReadListFile(filename)
	var MapExist = make(map[string]bool)
	var result []string
	for _, str := range strList {
		if MapExist[str] == true {
			continue
		}

		result = append(result, str)
		MapExist[str] = true
	}
	return result
}

// WriteListFile : 按行存储文件
func WriteListFile(filename string, value []string) {
	fileHandler, err := os.Create(filename)
	defer fileHandler.Close()
	if err != nil {
		return
	}
	buf := bufio.NewWriter(fileHandler)

	for _, v := range value {
		fmt.Fprintln(buf, v)
	}
	buf.Flush()
}

func DelFile(filename string) {
	os.Remove(filename)
}

// ListDir : 获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}

		files = append(files, dirPth+PthSep+fi.Name())
	}
	return files, nil
}

// AppendToFile :
func AppendToFile(fileName string, content string) error {
	b, _ := IsPathExist(fileName)
	if !b {
		os.Create(fileName)
	}

	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("OpenFile Error", err)
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

// AppendToFile :
func AppendToFileLn(fileName string, content string) error {
	content = content + "\n"
	return AppendToFile(fileName, content)
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// MkDir :
func MkDir(path string) error {
	if path == "" {
		_, file := filepath.Split(os.Args[0])
		path = "/var/log/" + file
	}

	b, err := IsPathExist(path)
	if err != nil {
		fmt.Println("IsPathExist Error", err)
	}

	if !b {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println("Mkdir Error", err)
		}
	}

	return nil
}
