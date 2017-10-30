package main

//1.0

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type dictNeed struct {
	url            string
	key            string
	word           string
	types          string
	emptyRes       string
	pronounciation string
	toEng          string
}

type Jinshan struct {
	English
	Chinese
}

type English struct {
	WordName string `json:"word_name"`
	//单词的各种时态
	Symbols []SymbolsEnglish `json:"symbols"`
}

type SymbolsEnglish struct {
	//英式音标
	Ph_en string `json:"ph_en"`
	//美式音标
	Ph_am string      `json:"ph_am"`
	Parts []PartsStru `json:"parts"`
}

type PartsStru struct {
	Part  string   `json:"part"`
	Means []string `json:"means"`
}

type Chinese struct {
	Word_name string           `json:"word_name"` //#所查的词
	Symbols   []SymbolsChinese `json:"symbols"`   //#词各种信息 下面字段都是这个字段下面的
}

type SymbolsChinese struct {
	Word_symbol string        `json:"word_symbol"`
	Part        []PartChinese `json:"parts"`
}

type PartChinese struct {
	ChineseMeans []ChineseMeanStru `json:"means"`
}

type ChineseMeanStru struct {
	Word_mean string `json:"word_mean"`
}

func main() {
	d := dictNeed{}
	d.key = "9E38A48B2D79AF5E10EC5D8B0AC63214"
	d.url = "http://dict-co.iciba.com/api/dictionary.php?"
	d.types = "json"
	d.emptyRes = "look up result is empty"
	d.pronounciation = "Pronoun: "
	d.toEng = "Eng: "
	d.word = d.getWord()
	jinshanResp := d.getMeaning()
	jinshanStru := d.parseMeans(jinshanResp)
	d.display(jinshanStru)
}

func (d dictNeed) getWord() string {
	if len(os.Args) == 1 {
		return ""
	} else {
		return os.Args[1]
	}
}

func (d dictNeed) getMeaning() []byte {
	url := fmt.Sprintf("%sw=%s&key=%s&type=%s", d.url, d.word, d.key, d.types)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("getting resp meets error", err)
	}
	resRead, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read error", err)
	}
	return resRead
}

func (d dictNeed) judgeChEng() bool {
	pattern, err := regexp.Compile("[A-Za-z]")
	if err != nil {
		log.Fatal(err)
	}
	res := pattern.MatchString(d.word)
	return res
}

func (d dictNeed) parseMeans(resp []byte) Jinshan {
	var j Jinshan
	if d.judgeChEng() {
		json.Unmarshal(resp, &j.English)
	} else {
		json.Unmarshal(resp, &j.Chinese)
	}

	return j
}

func (d dictNeed) display(j Jinshan) {
	fmt.Println("word: ", d.word)
	if d.judgeChEng() {
		d.displayEng(j)
	} else {
		d.displayChinese(j)
	}
}

func (d dictNeed) displayChinese(j Jinshan) {
	symbol := j.Chinese.Symbols
	if len(symbol) == 0 {
		fmt.Println(d.emptyRes)
		os.Exit(0)
	}
	fmt.Println(d.pronounciation, symbol[0].Word_symbol)
	for _, means := range symbol[0].Part {
		for i, mean := range means.ChineseMeans {
			iStr := strconv.FormatInt(int64(i), 10)
			prefix := d.toEng + iStr
			fmt.Println(prefix, mean.Word_mean)
		}
	}
}

func (d dictNeed) displayEng(j Jinshan) {
	symbol := j.English.Symbols
	if len(symbol) == 0 {
		fmt.Println(d.emptyRes)
		os.Exit(0)
	}
	fmt.Println(d.pronounciation, symbol[0].Ph_am)
	for _, part := range symbol[0].Parts {
		var means string
		for _, mean := range part.Means {
			means += mean + ","
		}
		fmt.Println(part.Part, means)
	}
}
