package main

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"
)

type Status struct {
	StatusCode int
	StatusText string
}

type Response struct {
	Method  string
	Status  Status
	Headers map[string]string
	Body    []byte
}

var statusList = map[int]string{
	200: "OK",
	400: "Bad request",
	404: "Not found",
	500: "Internal Server Error",
	501: "Not Implemented",
}

func (res *Response) setStatusCode(code int) error {
	res.Status.StatusCode = code
	res.Status.StatusText = statusList[code]
	if res.Status.StatusText == "" {
		return errors.New("error setting status code ")
	}
	return nil
}

func (res *Response) setHeader(key string, val string) {
	res.Headers[key] = val
}

var badRequest = []byte(`
	<!DOCTYPE html>
	<html>
	<head>
	<title>Bad request</title
	</head>
	<body>
	<h3>Bad request</h3>
	</body>
	</html>
`)

var mainPage = []byte(`
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

var loginPage = []byte(`
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

type loginData struct {
	Login    string
	Password string
}

// Converts request body in Data struct
// The received data is transferred to a template for generating html markup.
func proccessLoginBody(body []byte) loginData {
	d := loginData{}

	credentials := []string{}
	params := strings.Split(string(body), "&") //[login=l,password=p]
	for _, p := range params {
		temp := strings.Split(p, "=")              //[login,l]
		credentials = append(credentials, temp[1]) //[l,p]
	}

	if len(credentials) > 0 {
		d.Login = credentials[0]
		d.Password = credentials[1]
	}
	return d
}

func createResponseBody(d loginData) ([]byte, error) {
	var postAnswer = string(`<!DOCTYPE html>
<html>
<head>
<title>Answer</title>
</head>
<body>
You entered: {{.Login}} {{.Password}}
</body>
</html>
`)

	t, err := template.New("webPage").Parse(postAnswer)
	if err != nil {
		return nil, fmt.Errorf("cannot parse template: %v", err)
	}

	//creates temporary file to store generated html
	f, e := os.Create("temp.txt")
	if e != nil {
		return nil, fmt.Errorf("cannot create temp file: %v", err)
	}
	defer f.Close()

	//writes html in file
	err = t.Execute(f, d)
	if err != nil {
		return nil, fmt.Errorf("cannot generate html markup: %v", err)
	}

	body, err := os.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}
	os.Remove(f.Name())

	return body, nil
}

func createResponse(req *Request) (*Response, error) {
	res := &Response{}

	res.Method = "HTTP/1.0"

	res.Headers = make(map[string]string)

	if req.Method == "GET" {
		res.setStatusCode(200)
		res.setHeader("Content-Type", "text/html; charset=utf-8")
		res.setHeader("Connection", "Keep-Alive")

		switch req.URL.String() {
		case "/":
			res.Body = mainPage
		case "/login":
			res.Body = loginPage
		default:
			res.setStatusCode(404)
			res.Body = badRequest
		}

	} else if req.Method == "POST" {

		if req.Body == nil {
			res.setStatusCode(400)
			res.setHeader("Content-Type", "text/html; charset=utf-8")
			res.Body = badRequest
			return res, nil
		}

		res.setStatusCode(200)
		res.setHeader("Content-Type", "text/html; charset=utf-8")
		res.setHeader("Connection", "Keep-Alive")

		data := proccessLoginBody(req.Body)

		responseBody, err := createResponseBody(data)
		if err != nil {
			return nil, err
		}
		res.Body = responseBody
	}
	return res, nil
}

func (res *Response) toByte() []byte {
	body := fmt.Sprintf("%s %s %s\n", res.Method, strconv.Itoa(res.Status.StatusCode), res.Status.StatusText)
	for k, v := range res.Headers {
		body += fmt.Sprintf("%s: %v\n", k, v)
	}
	body += "\n"
	body += string(res.Body)
	return []byte(body)
}
