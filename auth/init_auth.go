package auth

import (
  "fmt"
)

func init() {
  fmt.Println("init auth");
  initAuthRouter();
  initDBConnect();
}
