package response

import (
	"github.com/baitulakova/TCP-HTTPServer/request"
	"errors"
)

type Response struct{
	Method string
	Status Status
	Headers map[string]string
	Body []byte
}

type Status struct {
	StatusCode int
	StatusText string
}

var StatusList=map[int]string{
	200:"OK",
	404:"Not found",
	500:"Internal Server Error",
	501:"Not Implemented",
}

func (res *Response) SetStatusCode(code int)error{
	res.Status.StatusCode=code
	res.Status.StatusText=StatusList[code]
	if res.Status.StatusText==""{
		err:=errors.New("Error setting status code ")
		return err
	}
	return nil
}

func (res *Response) SetHeader(HeaderKey string,HeaderValues string){
	res.Headers=make(map[string]string)
	res.Headers[HeaderKey]=HeaderValues
}

var NotFound=[]byte(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>Not Found</title
	</head>
	<body>
	<h5>Not Found</h5>
	</body>
	</html>
`)

var MainPage=[]byte(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>Hello</title
	</head>
	<body>
	<h1>Server is working</h1>
	</body>
	</html>
`)

func FormResponse(req request.Request)(res Response){
	res.Method = "HTTP/1.0"
	if req.Method=="GET" {
		res.SetStatusCode(200)
		res.SetHeader("Content-Type", "text/html; charset=utf-8")
		res.SetHeader("Connection", "Close")
		if req.URL.String()=="/"{
			res.Body=MainPage
		}
	}else if req.Method=="POST"{
		if req.Body==nil{
			res.SetStatusCode(404)
			res.SetHeader("Content-Type","text/html; charset=utf-8")
			res.Body=NotFound
		}
	}
	return res
}

