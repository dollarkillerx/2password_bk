package main

import (
	"github.com/dollarkillerx/2password/internal/server"
	"github.com/dollarkillerx/2password/internal/utils"

	"log"
)

func main() {
	utils.InitJWT()

	log.SetFlags(log.LstdFlags | log.Llongfile)

	server := server.NewServer()
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
