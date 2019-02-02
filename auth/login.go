package auth

import (
  "fmt"
  "io/ioutil"
  "os"
  "net/http"

  "../defines"
)

func authLogin(name string, pw string) (string, string) {
  if auth_conn == nil {
    fmt.Println("no conn");
    return "", "no conn";
  }

  if len(name) <= 0 {
    return "", "argue err";
  }
  if len(pw) <= 0 {
    return "", "argue err";
  }

  rows, err := auth_conn.Query("select uid, password from users where name=?", name);
  if err != nil {
    fmt.Println(err);
    return "", "db error";
  }

  defer rows.Close();

  if !rows.Next() {
    return "", "user not exist";
  }
  var uid string;
  var pw2 string;
  rows.Scan(&uid, &pw2);
  fmt.Println(uid, pw2);
  if pw != pw2 {
    return "", "pw error";
  }

  return uid, "";
}


func readLoginPage() ([]byte, error) {
  path := defines.ResRoot;
  path += "auth/login.html"
  f, err := os.Open(path);
  if err != nil {
    fmt.Println(err);
    return nil, err;
  }

  return ioutil.ReadAll(f);
}

func loginPage(w http.ResponseWriter, r *http.Request) {
  b, err := readLoginPage();
  if err != nil {
    fmt.Println(nil);
    w.WriteHeader(202);
    return;
  }

  w.Write(b);
}

func login(w http.ResponseWriter, r *http.Request) {
  r.ParseForm();
  //fmt.Println(r);
  fmt.Println("login", r.Form["username"][0], r.Form["password"][0]);
  uid, err := authLogin(r.Form["username"][0], r.Form["password"][0]);

  if len(err) > 0 {
    w.Write([]byte(err));
    return;
  }
  w.Write([]byte("success:" + uid));
}
