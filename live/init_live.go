package live

import (
  "fmt"
  "net/http"
)

func initLiveRouter() {
  http.HandleFunc("/live/livelist", getLiveList);
  http.HandleFunc("/live/requireRoom", requireLiveRoom);
  http.HandleFunc("/live/requireRoom.html", requireLiveRoomPage);
}


func InitModule() {
  fmt.Println("init live");
  initLiveRouter();
}