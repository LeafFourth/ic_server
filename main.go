package main

import (
  "fmt"
  _ "ic_server/auth"
  "ic_server/services"
)

func main() {
  fmt.Println("main()");
  services.SetupServices();
} 
