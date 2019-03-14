package auth

import (
  "fmt"
  "net/http"
)

func initAuthRouter() {
  http.HandleFunc("/login.html", loginPage);
  http.HandleFunc("/auth", login)
}

func InitModule() {
  fmt.Println("init auth");
  initAuthRouter();
  initDBConnect();
}
