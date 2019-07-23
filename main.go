package main

import (
	"getbox/config"
	"getbox/pbx"
)

func main() {
	config := config.Init(".")
	getbox := &pbx.Getbox{config}
	getbox.Run()
}
