package dbbin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type File struct {
	Path string
}

func (this *File) readMetadata() (Metadata, []string) {
	fread, _ := os.Open(this.Path)
	defer fread.Close()

	byteArray := make([]byte, 2)
	buff := bytes.NewBuffer(byteArray)

	var metadataSize uint16 = 0
	fread.ReadAt(byteArray, 0)
	binary.Read(buff, binary.LittleEndian, &metadataSize)

	meta := Metadata{}
	byteArray = make([]byte, metadataSize)
	buff = bytes.NewBuffer(byteArray)
	fread.ReadAt(byteArray, 0)
	binary.Read(buff, binary.LittleEndian, &meta)

	var fieldsSize uint16 = 0
	buff = bytes.NewBuffer(byteArray)
	fread.ReadAt(byteArray, 2)
	binary.Read(buff, binary.LittleEndian, &fieldsSize)

	fields := make([]byte, fieldsSize)
	byteArray = make([]byte, fieldsSize)
	buff = bytes.NewBuffer(byteArray)
	fread.ReadAt(byteArray, int64(metadataSize))
	binary.Read(buff, binary.LittleEndian, &fields)

	return meta, strings.Split(string(fields), "|")
}

func (this *File) InsertRecords(records []interface{}) error {
	if len(records) == 0 {
		return nil
	}
	// test
	//os.Remove(this.Path)
	// ...

	// buffer
	byteArray := make([]byte, 0)
	buff := bytes.NewBuffer(byteArray)

	var f *os.File = nil
	var err error = nil
	recordSize := Sizeof(records[0])
	fields := ""
	val := reflect.ValueOf(records[0])
	for i := 0; i < val.NumField(); i++ {
		if fields != "" {
			fields += "|"
		}
		fields += val.Type().Field(i).Name
	}
	isAppend := false

	if FileExists(this.Path) {
		meta, fields_ := this.readMetadata()
		if strings.Join(fields_, "|") != fields {
			return errors.New("Invalid record fields")
		}
		if meta.Size == 0 {
			return errors.New("Invalid metadata")
		}
		if recordSize != int(meta.RecordSize) {
			return errors.New("Invalid record struct size")
		}

		f, err := os.OpenFile(this.Path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		isAppend = true

	} else { // new file
		fmt.Println("record size:", recordSize)

		meta := Metadata{RecordSize: uint16(recordSize), FieldsSize: uint16(len([]byte(fields))), Created: uint64(NowTimestamp())}
		meta.Size = uint16(Sizeof(meta))
		fmt.Println("metadata size:", meta.Size)
		fmt.Println("fields:", meta.FieldsSize)
		err = binary.Write(buff, binary.LittleEndian, meta)
		if err != nil {
			return err
		}
		_, err = buff.Write([]byte(fields))
		if err != nil {
			return err
		}
		//fmt.Println(buff.Bytes())
		//fmt.Println([]byte(fields))

		// creates file
		f, err = os.Create(this.Path)
		if err != nil {
			return errors.New("Couldn't open file")
		}
		defer f.Close()
	}

	// write records
	for _, r := range records {
		err = binary.Write(buff, binary.LittleEndian, r)
		if err != nil {
			return errors.New("Write buffer failed")
		}
		//fmt.Println("buffer:", buff.Bytes())
	}

	if isAppend {
		err = AppendToFile(this.Path, buff.Bytes())
	} else {
		_, err = f.Write(buff.Bytes())
	}
	if err != nil {
		return err
	}

	return nil
}

func (this *File) Walk(record interface{}, callback func(interface{}, int) bool) error {
	if FileExists(this.Path) {
		meta, _ := this.readMetadata()
		initOffset := int(meta.Size + meta.FieldsSize)
		recordSize := Sizeof(record)
		if recordSize != int(meta.RecordSize) {
			return errors.New("Invalid record struct size")
		}
		fread, _ := os.Open(this.Path)
		defer fread.Close()
		byteArray := make([]byte, recordSize)
		count := 0
		stat, err := fread.Stat()
		if err != nil {
			return err
		}
		fileSize := stat.Size()
		for {
			pos := int64(initOffset + (count * recordSize))
			buff := bytes.NewBuffer(byteArray)
			fread.ReadAt(byteArray, pos)
			binary.Read(buff, binary.LittleEndian, record)
			if !callback(record, count) {
				break
			}
			if pos+int64(recordSize) >= fileSize {
				break
			}
			count++
		}
	}
	return nil
}

func OpenFile(path string) *File {
	file := &File{Path: path}
	return file
}
