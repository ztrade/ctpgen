package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

var (
	src    = flag.String("s", "", "source dir")
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
	dataTypeFile := filepath.Join(*src, "ThostFtdcUserApiDataType.h")
	structFile := filepath.Join(*src, "ThostFtdcUserApiStruct.h")
	structList, err := ParseStructData(structFile, dataTypeFile)
	err = structList.Generate(*pkg, *dir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	mdApi := filepath.Join(*src, "ThostFtdcMdApi.h")
	spi, err := ParseSpi(mdApi)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	spi.Generate(*pkg, "md", *dir)

}
