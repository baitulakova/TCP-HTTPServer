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
	Body []byte
}

//Reads request from file, reads line by line
// and return array of single lines
func HandleRequest(filepath string)[]string{
	logrus.Info("Started handling request")
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
	return RequestLines
}

//Removes given string from slice
func RemoveLine(List []string,stringToRemove string)(ResultList []string){
	for _,line:=range List{
		if line==stringToRemove{
			continue
		}else {
			ResultList=append(ResultList,line)
		}
	}
	return
}

//reads request lines and converts to Request struct
func FormRequest(RequestLines []string)(req Request){
	req.Headers=make(map[string][]string)
	if len(RequestLines)>0{
		firstLine:=RequestLines[0] //method,URL,protocol version
		StartingLine:=strings.Split(firstLine," ")
		if len(StartingLine)>0{
			req.Method=StartingLine[0]
			u,err:=url.Parse(StartingLine[1])
			if err!=nil{
				logrus.Error("Error parsing url: ",err)
			}
			req.URL=u
			req.Protocol=StartingLine[2]
		}
		//delete first line from request
		ListOfHeaders:=RemoveLine(RequestLines,firstLine)
		if len(ListOfHeaders)>0 {
			for i := 0; i < len(ListOfHeaders); i++ {
				header := strings.SplitN(ListOfHeaders[i], ":", 2)
				if len(header) == 2 {
					keyHeader := header[0]
					valueHeader := header[1]
					var values []string
					values = append(values, valueHeader)
					req.Headers[keyHeader] = values
				}else{
					if header[0]==""{
						continue
					}else {
						req.Body=[]byte(header[0])
					}
				}
			}
		}else{
			logrus.Error("Length of headers is 0")
		}
	}
	return
}