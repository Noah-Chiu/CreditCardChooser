package util

import (
	"fmt"
)

// return num in what digits you want
func IntToDigits(num int, digits int) (numStr string) {
	// 將數量改為 001,002...033 等等形式
	numStr = fmt.Sprint(num)
	for len(numStr) < digits {
		numStr = "0" + numStr
	}
	return
}
