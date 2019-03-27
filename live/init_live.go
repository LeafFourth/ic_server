package live

import (
  "fmt"
  "net"
  "net/http"

  "github.com/gwuhaolin/livego/protocol/httpflv"
  "github.com/gwuhaolin/livego/protocol/rtmp"
  "github.com/gwuhaolin/livego/configure"

  "ic_server/defines"
)

var liveStreams *rtmp.RtmpStream;

func initLiveRouter() {
  http.HandleFunc("/live/livelist", getLiveList);
  http.HandleFunc("/live/requireRoom", requireLiveRoom);
  http.HandleFunc("/live/requireRoom.html", requireLiveRoomPage);
}

func initLiveStreams() {
  liveStreams = rtmp.NewRtmpStream(); 
}


func InitModule() {
  fmt.Println("init live");
  initLiveRouter();
  initLiveStreams();
}

func SetupRtmpServer(ln net.Listener) {
  err := configure.LoadConfig(defines.ResRoot + "live/configure.json");
  if nil != err {
    fmt.Println(err);
  }

  rtmp.NewRtmpServer(liveStreams, nil).Serve(ln);
}

func SetupFlvServer(ln net.Listener) {
  httpflv.NewServer(liveStreams).Serve(ln);
}