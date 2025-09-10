package ws

import (
    "encoding/json"
	"net/http"
    "github.com/gorilla/websocket"
)

type MissionStatus struct {
    Status bool `json:"status"`
}

func DroneMissionHandle(w http.ResponseWriter, r *http.Request) {
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

func DroneMissionBroadcast(s MissionStatus) {
    b, _ := json.Marshal(s)
    conns.Lock()

    broadcastTo := conns.groups["/drone/mission/status"]

    for i := 0; i < len(broadcastTo); {
        if err := broadcastTo[i].WriteMessage(websocket.TextMessage, b); err != nil {
            broadcastTo = append(broadcastTo[:i], broadcastTo[i+1:]...)
            continue
        }
        i++
    }
    conns.Unlock()
}
