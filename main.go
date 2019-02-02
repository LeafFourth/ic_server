package main

import (
  "fmt"
  "net/http"

  _ "./auth"
)


func main() {
  fmt.Println("main()");
  err := http.ListenAndServe(":9090", nil);
  if err != nil {
    fmt.Println("connnect error:", err);
  }
}
