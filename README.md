# gozip
gozip is a simple edit zip file library add file to zip file, edit sub file in zip ... 
everything is easy.



### Example:

```go
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
        //Does not exist, will be created, there is an overwrite
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
```

