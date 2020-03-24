package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/itmo-eve/eden/pkg/einfo"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s file [field:regexp ...]\n", os.Args[0])
		os.Exit(-1)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	q := make(map[string]string)
	for _, a := range os.Args[2:] {
		s := strings.Split(a, ":")
		q[s[0]] = s[1]
	}

	im, err := einfo.ParseZInfoMsg(data)
	if err != nil {
		fmt.Println("ParseZInfoMsg error", err)
		return
	}

	ds := einfo.ZInfoDevSWFind(&im, q)
	if (ds != nil) {
		einfo.ZInfoDevSWPrn(&im, ds)
	}
}
