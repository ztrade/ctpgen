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
	err = generateApi(mdApi, "ctp", "md", *dir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tdApi := filepath.Join(*src, "ThostFtdcTraderApi.h")
	err = generateApi(tdApi, "ctp", "td", *dir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func generateApi(file, pkg, prefix, dir string) (err error) {
	spi, err := ParseSpi(file)
	if err != nil {
		return
	}
	err = spi.Generate(pkg, prefix, dir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	api, err := ParseApi(file)
	if err != nil {
		return
	}
	err = api.Generate(pkg, prefix, dir)
	return
}
