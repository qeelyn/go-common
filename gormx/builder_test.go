package gormx_test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qeelyn/go-common/gormx"
	"testing"
)

var (
	Db *gorm.DB
)

func init() {
	Db, _ = gorm.Open("mysql", "root:@tcp(localhost:3306)/test")
}

func TestBuilder_Where(t *testing.T) {
	bl := gormx.NewBuild(Db)
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
