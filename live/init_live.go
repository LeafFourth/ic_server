package live

import (
  "fmt"
  "net/http"
)

func initLiveRouter() {
  http.HandleFunc("/live/livelist", getLiveList);
}


func InitModule() {
  fmt.Println("init live");
  initLiveRouter();
}