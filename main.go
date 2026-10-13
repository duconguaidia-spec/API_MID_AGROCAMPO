<<<<<<< HEAD
package main

import (
	_ "API_MID_AGROCAMPO/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

=======
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

>>>>>>> 1b09d886ad752a742b50653834eb8b404341653a
