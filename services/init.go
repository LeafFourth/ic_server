package services

import "ic_server/live"
import "fmt"
import "net"
import "net/http"
import "sync"

func SetupServices() {
	var wg sync.WaitGroup;
	wg.Add(2);
	go setupHttpService(&wg);
	go setupLiveTcpService(&wg);

	wg.Wait();
 }

func setupHttpService(wg *sync.WaitGroup) {
  fmt.Println("setup http server");
	
	err := http.ListenAndServe(":9090", nil);
  if err != nil {
    fmt.Println("connnect error:", err);
	}
	
	wg.Done();
}

func setupLiveTcpService(wg *sync.WaitGroup) {
	fmt.Println("setup live server");

	bus, err := net.Listen("tcp", ":8033");
	
	if err != nil {
		fmt.Println(err);
		return;
	}

	live.SetupServer(bus);
	
  wg.Done();
}