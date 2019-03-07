package live

import "fmt"
import "net"


func HandleConnect(conn net.Conn) {
  data := make([]byte, 512);
  
  for {
    n, err := conn.Read(data);
    if err != nil {
      fmt.Println(err);
      return;
    }
    fmt.Println(data[:n]);
  }
}