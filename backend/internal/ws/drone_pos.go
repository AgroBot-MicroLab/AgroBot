package ws

import (
    "sync"
    "encoding/json"
	"net/http"
    "github.com/gorilla/websocket"
)

type Pos struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

var up = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

var conns = struct {
    sync.Mutex
    list []*websocket.Conn
}{}

func DronePosHandle(w http.ResponseWriter, r *http.Request) {
    c, err := up.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    conns.Lock()
    conns.list = append(conns.list, c)
    conns.Unlock()
}

func DronePosBroadcast(p Pos) {
    b, _ := json.Marshal(p)
    conns.Lock()
    for i := 0; i < len(conns.list); {
        if err := conns.list[i].WriteMessage(websocket.TextMessage, b); err != nil {
            conns.list = append(conns.list[:i], conns.list[i+1:]...)
            continue
        }
        i++
    }
    conns.Unlock()
}
