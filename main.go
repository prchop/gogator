package main

import (
	"fmt"
	"log"

	"github.com/prchop/gogator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Println(err)
	}
	cfg.SetUser("prchop")

	cfg, err = config.Read()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(cfg.DBURL)
	fmt.Println(cfg.UserName)
}
