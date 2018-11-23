package response

import (
	"github.com/baitulakova/TCP-HTTPServer/request"
	"errors"
	"strconv"
	"html/template"
	"github.com/Sirupsen/logrus"
	"strings"
	"os"
	"io/ioutil"
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

var statusList=map[int]string{
	200:"OK",
	400:"Bad request",
	404:"Not found",
	500:"Internal Server Error",
	501:"Not Implemented",
}

func (res *Response) SetStatusCode(code int)error{
	res.Status.StatusCode=code
	res.Status.StatusText=statusList[code]
	if res.Status.StatusText==""{
		err:=errors.New("Error setting status code ")
		return err
	}
	return nil
}

func (res *Response) SetHeader(HeaderKey string,HeaderValues string){
	res.Headers[HeaderKey]=HeaderValues
}

var notFound=[]byte(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>Not Found</title
	</head>
	<body>
	<h3>Not Found</h3>
	</body>
	</html>
`)

var mainPage=[]byte(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>Hello</title>
	</head>
	<body>
	<h1>Server is working</h1>
	</body>
	</html>
`)

var loginPage=[]byte(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>Login</title>
	</head>
	<body>
	<h3>Форма входа</h3>
	<form action="/login" method="post">
		Логин:
		<input type="text" name="login">
		<br><br>
		Пароль:
		<input type="password" name="password">
		<br><br>
		<input type="submit" value="Войти">
	</form>
	</body>
	</html>
`)

type Data struct{
	Login string
	Password string
}

func GetBody(req request.Request)(d Data){
	var logAndPass []string
	params:=strings.Split(string(req.Body),"&") //[login=l,password=p]
	for _,p:=range params{
		temp:=strings.Split(p,"=") //[login,l]
		logAndPass=append(logAndPass,temp[1]) //[l,p]
	}
	if len(logAndPass)>0{
		d.Login=logAndPass[0]
		d.Password=logAndPass[1]
	}
	return d
}

func (d *Data) MakePostAnswer() (responseBody []byte){
	var postAnswer =string(
`<!DOCTYPE html>
<html>
<head>
<title>Answer</title>
</head>
<body>
You entered: {{.Login}} {{.Password}}
</body>
</html>
`)
	t,err:=template.New("webpage").Parse(postAnswer)
	if err!=nil{
		logrus.Error("Error parsing html: ",err)
	}
	f,e:=os.Create("temp.txt")
	if e!=nil{
		logrus.Error("Error creating file temp.txt")
	}
	err=t.Execute(f,d)
	f.Close()
	if err!=nil{
		logrus.Error("Error execute: ",err)
	}
	responseBody,_=ioutil.ReadFile("temp.txt")
	os.Remove("temp.txt")
	return responseBody
}

func FormResponse(req request.Request)(res Response){
	res.Method = "HTTP/1.0"
	res.Headers=make(map[string]string)
	if req.Method=="GET" {
		res.SetStatusCode(200)
		res.SetHeader("Content-Type", "text/html; charset=utf-8")
		res.SetHeader("Connection", "Keep-Alive")
		if req.URL.String()=="/"{
			res.Body=mainPage
		}else if req.URL.String()=="/login"{
			res.SetStatusCode(200)
			res.SetHeader("Content-Type", "text/html; charset=utf-8")
			res.SetHeader("Connection", "Keep-Alive")
			res.Body=loginPage
		} else {
			res.SetStatusCode(404)
			res.Body=notFound
		}
	}else if req.Method=="POST"{
		if req.Body==nil{
			res.SetStatusCode(400)
			res.SetHeader("Content-Type","text/html; charset=utf-8")
			res.Body=notFound
		}else {
			res.SetStatusCode(200)
			res.SetHeader("Content-Type","text/html; charset=utf-8")
			data:=GetBody(req)
			responseBody:=data.MakePostAnswer()
			res.Body = responseBody
		}
	}
	return res
}

func (res *Response)MakingResponse()(response []byte){
	responseString:=res.Method+" "+strconv.Itoa(res.Status.StatusCode)+" "+res.Status.StatusText+"\n"
	for k,v:=range res.Headers{
		responseString+=k+": "+v+"\n"
	}
	responseString+="\n"
	responseString+=string(res.Body)
	response=[]byte(responseString)
	return response
}