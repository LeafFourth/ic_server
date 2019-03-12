package auth

import (
  "fmt"
  "net/http"
)

func initAuthRouter() {
  http.HandleFunc("/login.html", loginPage);
  http.HandleFunc("/auth", login)
}

func init() {
  fmt.Println("init auth");
  initDBConnect();
}
