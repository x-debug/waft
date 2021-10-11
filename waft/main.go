package main

import (
	"log"
	"os"
	"waft/config"
	"waft/internal"
)

func main() {
	file, err := os.Open("./conf_example/waft.yml")
	if err != nil {
		log.Fatalln("open config file error: ", err)
	}
	defer file.Close()

	conf, err := config.LoadProxyConf(file)
	if err != nil {
		log.Fatalln("load proxy error: ", err)
	}

	proxy := internal.NewProxy(conf)
	log.Fatalln(proxy.Run())
}
