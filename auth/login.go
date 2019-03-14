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

  if len(name) <= 0 || len(pw) <= 0 {
    return 0, "argue err";
  }
  if len(pw) <= 0 {
    return 0, "argue err";
  }

  rows, err := auth_conn.Query("select uid, password from users where name=?;", name);
  if err != nil {
    fmt.Println("query uid error");
    fmt.Println(err);
    return 0, "db error";
  }
  auth_conn.Exec("COMMIT;");

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
  //return "";
  tx, err := auth_conn.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable});
  if err != nil {
    fmt.Println(err);
    return "db error";
  }
  defer tx.Commit();
  fmt.Println("begin");

  rows, err := tx.Query("select uid from user_tokens where uid=?;", uid);
  if err != nil {
    fmt.Println(err);
    return "db error";
  }
  defer rows.Close();
  fmt.Println("find uid succ");

  if (rows.Next())  {
    fmt.Println("need update");
    rows2, err := tx.Query("UPDATE user_tokens SET token=? WHERE uid=?;", token, uid);
    if err != nil {
      fmt.Println(err);
      return "db error";
    }
    defer rows2.Close();

    fmt.Println("update succ");
  } else {
    fmt.Println("need insert");
    rows3, err := tx.Query("INSERT INTO user_tokens(uid, token) VALUES(?, ?);", uid, token);
    if err != nil {
      fmt.Println(err);
      return "db error";
    }
    defer rows3.Close();

    fmt.Println("insert succ");
  }

  if (rows.Next()) {
    fmt.Println("db err, multi records");
  }
  return "";
}

func login(w http.ResponseWriter, r *http.Request) {
  r.ParseForm();
  //fmt.Println(r);
  fmt.Println("login", r.Form["username"][0], r.Form["password"][0]);
  uid, err := authLogin(r.Form["username"][0], r.Form["password"][0]);

  if uid == 0 {
    w.WriteHeader(401);
    w.Write([]byte(err));
    return;
  }
  fmt.Println("get uid:", uid);

  token := genToken(uid);
  if len(token) == 0 {
    w.WriteHeader(500);
    w.Write([]byte("fatal error"));
    return;
  }
  fmt.Println("get token:", token);

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
