package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

func responder(c *web.Controller, status int, success bool, message string, data interface{}) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{
		"Success": success,
		"Status":  status,
		"Message": message,
		"Data":    data,
	}
	c.ServeJSON()
}