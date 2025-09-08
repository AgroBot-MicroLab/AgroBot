export function useWebSocket(url, cb) {
  let ws = null;
  let timer = null;

  const connect = () => {
    ws = new WebSocket(url);

    ws.onopen = () => {
      console.log("WS connected");
      timer = setInterval(() => ws && ws.readyState === 1 && ws.send("ping"), 15000);
    };

    ws.onmessage = (e) => {
      try {
        const data = JSON.parse(e.data);
        cb(data);
      } catch {
        console.log("msg:", e.data);
      }
    };

    ws.onclose = () => {
      console.log("WS closed, reconnecting…");
      if (timer) clearInterval(timer);
      setTimeout(connect, 1000); // reconnect after 1s
    };

    ws.onerror = () => ws && ws.close();
  };

  const send = (obj) => {
    if (ws && ws.readyState === 1) {
      ws.send(JSON.stringify(obj));
    }
  };

  const close = () => {
    if (timer) clearInterval(timer);
    ws && ws.close();
  };

  connect();
  return { send, close };
}

