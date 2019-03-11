package live

import "fmt"
import "github.com/gwuhaolin/livego/protocol/rtmp"
import "net"

func SetupServer(ln net.Listener) {
  rtmp.NewRtmpServer(ln).Start();


}