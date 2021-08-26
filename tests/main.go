package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type Record struct {
	IsAnswered bool
	Status     uint16
	Billsec    uint16
	Text       [5]byte
}

func main() {
	os.Remove("file.bin")
	f, err := os.Create("file.bin")
	if err != nil {
		fmt.Println("Couldn't open file")
	}

	var arr [5]byte
	copy(arr[:], "123456")
	s := Record{true, 200, 15000, arr}
	err = binary.Write(f, binary.LittleEndian, s)
	if err != nil {
		fmt.Println("Write failed")
	}
	f.Close()

	// deserialize
	r := &Record{}
	fread, _ := os.Open("file.bin")
	err = binary.Read(fread, binary.LittleEndian, r)
	fread.Close()
	if err != nil {
		fmt.Println("Read failed")
	}
	fmt.Println("deserialize:", r, string(r.Text[:]))

	// read bit
	fread, _ = os.Open("file.bin")
	defer fread.Close()

	buff := make([]byte, 2)
	i, _ := fread.ReadAt(buff, 3)
	buff2 := bytes.NewBuffer(buff)
	var billsec uint16 = 0
	binary.Read(buff2, binary.LittleEndian, &billsec)
	fmt.Println("bit:", billsec, i)
}
