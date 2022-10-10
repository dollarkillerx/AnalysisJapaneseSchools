package utils

import (
	"strconv"
	"strings"
)

// WhitespaceOptimization 去除多余空格
func WhitespaceOptimization(str string) string {
	var result string
	for i, v := range str {
		if string(v) == " " {
			if i+1 >= len(str) {
				continue
			} else {
				if string(str[i+1]) == " " {
					continue
				}
			}
		}

		result += string(v)
	}

	return result
}

// ExtractNumbers 提取数字
func ExtractNumbers(str string) []int {
	str = strings.TrimSpace(str)

	var result []int

	var preIsNumber bool
	var ntr string

	var rList = []rune(str)
	for _, v := range rList {
		ic := string(v)
		_, err := strconv.Atoi(ic)
		if err == nil {
			preIsNumber = true
			ntr += ic
		} else {
			preIsNumber = false
		}
		//fmt.Println(ic, "  ", preIsNumber)

		if !preIsNumber {
			if ntr != "" {
				itoa, err := strconv.Atoi(ntr)
				if err == nil {
					result = append(result, itoa)
				}
				ntr = ""
			}
		}
	}

	if ntr != "" {
		atoi, err := strconv.Atoi(strings.TrimSpace(ntr))
		if err == nil {
			result = append(result, atoi)
		}
	}

	return result
}
