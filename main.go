package main

import (
	"rosella/setup"

	"github.com/astaxie/beego"
)

func main() {
	setup.Setup()
	beego.Run()
}
