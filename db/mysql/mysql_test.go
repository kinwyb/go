package mysql

import (
	"testing"

	"github.com/kinwyb/go/db"
	"github.com/smartystreets/goconvey/convey"
)

//测试mysql
func Test_Mysql(t *testing.T) {
	//mysql://lcfgly:wang93426@tcp(api.zhifangw.cn:3306)/rfid?loc=Local&multiStatements=true
	convey.Convey("测试Mysql", t, func() {
		conn, err := Connect("api.zhifangw.cn:3306", "lcfgly", "wang93426", "rfid", "loc=Local&multiStatements=true")
		convey.So(err, convey.ShouldBeNil)
		row := conn.QueryRows("SELECT iid,company_name,company_user FROM rfid_company_user ORDER BY id DESC LIMIT 5 ")
		row.Error(func(error error) {
			convey.Printf("错误:%s\n", error.Error())
		}).ForEach(func(result map[string]interface{}) bool {
			convey.Printf("%s,%s,%d\n", db.StringDefault(result["company_name"], "无数据"), result["company_user"], result["id"])
			return true
		})
		row = conn.QueryRow("SELECT company_user,id FROM rfid_company_user ORDER BY id ASC")
		row.Error(func(error error) {
			convey.Printf("错误:%s\n", error.Error())
		}).ForEach(func(result map[string]interface{}) bool {
			convey.Printf("%s,%s,%d\n", db.StringDefault(result["company_name"], "无数据"), result["company_user"], result["id"])
			return true
		})
		row = conn.QueryRows("SELECT creator,paypassword,company_user,id,company_password FROM rfid_company_user ORDER BY id ASC LIMIT 1,1 ")
		row.Error(func(error error) {
			convey.Printf("错误:%s\n", error.Error())
		}).ForEach(func(result map[string]interface{}) bool {
			convey.Printf("%s,%s,%d\n", db.StringDefault(result["company_name"], "无数据"), result["company_user"], result["id"])
			return true
		})
		result := conn.Exec("UPDATE rfid_company_user SET enablestate = 1 WHERE id = ? ", 25)
		result.Error(func(i error) {
			convey.Printf("更新错误:%s", i.Error())
		})
		row = conn.QueryRow("SELECT enablestate FROM rfid_company_user WHERE id = ? ", 25)
		row.Error(func(i error) {
			convey.Printf("查询错误:%s", i.Error())
		}).ForEach(func(result map[string]interface{}) bool {
			id := db.Int64Default(result["enablestate"], 0)
			convey.So(id, convey.ShouldEqual, 1)
			return true
		})
		conn.Close()
	})
}

func Benchmark_Mysql(b *testing.B) {
	b.StopTimer()
	conn, err := Connect("api.zhifangw.cn:3306", "lcfgly", "wang93426", "rfid", "loc=Local&multiStatements=true")
	if err != nil {
		b.Fatalf("错误:%s", err.Error())
		return
	}
	b.N = 10
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		row := conn.QueryRows("SELECT id,company_name,company_user FROM rfid_company_user ORDER BY id DESC LIMIT 5 ")
		row.Error(func(error error) {
			b.Fatalf("错误:%s\n", error.Error())
		}).ForEach(func(res map[string]interface{}) bool {
			return true
		})
	}
	conn.Close()
}

//上一版本查询测试: github.com/kinwyb/golang/gosql
func BenchmarkConnect(b *testing.B) {
	b.StopTimer()
	conn, err := Connect("api.zhifangw.cn:3306", "lcfgly", "wang93426", "rfid", "loc=Local&multiStatements=true")
	if err != nil {
		b.Fatalf("错误:%s", err.Error())
		return
	}
	b.N = 10
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		result := conn.QueryRows("SELECT id,company_name,company_user FROM rfid_company_user ORDER BY id DESC LIMIT 5 ")
		if result.HasError() != nil {
			b.Fatalf("错误:%s\n", err.Error())
		}
	}
	conn.Close()
}
