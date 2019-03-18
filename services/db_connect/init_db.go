package db_connect

import (
	"fmt"
	"database/sql"

	_ "github.com/Go-SQL-Driver/MySQL"
)

var ServerDB *sql.DB = nil;


func InitDBConnect() {
  fmt.Println("initDBConnect");
  var err error;
  ServerDB, err = sql.Open("mysql", "ic_server:ic123456@tcp(127.0.0.1:3306)/IC?charset=utf8");
  if err != nil {
      fmt.Println(err);
      return;
  }
  //ServerDB.Exec("SET AUTOCOMMIT=0;");
  fmt.Println(ServerDB);
}