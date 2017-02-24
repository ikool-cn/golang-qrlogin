package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func nilFunc(w http.ResponseWriter, r *http.Request) {
	return
}

func getId() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	return ConvertTo62(rndNum)
}

type LoginInfo struct {
	ChanId      string `json:"channel_id"`
	RedirectUrl string `json:"redirect_url"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	info := &LoginInfo{
		ChanId:      getId(),
		RedirectUrl: config.RedirectUrl,
	}

	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println(err)
	}
	err = t.Execute(w, info)
	if err != nil {
		log.Println(err)
	}
}

/*
	generate qrcode and to query specified channel to check if login success
	return:
		channel id
		redirect url
*/
func getChannel(w http.ResponseWriter, r *http.Request) {
	setHeader(w)
	log.Println("getLoginInfo: begin")

	info := &LoginInfo{
		ChanId:      getId(),
		RedirectUrl: config.RedirectUrl,
	}
	body, _ := json.Marshal(info)
	w.Write(body)
}

// check if scan from app
func checkLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("checkLogin: begin")
	setHeader(w)

	r.ParseForm()
	chanId := r.Form.Get("channel_id")
	channel := channelPool.GetChannelById(chanId)
	defer func() {
		channelPool.DestroyChannelById(chanId)
		log.Println("CheckLogin: end")
	}()

	var token string
	select {
	case token = <-channel:

	case <-time.After(time.Duration(config.Timeout) * time.Second):
		w.Write(getResult("timeout, no response from app", 1))
		return
	}

	result := Result{
		Code: 0,
		Msg:  "success",
		Data: token,
	}
	body, _ := json.Marshal(result)

	w.Write(body)
	return
}

/*
	app scan qrcode and then visit this url with token to login
*/
func readyToLogin(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	chanId := r.Form.Get("channel_id")
	token := r.Header.Get("Authorization")
	log.Println("readyToLogin: begin", token)
	defer func() {
		log.Println("readyToLogin: end")
	}()

	if token == "" {
		render(w, "template/failure.html", nil)
		return
	}

	// Bearer token
	token = token[7:]

	//	valid := checkUserValidation(token)
	//	if !valid {
	//		render(w, "template/failure.html", nil)
	//		return
	//	}

	channel := channelPool.GetChannelById(chanId)

	select {
	case channel <- token:
		render(w, "template/success.html", nil)
	case <-time.After(time.Duration(config.Timeout/5) * time.Second):
		channelPool.DestroyChannelById(chanId)
		render(w, "template/failure.html", nil)
		return
	}

}

type Organization struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func checkUserValidation(token string) bool {
	url := config.CreatedOrganizationsUrl
	data, err := request(url, token)
	if err != nil {
		log.Println("Check user validation failed, token= ", token, err)
		return false
	}

	var orgs []Organization
	err = json.Unmarshal(data, &orgs)
	if err != nil {
		log.Println("Unmashal json failed, token=", token, string(data), err)
		return false
	}

	if len(orgs) > 0 {
		return true
	}

	return false
}

func render(w http.ResponseWriter, t string, p interface{}) {
	tpl, err := template.ParseFiles(t)
	if err != nil {
		log.Println(err)
	}

	err = tpl.Execute(w, p)
	if err != nil {
		log.Println(err)
	}
}

/*
browser call this url to confirm login successful
*/
func confirmLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("confirmLogin: begin")
	setHeader(w)

	r.ParseForm()
	chanId := r.Form.Get("channel_id")
	channelPool.DestroyChannelById(chanId)

	w.Write(getResult("login success", 0))
	return
}

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json;charset=utf-8")
}

func getResult(msg string, code int) []byte {
	result := Result{
		Code: code,
		Msg:  msg,
	}
	body, _ := json.Marshal(result)

	log.Println(string(body))
	return body
}

func request(url string, token string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Get response failed from ", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return data, err
}
