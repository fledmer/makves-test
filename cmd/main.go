package main

import (
	"context"
	"flag"
	"items-service/api/rest"
	"items-service/loader"
	"log"
)

var resourcePath string
var serverPort int

// делать целый конфиг это какой-то rocker science под такую задачу,
func init() {
	flag.StringVar(&resourcePath, "r", "ueba.csv", "path to resource file")
	flag.IntVar(&serverPort, "p", 8090, "rest api server port")
}

func main() {
	flag.Parse()
	loader := loader.NewBufferLoader()
	err := loader.LoadCSVItems(resourcePath)
	if err != nil {
		log.Fatal("failed to Load CSVItems", err)
	}
	rest.New(context.TODO(), loader).Run(serverPort)
	if err != nil {
		log.Fatal("failed to run rest api", err)
	}
}
