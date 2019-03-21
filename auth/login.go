package auth

import (
  "context"
  "crypto/md5"
  "database/sql"
  "encoding/hex"
  "errors"
  "fmt"
  "net/http"
  "strconv"
  "time"

  "ic_server/services/db_connect"
  "ic_server/common"
)

func authLogin(name string, pw string) (uint, string) {
  if db_connect.ServerDB == nil {
    fmt.Println("no conn");
    return 0, "no conn";
  }

  if len(name) <= 0 || len(pw) <= 0 {
    return 0, "argue err";
  }
  if len(pw) <= 0 {
    return 0, "argue err";
  }

  rows, err := db_connect.ServerDB.Query("select uid, password from users where name=?;", name);
  if err != nil {
    fmt.Println("query uid error");
    fmt.Println(err);
    return 0, "db error";
  }
  db_connect.ServerDB.Exec("COMMIT;");

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

func loginPage(w http.ResponseWriter, r *http.Request) {
  b, err := common.ReadPage("auth/login.html");
  if err != nil {
    fmt.Println(err);
    w.WriteHeader(http.StatusNotFound);
    w.Write([]byte(""));
    return;
  }

  w.Write(b);
}

func updateUserToken(token string, uid uint) (string) {
  //return "";
  tx, err := db_connect.ServerDB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable});
  if err != nil {
    fmt.Println(err);
    return "db error";
  }

  hasErr := false;
  retErr := "";
  for {
    rows, err := tx.Query("select uid from user_tokens where uid=?;", uid);
    if err != nil {
      fmt.Println(err);
      retErr = "db error";
      hasErr = true;
      break;
    }
  
    if (rows.Next())  {
      if rows.Next() {
        fmt.Println("db err, multi records");
      }
      rows.Close();
      rows2, err := tx.Query("UPDATE user_tokens SET token=? WHERE uid=?;", token, uid);
      if err != nil {
        fmt.Println(err);
        retErr = "db error";
        hasErr = true;
        break;
      }
      rows2.Close();
  
    } else {
      rows.Close();  
      fmt.Println(4);
      rows3, err := tx.Query("INSERT INTO user_tokens(uid, token) VALUES(?, ?);", uid, token);
      if err != nil {
        fmt.Println(err);
        retErr = "db error";
        hasErr = true;
        break;
      }
      rows3.Close();
    }

    break;
  }

  if (hasErr) {
    fmt.Println("transaction error");
    tx.Rollback();
    return retErr;
  }

  if err3 := tx.Commit();err3 != nil {
    return "db error";
  }
  fmt.Println("update token succ");

  return "";
}

func login(w http.ResponseWriter, r *http.Request) {
  r.ParseForm();
  uname := r.Form["username"];
  pwd := r.Form["password"];
  if uname == nil || pwd == nil {
    fmt.Println("args error");
    w.WriteHeader(http.StatusBadRequest );
    w.Write([]byte("args err"));
    return;
  }
  fmt.Println("login", r.Form["username"][0], r.Form["password"][0]);
  uid, err := authLogin(r.Form["username"][0], r.Form["password"][0]);

  if uid == 0 {
    w.WriteHeader(http.StatusUnauthorized);
    w.Write([]byte(err));
    return;
  }
  fmt.Println("get uid:", uid);

  token := genToken(uid);
  if len(token) == 0 {
    w.WriteHeader(http.StatusInternalServerError);
    w.Write([]byte("fatal error"));
    return;
  }
  fmt.Println("get token:", token);

  if db_connect.ServerDB == nil {
    w.WriteHeader(http.StatusInternalServerError );
    w.Write([]byte("mysql err"));
    return;
  }

  err = updateUserToken(token, uid);
  if len(err) != 0 {
    fmt.Println(err);
    w.WriteHeader(http.StatusInternalServerError );
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

func CheckToken(token string) (uint, error) {
  if len(token) == 0 {
    return 0, errors.New("emptr token");
  }

  rows, err := db_connect.ServerDB.Query("select uid from user_tokens where token=?", token);
  if (err != nil) {
    return 0, errors.New("db error");
  }

  defer rows.Close();

  if !rows.Next() {
    return 0, errors.New("token error");
  }

  var uid uint = 0;
  rows.Scan(&uid);
  return uid, nil;
}

func registerPage(w http.ResponseWriter, r *http.Request) {
  b, err := common.ReadPage("auth/register.html");
  if err != nil {
    fmt.Println(err);
    w.WriteHeader(http.StatusNotFound);
    w.Write([]byte(""));
    return;
  }

  w.Write(b);
}

func register(w http.ResponseWriter, r *http.Request) {
  r.ParseForm();
  name := r.Form["username"];
  pwd  := r.Form["password"];

  if name == nil || pwd == nil {
    w.WriteHeader(http.StatusBadRequest);
    w.Write([]byte("args error"));
    return;
  }

  if len(name[0]) == 0 || len(pwd[0]) == 0 {
    w.WriteHeader(http.StatusBadRequest);
    w.Write([]byte("args error"));
    return;
  }

  if _, err := db_connect.ServerDB.Exec("INSERT INTO users(name, password) VALUES(?, ?)", name[0], pwd[0]); err != nil {
    w.WriteHeader(http.StatusBadRequest);
    w.Write([]byte("user exists!"));
    return;
  }

  w.WriteHeader(http.StatusOK);
  w.Write([]byte("success"));
  return;
}