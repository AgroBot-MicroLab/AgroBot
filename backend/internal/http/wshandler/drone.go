package wshandler

import (
    "encoding/json"
	"net/http"
    "github.com/gorilla/websocket"

    "agro-bot/internal"
    "agro-bot/internal/shared"
)

type DroneHandlerWS struct {
    App       *internal.App
}

func (h *DroneHandlerWS) DroneMissionHandle(w http.ResponseWriter, r *http.Request) {
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

func (h *DroneHandlerWS) DroneMissionBroadcast(s shared.MissionEvent) {
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

func (h *DroneHandlerWS)DronePosHandle(w http.ResponseWriter, r *http.Request) {
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

func (h *DroneHandlerWS)DronePosBroadcast(p shared.Pos) {
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
