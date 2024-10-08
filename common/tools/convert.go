package tools

import (
	"log"
	"strconv"
)

func ToInt64(str string) int64 {
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return res
}

func ToFloat64(str string) float64 {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return res
}
