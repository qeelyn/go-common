package gormx_test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/qeelyn/go-common/gormx"
	"testing"
)

var (
	Db       *gorm.DB
	mysql    = "mysql"
	mysqlDsn = "root:123456@tcp(localhost:3306)/yak"
)

func setDefaultDb(t *testing.T) {
	cfg := map[string]interface{}{
		"dialect": mysql,
		"dsn":     mysqlDsn,
	}
	var err error
	Db, err = gormx.NewDb(cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewDb(t *testing.T) {
	cfg := map[string]interface{}{
		"dialect":         mysql,
		"dsn":             mysqlDsn,
		"maxidleconns":    10,
		"maxopenconns":    10,
		"connmaxlifetime": 200,
	}
	if Db, err := gormx.NewDb(cfg); err != nil {
		t.Fatal(err)
	} else {
		Db.Close()
	}
}

func TestBuilder_Where(t *testing.T) {
	setDefaultDb(t)
	bl := gormx.NewBuilder(Db)
	wstr := "id = ? and date between ? and ?"
	wps := map[string]string{
		"1":   "2017-01-01",
		"0":   "1",
		"2":   "2017-02-02",
		"sec": "adfasdf",
	}
	bl.Where(wstr, wps)
	fmt.Println(bl.Prepare().SubQuery())
}
