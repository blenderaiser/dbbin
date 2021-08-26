package dbbin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

var (
	RootPath, _ = os.Getwd()
)

// FILE
func PathJoin(pathA string, pathB string) string {
	return filepath.Join(pathA, pathB)
}

func RelToAbsPath(path string) string {
	return filepath.Join(RootPath, path)
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FileRead(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func FileCreate(path string, data string) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err2 := f.WriteString(data)
	if err2 != nil {
		return err2
	}

	return nil
}

func AppendToFile(fileName string, content []byte) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	} else {
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt(content, n)
	}
	defer f.Close()
	return err
}

// STRUCT
func Sizeof(v interface{}) int {
	size := 0
	s := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).CanInterface() {
			//size += Sizeof(s.Field(i).Interface())
			size += int(reflect.TypeOf(s.Field(i).Interface()).Size())
		}
	}
	return size
}

// TIME
func NowTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
