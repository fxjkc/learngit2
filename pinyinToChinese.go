package main

import (
	"fmt"
	"os"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"regexp"
)

type qihu struct {
	url string
	qihuResp
}

type qihuResp struct {
	Query string `json:"query"`
	Result []ResultArr
}

type ResultArr struct {
	Word string `json:"word"`
}

func main(){
	var q qihu
	q.url = "https://sug.so.360.cn/suggest?callback=suggest_so&encodein=utf-8" +
		"&encodeout=utf-8&format=json&fields=word&word="
	word := q.clnWord()
	url := q.makeLink(word)
	resp := q.clnResp(url)
	q.display(resp)
}

func (q qihu)display(resp qihuResp)(int,string){
	var a int
	var b string
	for seq,res := range resp.Result{
		fmt.Println(seq,res.Word)
		a=seq
		b=res.Word
	}
	return a,b
}

func (q qihu)clnResp(url string)qihuResp{
	var qi qihuResp
	resp,err := http.Get(url)
	q.fatalErr(err,url)
	readRes,err := ioutil.ReadAll(resp.Body)
	q.fatalErr(err,"ioutil.ReadAll")
	pattern,err := regexp.Compile("{.*}")
	q.fatalErr(err,"compile pattern")
	regexpRes := pattern.FindAllString(string(readRes),-1)
	var unmarshalObj []byte
	if len(regexpRes) > 0{
		unmarshalObj = []byte(regexpRes[0])
	}
	err = json.Unmarshal(unmarshalObj,&qi)
	q.fatalErr(err,string(readRes))
	return qi
}

func (q qihu)fatalErr(err error,msg interface{}){
	if err != nil{
		log.Println("ERROR MSG",msg)
		log.Fatal(err)
	}
}

func (q qihu)makeLink(word string)string{
	return fmt.Sprintf("%s%s",q.url,word)
}

func (q qihu)clnWord()string{
	args := os.Args
	if len(args) == 1{
		return ""
	} else {
		return args[1]
	}
}












