package live

import "fmt"
import "net"

import "github.com/gwuhaolin/livego/protocol/rtmp"
import "github.com/gwuhaolin/livego/configure"

import "ic_server/defines"

func SetupServer(ln net.Listener) {
  err := configure.LoadConfig(defines.ResRoot + "live/configure.json");
  if nil != err {
    fmt.Println(err);
  }

  rtmp.NewRtmpServer(rtmp.NewRtmpStream(), nil).Serve(ln);
}