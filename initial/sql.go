package initial

import (
	"fmt"

	_ "github.com/angui001/CZDock/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func InitSql() {
	user := beego.AppConfig.String("psqluser")
	passwd := beego.AppConfig.String("psqlpass")
	host := beego.AppConfig.String("psqlurls")
	port, err := beego.AppConfig.Int("psqlport")
	dbname := beego.AppConfig.String("psqldb")
	if nil != err {
		port = 5432
	}
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	orm.RegisterDriver("postgres", orm.DRPostgres) // 注册驱动
	orm.RegisterDataBase("default",
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
			user, passwd, dbname, host, port))
	orm.DefaultRowsLimit = -1
	orm.RunSyncdb("default", false, true)
}
