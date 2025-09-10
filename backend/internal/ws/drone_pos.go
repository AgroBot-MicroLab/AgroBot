package ws

import (
    "encoding/json"
	"net/http"
    "github.com/gorilla/websocket"
)

type Pos struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

func DronePosHandle(w http.ResponseWriter, r *http.Request) {
    c, err := up.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    conns.Lock()
    endpoint := r.URL.Path
    conns.groups[endpoint] = append(conns.groups[endpoint], c)
    conns.Unlock()
}

func DronePosBroadcast(p Pos) {
    b, _ := json.Marshal(p)
    conns.Lock()

    broadcastTo := conns.groups["/drone/position"]

    for i := 0; i < len(broadcastTo); {
        if err := broadcastTo[i].WriteMessage(websocket.TextMessage, b); err != nil {
            broadcastTo = append(broadcastTo[:i], broadcastTo[i+1:]...)
            continue
        }
        i++
    }
    conns.Unlock()
}
