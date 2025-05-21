package main

import (
	"fmt"
	"github.com/t6kke/gator/internal/config"
)

func main() {
	conf := config.ReadConfig()
	fmt.Printf("%+v\n", conf)
	conf.SetUser("t6kke")
	new_conf := config.ReadConfig()
	fmt.Printf("%+v\n", new_conf)
}
