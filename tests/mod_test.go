package tests

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	"github.com/blenderaiser/dbbin"
)

type Record struct {
	Status  uint16
	Billsec uint16
	From    [16]byte
	To      [16]byte
	TechA   [4]byte
	TechB   [4]byte
	//G       uint16
}

func TestViewBin(t *testing.T) {
	fread, _ := os.Open("files/records.bin")
	defer fread.Close()

	byteArray := make([]byte, 2)
	buff := bytes.NewBuffer(byteArray)

	var metadataSize uint16 = 0
	fread.ReadAt(byteArray, 0)
	binary.Read(buff, binary.LittleEndian, &metadataSize)
	//fmt.Println("bit:", billsec, i)
	fmt.Println("metadata size:", metadataSize)

	var fieldsSize uint16 = 0
	buff = bytes.NewBuffer(byteArray)
	fread.ReadAt(byteArray, 2)
	binary.Read(buff, binary.LittleEndian, &fieldsSize)
	fmt.Println("fields size:", fieldsSize)

	fields := make([]byte, fieldsSize)
	byteArray = make([]byte, fieldsSize)
	buff = bytes.NewBuffer(byteArray)
	fread.ReadAt(byteArray, int64(metadataSize))
	binary.Read(buff, binary.LittleEndian, &fields)
	fmt.Println("fields:", string(fields[:]))

	rec := Record{}
	byteArray = make([]byte, dbbin.Sizeof(rec))
	buff = bytes.NewBuffer(byteArray)
	fread.ReadAt(byteArray, int64(metadataSize+fieldsSize)+int64(dbbin.Sizeof(rec)))
	binary.Read(buff, binary.LittleEndian, &rec)
	fmt.Println("rec:", rec)
}

func TestInsert(t *testing.T) {
	file := dbbin.OpenFile("files/records.bin")
	err := file.InsertRecords(append(make([]interface{}, 0), Record{Status: 5, Billsec: 2}))
	if err != nil {
		fmt.Println(err)
	}
}

func TestWalk(t *testing.T) {
	file := dbbin.OpenFile("files/records.bin")
	err := file.Walk(&Record{}, func(record interface{}, index int) bool {
		fmt.Println(index, record.(*Record).Billsec)
		return true //index < 4
	})
	if err != nil {
		fmt.Println(err)
	}
}
