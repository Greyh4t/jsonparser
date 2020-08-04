package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	typeNone = iota
	typeObjectBegin
	typeObjectEnd
	typeArrBegin
	typeArrEnd
	typeComma
	typeColon
	typeStr
	typeNumber
	typeBool
	typeNull
)

type Token struct {
	Type  int
	Value string
}

type Tokens struct {
	tokens       []Token
	currentIndex int
}

func (ts *Tokens) Next() Token {
	ts.currentIndex++
	return ts.tokens[ts.currentIndex]
}

func (ts *Tokens) HasMore() bool {
	return ts.currentIndex < len(ts.tokens)-1
}

func (ts *Tokens) Pre() Token {
	return ts.tokens[ts.currentIndex-1]
}

func expect(typeList []int, t int) {
	for _, v := range typeList {
		if v == t {
			return
		}
	}
	panic("parse err")
}

func parse(tokens *Tokens) map[string]interface{} {
	if tokens.HasMore() {
		token := tokens.Next()
		if token.Type == typeObjectBegin {
			return parseObject(tokens)
		}
		panic("parse err")
	}

	return nil
}

func parseArr(tokens *Tokens) []interface{} {
	var r []interface{}
	var expectToken = []int{typeArrEnd, typeBool, typeNull, typeStr, typeNumber, typeObjectBegin}
	for {
		if tokens.HasMore() {
			token := tokens.Next()
			switch token.Type {
			case typeArrEnd:
				expect(expectToken, token.Type)
				return r
			case typeArrBegin:
				expect(expectToken, token.Type)
				r = append(r, parseArr(tokens))
				expectToken = []int{typeComma, typeArrEnd}
			case typeStr, typeBool, typeNumber, typeNull:
				expect(expectToken, token.Type)
				r = append(r, parseValue(token))
				expectToken = []int{typeComma, typeArrEnd}
			case typeComma:
				expect(expectToken, token.Type)
				expectToken = []int{typeStr, typeNumber, typeBool, typeNull, typeArrBegin, typeObjectBegin}
			case typeObjectBegin:
				expect(expectToken, token.Type)
				r = append(r, parseObject(tokens))
				expectToken = []int{typeArrEnd, typeComma}
			}
		}
	}
}

func parseValue(token Token) interface{} {
	switch token.Type {
	case typeBool:
		if token.Value == "true" {
			return true
		} else {
			return false
		}
	case typeNull:
		return nil
	case typeNumber:
		n, err := strconv.Atoi(token.Value)
		if err != nil {
			panic(err)
		}
		return n
	case typeStr:
		return token.Value
	}
	panic("parse err")
}

func parseObject(tokens *Tokens) map[string]interface{} {
	var r = map[string]interface{}{}
	var key string
	var expectToken = []int{typeObjectEnd, typeStr}
	for {
		if tokens.HasMore() {
			token := tokens.Next()
			switch token.Type {
			case typeObjectEnd:
				expect(expectToken, token.Type)
				return r
			case typeObjectBegin:
				expect(expectToken, token.Type)
				r[key] = parseObject(tokens)
				expectToken = []int{typeObjectEnd, typeComma}
			case typeStr:
				expect(expectToken, token.Type)
				if tokens.Pre().Type == typeObjectBegin || tokens.Pre().Type == typeComma {
					key = token.Value
					expectToken = []int{typeColon} //:
				} else if tokens.Pre().Type == typeColon {
					expectToken = []int{typeComma, typeObjectEnd} //,}
					r[key] = token.Value
				}
			case typeBool, typeNull, typeNumber:
				expect(expectToken, token.Type)
				r[key] = parseValue(token)
				expectToken = []int{typeComma, typeObjectEnd}
			case typeColon:
				expect(expectToken, token.Type)
				expectToken = []int{typeStr, typeBool, typeNumber, typeNull, typeArrBegin, typeObjectBegin}
			case typeComma:
				expect(expectToken, token.Type)
				expectToken = []int{typeStr}
			case typeArrBegin:
				expect(expectToken, token.Type)
				r[key] = parseArr(tokens)
				expectToken = []int{typeObjectEnd, typeComma}
			}
		}
	}
}

func main() {
	jsonStr := `
	{
	    "body": "",
	    "nic": [{
	            "ip": "172.17.0.1",
	            "name": "docker0"
	        }, {
	            "ip": "192.168.109.128",
	            "name": "aaa\\bbb\"ccc\\"
	        },null
	    ],
	    "header": {
	        "referer": null,
	        "accept-language": 12345,
	        "cookie": true,
			"host": false
	    }
	}`
	tokens := lex(jsonStr)
	// data, _ := json.MarshalIndent(tokens, "", "  ")
	// fmt.Println(string(data))
	obj := parse(&Tokens{tokens: tokens, currentIndex: -1})
	data1, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(data1))
	fmt.Println("end")
}
