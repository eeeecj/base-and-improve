package sql

import (
	"github.com/balance/boolfilter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

var sqlDSN = "root:eeeecj@tcp(localhost:3306)/gobasic?charset=utf8mb4&parseTime=True&loc=Local"

func TestSqlFilter(t *testing.T) {
	// Init gorm.DB
	db, err := gorm.Open(mysql.Open(sqlDSN))
	if err != nil {
		t.Fatal(err)
	}
	// Init SQLFilter
	sqlFilter, err := NewSqlFilter(db, "bloom", "test", 1000, boolfilter.DefaultHash...)
	if err != nil {
		t.Fatal(err)
	}
	// Push 250-300 numbers to the filter.
	// 把250-300的数字压入过滤器
	for i := 200; i <= 300; i++ {
		sqlFilter.Push([]byte(strconv.Itoa(i)))
	}
	sqlFilter.Write()
}
func TestSqlFilterExist(t *testing.T) {
	db, err := gorm.Open(mysql.Open(sqlDSN))
	if err != nil {
		t.Fatal(err)
	}
	sqlFilter, err := NewSqlFilter(db, "bloom", "test", 1000, boolfilter.DefaultHash...)
	if err != nil {
		t.Fatal(err)
	}
	// 280-300 should exist in filter, and 301-320 doesn't.
	// 280-300应该在过滤器中，而301-320不应该在。
	for i := 280; i < 320; i++ {
		r, _ := sqlFilter.Exists([]byte(strconv.Itoa(i)))
		t.Logf("%d: %t", i, r)
	}
}
