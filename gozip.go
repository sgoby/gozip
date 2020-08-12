package gozip

import (
	"io"
	"bytes"
	"archive/zip"
	"time"
	"os"
	"fmt"
	"sync"
)

type Gozip struct {
	ZipFileName     string
	zipInnerFileMap map[string]io.Writer //
	dataBuffer      *bytes.Buffer        //
	zipWriter       *zip.Writer
	zipFile         *os.File
	zipReader       *zip.Reader
	deleteFileMap   map[string]bool
	mutex           *sync.Mutex
}

//
func NewGozip() *Gozip{
	//
	gz := &Gozip{}
	gz.zipInnerFileMap = make(map[string]io.Writer)
	gz.dataBuffer = bytes.NewBuffer([]byte{})
	gz.zipWriter = zip.NewWriter(gz.dataBuffer)
	gz.deleteFileMap = make(map[string]bool)
	gz.mutex = &sync.Mutex{}
	//
	return gz
}

//
func (gz *Gozip) Open(fileName string)(err error){
	gz.ZipFileName = fileName
	return gz.openZipFile()
}

//
func (gz *Gozip) RemoveFile(fileName string){
	gz.mutex.Lock()
	defer gz.mutex.Unlock()
	gz.deleteFileMap[fileName] = true
}

//
func (gz *Gozip) WriteFile(fileName string,data []byte) (n int, err error){
	gz.mutex.Lock()
	defer gz.mutex.Unlock()
	//
	if w,ok := gz.zipInnerFileMap[fileName];ok{
		return w.Write(data)
	}
	//
	w,err := gz.createInnerFileWriter(fileName)
	if err != nil{
		return 0,err
	}
	gz.zipInnerFileMap[fileName] = w
	return w.Write(data)
}

//
func (gz *Gozip) Close() (err error) {
	err = gz.Flush()
	if err != nil{
		return err
	}
	//
	if gz.zipReader != nil{
		gz.zipReader = nil
	}
	//
	return gz.closeZipFile()
}

//
func (gz *Gozip) Flush() error {
	gz.mutex.Lock()
	defer gz.mutex.Unlock()
	//
	if gz.zipFile == nil{
		return nil
	}
	//
	err := gz.readeZipFile()
	if err != nil{
		return err
	}
	//
	if gz.zipWriter != nil{
		//
		err = gz.zipWriter.Close()
		if err != nil{
			return err
		}

	}
	//
	dateBuffer := gz.dataBuffer.Bytes()
	//
	gz.reset()
	//
	gz.zipFile.Seek(0,os.SEEK_SET)
	//change file size, no use when  O_APPEND
	err = gz.zipFile.Truncate(int64(len(dateBuffer)))
	if err != nil{
		return err
	}
	_,err = gz.zipFile.Write(dateBuffer)
	return err
}

//reset zipwriter after Flush
func (gz *Gozip) reset() {
	gz.dataBuffer.Reset()
	gz.zipWriter = zip.NewWriter(gz.dataBuffer)
	//
	gz.deleteFileMap = make(map[string]bool)
	gz.zipInnerFileMap = make(map[string]io.Writer)
}

//
func (gz *Gozip) openZipFile() (err error) {
	gz.zipFile, err = os.OpenFile(gz.ZipFileName, os.O_RDWR|os.O_CREATE, 0777)
	return err
}

//
func (gz *Gozip) closeZipFile() (err error) {
	if gz.zipFile != nil{
		err = gz.zipFile.Close()
		gz.zipFile = nil
		return err
	}
	return fmt.Errorf("can not close closed file.")
}

//read hoistory zip file befor every Flush
func (gz *Gozip) readeZipFile() (err error){
	fi,err := gz.zipFile.Stat()
	if err != nil{
		return err
	}
	if fi.Size() <= 0{
		return nil
	}
	//
	// directoryEndSignature    = 0x06054b50
	gz.zipReader,err = zip.NewReader(gz.zipFile,fi.Size()) // 292
	//
	for _, file := range gz.zipReader.File {
		//log.Println(file.FileHeader.Name)
		//
		if _,ok := gz.zipInnerFileMap[file.FileHeader.Name];ok{
			continue
		}
		if _,ok := gz.deleteFileMap[file.FileHeader.Name];ok{
			continue
		}
		//
		wc,err := gz.createInnerFileWriter(file.FileHeader.Name)
		if err != nil {
			return err
		}
		//gz.zipInnerFileMap[file.FileHeader.Name] = wc
		//
		rc, err := file.Open()
		if err != nil {
			return err
		}
		_,err = io.Copy(wc,rc)
		if err != nil {
			return err
		}
		//
		err = rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//
func (gz *Gozip) createInnerFileWriter(fileName string)(w io.Writer, err error){
	header := zip.FileHeader{
		Name:     fileName,
		Method:   zip.Deflate,
		Modified: time.Now(),
	}
	return  gz.createInnerFileIo(header)
}

//
func (gz *Gozip) createInnerFileIo(header zip.FileHeader)(io.Writer,error){
	w,err := gz.zipWriter.CreateHeader(&header)
	if err != nil{
		return nil,err
	}
	return w,err
}