package auth

import (
  "net/http"
)

func initAuthRouter() {
  http.HandleFunc("/login.html", loginPage);
  http.HandleFunc("/auth", login)
}
