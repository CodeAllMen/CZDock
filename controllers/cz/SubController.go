/**
  create by yy on 2020/5/22
*/

package cz

import "github.com/astaxie/beego"

type SubController struct {
	beego.Controller
}

func (sub *SubController) StartSub() {

	sub.Data["json"] = "ok"

	sub.ServeJSON()
}
