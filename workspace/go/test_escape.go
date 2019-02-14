package main

import (
	"errors"
	"fmt"
	"strings"
)

func MysqlEscapeString(source string) (string, error) {
	if len(source) == 0 {
		return "", errors.New("source is null")
	}

	j := 0
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		case ':':
			flag = true
			escape = ':'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j]), nil
}

func main() {
	str, _ := MysqlEscapeString("http:::")
	fmt.Printf("str:%s\n", str)
	arr := strings.Split("0:小明\\n1:http://v1.dwstatic.com/bd/201709/09/pic/3ad7a91aa4f1b259d46640a041a00000\n.jpg?w=240&h=320\\n0:小花", "\\n")
	fmt.Printf("-----:+%v", arr)
}
