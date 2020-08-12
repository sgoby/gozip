package gozip

import (
	"encoding/binary"
	"testing"
	"log"
	"fmt"
	"time"
)

func Test_gozip(t *testing.T){
	gz := NewGozip()
	err := gz.Open("test.zip")
	if err != nil{
		log.Println(err)
		return
	}
	for i := 140;i < 150 ;i++{
		fileName := fmt.Sprintf("t-%d.txt",i)
		text := fmt.Sprintf("[%d]--askfjalsjklkjljasdfasdfasdfasdfasdf",i)
		_,err := gz.WriteFile(fileName,[]byte(text))
		if err != nil{
			log.Println("@@@",i,err)
			return
		}
		gz.Flush()
		log.Println("End: ",i)
		time.Sleep(time.Second * 2)
	}
	gz.RemoveFile("t-118.txt")
	err = gz.Close()
	if err != nil{
		log.Println(err)
		return
	}
}
func Test_binary(t *testing.T){
	b := make([]byte,2)
	binary.LittleEndian.PutUint16(b, 30)
	log.Printf("%X",b)
	//
	v := binary.LittleEndian.Uint16(b)
	log.Printf("%v",v)
}