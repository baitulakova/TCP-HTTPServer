package request

import (
	"net/url"
	"io"
	"os"
	"bufio"
	"strings"
	"github.com/Sirupsen/logrus"
)

type Request struct {
	Method string
	URL *url.URL
	Protocol string
	Headers map[string][]string
	Body io.ReadCloser
}

//Reads request from file, reads line by line
// and return array of single lines
func HandleRequest(filepath string)[]string{
	RequestLines:=make([]string,0)
	file,_:=os.Open(filepath)
	r:=bufio.NewReader(file)
	for {
		s, _, e := r.ReadLine()
		if e==io.EOF{
			break
		}
		RequestLines=append(RequestLines,string(s))
	}
	file.Close()
	os.Remove(filepath)
	RequestLines=RequestLines[:len(RequestLines)-1]
	return RequestLines
}

//reads request lines and converts to Request struct
func FormRequest(RequestLines []string)(req Request){
	if len(RequestLines)>0{
		StartingLine:=strings.Split(RequestLines[0]," ")
		if len(StartingLine)>0{
			req.Method=StartingLine[0]
			u,err:=url.Parse(StartingLine[1])
			if err!=nil{
				logrus.Error("Error parsing url: ",err)
			}
			req.URL=u
			req.Protocol=StartingLine[2]
		}
		//delete first element from slice
		RequestLines=RequestLines[:0]
		for i:=0;i<len(RequestLines);i++{
			header:=strings.Split(RequestLines[i],":")
			req.Headers[header[0]]=strings.Split(header[1],",")
		}
	}
	return
}