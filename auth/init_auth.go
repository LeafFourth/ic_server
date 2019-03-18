package auth

import (
  "fmt"
  "net/http"
)

func initAuthRouter() {
  http.HandleFunc("/auth/login.html", loginPage);
  http.HandleFunc("/auth/verify", login)
}

func InitModule() {
  fmt.Println("init auth");
  initAuthRouter();
}
