package main

import (
	_ "api_mid_AgroCampo/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

