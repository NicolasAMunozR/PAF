package util

import "fmt"

var areaCodes []string = []string{
	"AYC",
	"CDC",
	"DAF",
	"INT",
	"PEM",
	"PVS",
	"SMA",
}

func GetAreaCode(area int) string {
	i := 1
	fmt.Printf("area: %v\n", area)
	for _, code := range areaCodes {
		fmt.Printf("code: %v %v\n", code, i)
		if area == i {
			return code
		}
		i += 1
	}
	return ""
}
