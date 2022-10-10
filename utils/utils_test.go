package utils

import (
	"fmt"
	"testing"
)

func TestP(t *testing.T) {
	px := "asda sadsad   sads  sad   萨达萨  达第三   帝国   sdadefrg    dfrf     "
	fmt.Println(WhitespaceOptimization(px))
}

func TestP2(t *testing.T) {
	px := "31 名（うち専任：4 名）"
	fmt.Println(ExtractNumbers(px))
}

func TestP3(t *testing.T) {
	px := "中国\n\t\t16"
	fmt.Println(ExtractNumbers(px))
}
