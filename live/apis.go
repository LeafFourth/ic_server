package live

import (
  "context"
  "database/sql"
  "encoding/json"
  "fmt"
  "net/http"

  "ic_server/auth"
  "ic_server/common"
  "ic_server/services/db_connect"
)

type roomInfo struct {
  Uname string;
  Rname string;
  Rid   int;
  Uid   uint
}

type roomsInfo []roomInfo;

func getLiveList(w http.ResponseWriter, r *http.Request) {
  r.ParseForm();
  token := r.Form["token"];
  fmt.Println(token);
  if token == nil {
    w.WriteHeader(http.StatusUnauthorized);
    w.Write([]byte("not login!"));

    return;
  }

  {
    _, err := auth.CheckToken(token[0]);
    if err != nil {
      fmt.Println(err);
      w.WriteHeader(http.StatusUnauthorized);
      w.Write([]byte("not login!"));
      return;
    }
  }

  var infos roomsInfo = make([]roomInfo, 0);
  {
    rows, err := db_connect.ServerDB.Query("SELECT rooms.rid, rooms.uid, rooms.name, users.name  FROM rooms, users WHERE rooms.uid = users.uid");
    if err != nil {
      fmt.Println(err);
      w.WriteHeader(http.StatusInternalServerError);
      w.Write([]byte("db error"));
      return;
    }

    
    for ; rows.Next(); {
      var info roomInfo;
      rows.Scan(&info.Rid, &info.Uid, &info.Rname, &info.Uname);
      infos = append(infos, info);
    }
    rows.Close();
  }

  {
    var data []byte;
    data, err := json.Marshal(infos);
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError);
      w.Write([]byte("data error"));
      return;
    }
    w.Write(data);  
  }
}

func requireLiveRoom(w http.ResponseWriter, r *http.Request) {
  r.ParseForm();
  token := r.Form["token"];
  name  := r.Form["room_name"];
  fmt.Println(token);
  if token == nil {
    w.WriteHeader(http.StatusUnauthorized);
    w.Write([]byte("not login!"));

    return;
  }

  uid, err := auth.CheckToken(token[0]);
  if err != nil {
    fmt.Println(err);
    w.WriteHeader(http.StatusUnauthorized);
    w.Write([]byte("not login!"));
    return;
  }

  var name2 string;
  if name != nil {
    name2 = name[0];
  }

  tx, err2 := db_connect.ServerDB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable});
  if err2 != nil {
    w.WriteHeader(http.StatusInternalServerError);
    return;
  }

  
  var hasEerr = false;
  for {
    rows, err3 := tx.Query("SELECT rid FROM rooms WHERE uid=?", uid);
    if err3 != nil {
      fmt.Println(err3);
      w.WriteHeader(http.StatusInternalServerError);
      w.Write([]byte("db error"));
      hasEerr = true;
      break;
    }
    if rows.Next() {
      rows.Close();
      w.WriteHeader(http.StatusBadRequest);
      w.Write([]byte("has a room already"));
      hasEerr = true;
      break;
    }

    rows.Close();
    {
      _, err3 := tx.Exec("INSERT INTO rooms(uid, name) VALUES(?, ?)", uid, name2);
      if err3 != nil {
        fmt.Println(err3);
        w.WriteHeader(http.StatusBadRequest);
        w.Write([]byte("db error"));
        hasEerr = true;
        break;
      }
    }
    break;
  }

  {
    err3 := tx.Commit();
    if hasEerr {
      return;
    }

    if err3 != nil {
      fmt.Println(err);
      w.WriteHeader(http.StatusBadRequest);
      w.Write([]byte("db error"));
      return;
    }
  }
  fmt.Println(uid, "has a room");
  return;
}

func requireLiveRoomPage(w http.ResponseWriter, r *http.Request) {
  data, err := common.ReadPage("live/requireRoom.html");
  if err != nil {
    fmt.Println(err);
    w.WriteHeader(http.StatusNotFound);
    return;
  }

  w.Write(data);
}