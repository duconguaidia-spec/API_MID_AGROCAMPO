package main

import (
	_ "API_MID_AGROCAMPO/routers"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func main() {
	beego.Run()
}

