package main

import (
	"flag"
	"github.com/wy3148/cabdata/cabapp"
	"log"
)

func main() {

	config := flag.String("config", "", "location of the config file")
	flag.Parse()

	if len(*config) == 0 {
		log.Println("error:config file is not specified, use -config option")
		return
	}

	app := cabapp.NewCabDataApp(*config)
	app.Start()
}
