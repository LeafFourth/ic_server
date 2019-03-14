package auth

import (
  "database/sql"
  _ "github.com/Go-SQL-Driver/MySQL"
  "fmt"
)

var auth_conn *sql.DB = nil;

func initDBConnect() {
  fmt.Println("initDBConnect");
  var err error;
  auth_conn, err = sql.Open("mysql", "ic_server:ic123456@tcp(127.0.0.1:3306)/IC?charset=utf8");
  if err != nil {
      fmt.Println(err);
      return;
  }
  auth_conn.Exec("SET AUTOCOMMIT=0;");
  fmt.Println(auth_conn);
}

