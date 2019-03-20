package auth

import (
  "fmt"
  "net/http"
)

func initAuthRouter() {
  http.HandleFunc("/auth/login.html", loginPage);
  http.HandleFunc("/auth/verify", login);
  http.HandleFunc("/auth/register.html", registerPage);
  http.HandleFunc("/auth/signup", register);
}

func InitModule() {
  fmt.Println("init auth");
  initAuthRouter();
}
