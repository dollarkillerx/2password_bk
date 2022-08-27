package main

import (
	"github.com/dollarkillerx/2password/internal/server"
	"github.com/dollarkillerx/2password/internal/utils"

	"log"
)

func main() {
	utils.InitJWT()

	log.SetFlags(log.LstdFlags | log.Llongfile)

	s := server.NewServer()
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}
