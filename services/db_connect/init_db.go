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
  //fmt.Println(ServerDB);

  createTables();
}

func createTables() {
  sqlCmd1 := `CREATE table users(uid int unsigned PRIMARY KEY AUTO_INCREMENT, name char(100) NOT NULL unique, password char(100)) AUTO_INCREMENT=1000;`;
  sqlCmd2 := `CREATE table rooms(rid int unsigned PRIMARY KEY AUTO_INCREMENT, uid int unsigned, name char(100)) AUTO_INCREMENT=1000;`;
  sqlCmd3 := `CREATE table user_tokens(token char(100) PRIMARY KEY, uid int unsigned);`;   

  var err error;
  _, err = ServerDB.Exec(sqlCmd1);
  fmt.Println(err);

  _, err = ServerDB.Exec(sqlCmd2);
  fmt.Println(err);

  _, err = ServerDB.Exec(sqlCmd3);
  fmt.Println(err);
}