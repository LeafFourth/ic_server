import net

func SetuoServer() {
	bus, err : = net.Listen("tcp", ":80880");
	
	if err != nil {
		return;
	}
	
	
	for  {
		conn, err := bus.Accept();
		if err != nil {
			// handle error
			continue;
		}
		go handleConnection(conn);
	
	}
}

func handleConnection(conn Conn) {
  data := make([]byte, 2);
  n, err := Println(conn.Read(data));
  if err != nil {
    return;
  }
  Println(data[:n]);
}