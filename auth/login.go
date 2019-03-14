package auth

import (
  "crypto/md5"
  "encoding/hex"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "strconv"
  "time"

  "ic_server/defines"
)

func authLogin(name string, pw string) (uint, string) {
  if auth_conn == nil {
    fmt.Println("no conn");
    return 0, "no conn";
  }

  if len(name) <= 0 {
    return 0, "argue err";
  }
  if len(pw) <= 0 {
    return 0, "argue err";
  }

  rows, err := auth_conn.Query("select uid, password from users where name=?", name);
  if err != nil {
    fmt.Println(err);
    return 0, "db error";
  }

  defer rows.Close();

  if !rows.Next() {
    return 0, "user not exist";
  }
  var uid uint;
  var pw2 string;
  rows.Scan(&uid, &pw2);
  fmt.Println(uid, pw2);
  if pw != pw2 {
    return 0, "pw error";
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
    w.WriteHeader(401);
    w.Write([]byte(err));
    return;
  }

  token := genToken(uid);
  if len(token) == 0 {
    w.WriteHeader(500);
    w.Write([]byte("fatal error"));
    return;
  }

  if auth_conn == nil {
    w.WriteHeader(500 );
    w.Write([]byte("mysql err"));
    return;
  }

  _, err2 := auth_conn.Query("UPDATE user_tokens SET token = ? WHERE uid = ?", token, uid);
  if err2 != nil {
    fmt.Println(err2);
    w.WriteHeader(500 );
    w.Write([]byte("mysql err"));
    return;
  }

  w.Write([]byte(token));
}

func genToken(uid uint) string {
  oriKey := strconv.FormatUint(uint64(uid), 10);
  if len(oriKey) == 0 {
    return "";
  }

  timeStr := time.Now().Format(time.UnixDate);

  h := md5.New();
  h.Write([]byte(oriKey + "_" + timeStr));
  return hex.EncodeToString(h.Sum(nil));
}
