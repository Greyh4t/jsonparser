package main

func lex(jsonStr string) []Token {
	var tokens []Token

	jsonR := []rune(jsonStr)
	for i := 0; i < len(jsonR); i++ {
		r := jsonR[i]
		switch r {
		case ' ', '\t', '\n':
			continue
		case '{', '}', '[', ']', ',', ':':
			var token Token
			token.Value = string(r)
			switch r {
			case '{':
				token.Type = typeObjectBegin
			case '}':
				token.Type = typeObjectEnd
			case '[':
				token.Type = typeArrBegin
			case ']':
				token.Type = typeArrEnd
			case ',':
				token.Type = typeComma
			case ':':
				token.Type = typeColon
			}
			tokens = append(tokens, token)
		case '"':
			var str string
			for {
				r := jsonR[i+1]
				if r == '\\' {
					r1 := jsonR[i+2]
					if r1 == 'n' || r1 == 't' || r1 == '\\' || r1 == '"' {
						str += string(r1)
						i += 2
					} else {
						panic("parse err")
					}
				} else if r == '"' {
					tokens = append(tokens, Token{Type: typeStr, Value: str})
					i++
					break
				} else {
					str += string(r)
					i++
				}
			}
		case 't':
			if string(jsonR[i:i+4]) == "true" {
				tokens = append(tokens, Token{Type: typeBool, Value: "true"})
				i += 3
			} else {
				panic("parse failed")
			}
		case 'f':
			if string(jsonR[i:i+5]) == "false" {
				tokens = append(tokens, Token{Type: typeBool, Value: "false"})
				i += 4
			} else {
				panic("parse failed")
			}
		case 'n':
			if string(jsonR[i:i+4]) == "null" {
				tokens = append(tokens, Token{Type: typeNull})
				i += 3
			} else {
				panic("parse failed")
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var numStr = string(r)
			for {
				r := jsonR[i+1]
				if '0' <= r && r <= '9' {
					numStr += string(r)
					i++
				} else {
					break
				}
			}
			tokens = append(tokens, Token{Type: typeNumber, Value: numStr})
		default:
			panic("parse failed")
		}
	}

	return tokens
}
