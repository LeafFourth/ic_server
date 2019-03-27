package services

import (
	"fmt"
	"io/ioutil"
  "net"
	"net/http"
	"os"
	"strings"
	"sync"

	"ic_server/auth"
	"ic_server/defines"
	"ic_server/live"
	"ic_server/services/db_connect"
)

func SetupServices() {
	initServer();
	var wg sync.WaitGroup;
	wg.Add(3);
	go setupHttpService(&wg);
	go setupRtmpService(&wg);
	go setupFlvService(&wg);

	wg.Wait();
}

func initServer() {
	db_connect.InitDBConnect();
	auth.InitModule();
	live.InitModule();
	initGlobalRouter();
}

func initGlobalRouter() {
	http.HandleFunc("/", globalHandle);
}

func handleDefaultPage(w http.ResponseWriter, r *http.Request) bool {
  if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path += "index.html";
		http.DefaultServeMux.ServeHTTP(w, r);
		return true;
	}

	return false;
}

func globalHandle(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm();
	// token := r.Form["token"];
	// if token == nil {
	// 	fmt.Println("args error");
	// 	w.WriteHeader(401);
	// 	w.Write([]byte(""));
	// 	return;
	// }

	// if _, e := auth.CheckToken(token[0]); e != nil {
	// 	fmt.Println(e);
	// 	w.WriteHeader(401);
	// 	w.Write([]byte(""));
	// 	return;
	// }

	fmt.Println("require ", r.URL.Path);
	if handleDefaultPage(w, r) {
		return;
	}



  path := defines.ResRoot + r.URL.Path[1:];
  f, err := os.Open(path);
  if err != nil {
		fmt.Println(err);
		w.WriteHeader(404);
		w.Write([]byte(""));
		return;
	}

	data, err2 := ioutil.ReadAll(f);
	if err2 != nil {
		fmt.Println("read err");
		fmt.Println(err2);
		w.WriteHeader(404);
		w.Write([]byte(""));
		return;
	}
	
	w.Write(data);
}

func setupHttpService(wg *sync.WaitGroup) {
  fmt.Println("setup http server");
	
	err := http.ListenAndServe(":9090", nil);
  if err != nil {
    fmt.Println("connnect error:", err);
	}
	
	wg.Done();
}

func setupRtmpService(wg *sync.WaitGroup) {
	fmt.Println("setup live server");

	bus, err := net.Listen("tcp", ":8033");
	
	if err != nil {
		fmt.Println(err);
		return;
	}

	live.SetupRtmpServer(bus);
	
  wg.Done();
}

func setupFlvService(wg *sync.WaitGroup) {
	fmt.Println("setup httpflv server");

	bus, err := net.Listen("tcp", ":8034");
	
	if err != nil {
		fmt.Println(err);
		return;
	}

	live.SetupFlvServer(bus);
	
  wg.Done();
}

