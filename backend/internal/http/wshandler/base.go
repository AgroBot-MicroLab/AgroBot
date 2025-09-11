package wshandler

import (
    "sync"
	"net/http"
    "github.com/gorilla/websocket"
)

var up = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

var conns = struct {
    sync.Mutex
    groups map[string][]*websocket.Conn
}{
    groups: make(map[string][]*websocket.Conn),
}
