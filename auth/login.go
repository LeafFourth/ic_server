package auth

import (
  "context"
  "crypto/md5"
  "database/sql"
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
    fmt.Println(err);
    w.WriteHeader(404);
    w.Write([]byte(""));
    return;
  }

  w.Write(b);
}

func updateUserToken(token string, uid uint) (string) {
  tx, err := auth_conn.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable});
  if err != nil {
    fmt.Println(err);
    return "db error";
  }

  rows, err := tx.Query("select uid from user_tokens where uid=?", uid);
  if err != nil {
    fmt.Println(err);
    return "db error";
  }
  defer rows.Close();

  if (rows.Next())  {
    rows2, err := tx.Query("UPDATE user_tokens SET token=", token);
    if err != nil {
      fmt.Println(err);
      return "db error";
    }

    defer rows2.Close();
  } else {
    rows3, err := tx.Query("INSERT INTO user_tokens(uid, token) VALUES(?, ?)", uid, token);
    if err != nil {
      fmt.Println(err);
      return "db error";
    }

    defer rows3.Close();
  }

  if (rows.Next()) {
    fmt.Println("db err, multi records");
  }
  
  err = tx.Commit();
  if err != nil {
    fmt.Println(err);
    return "db error";
  }

  return "";
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

  err = updateUserToken(token, uid);
  if len(err) != 0 {
    fmt.Println(err);
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
