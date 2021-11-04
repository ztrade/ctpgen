package main

import (
	"flag"
	"fmt"
)

var (
	src    = flag.String("s", "", "source spi header file")
	pkg    = flag.String("pkg", "ctp", "pkg name")
	prefix = flag.String("p", "", "generate code prefix, md/td")
	dir    = flag.String("o", "./", "generate code dir")
)

func main() {
	flag.Parse()
	if *src == "" || *prefix == "" || *pkg == "" {
		flag.PrintDefaults()
		return
	}
	spi, err := ParseSpi(*src)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = spi.Generate(*pkg, *prefix, *dir)
	if err != nil {
		fmt.Println(err.Error())
	}
}
