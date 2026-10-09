package main

import (
	"fmt"

	_ "api_mid_AgroCampo/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func main() {
	beego.Run()
}

