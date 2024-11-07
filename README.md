
# wssession

`wssession` is a Go package that wraps a [Gorilla WebSocket](https://pkg.go.dev/github.com/gorilla/websocket) connection with session management, providing a simple way to handle WebSocket messages with automatic reconnection, message replay, and session caching.

## Features
- **Session Caching**: Messages are cached to allow replay when connections are lost. Messages are purged every 1-2 minutes in batches for efficiency. 
- **Automatic Reconnection**: Detects lost connections and automatically replays cached messages when a client reconnects.
- **Configurable Handlers**: Allows registering custom message handlers for specific message types.

## Installation
```bash
go get github.com/lordtatty/wssession
```

## Usage

### Setting up the Session Manager
```go
sessMgr := wssession.Mgr{}
sessMgr.RegisterHandler("ping", &pingHandler{})
```

### Upgrading a Connection
```go
u := &websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Allow all connections; adjust for production as needed
        return true
    },
}
conn, err := u.Upgrade(c.Writer, c.Request, nil)
if err != nil {
    log.Printf("error upgrading connection: %s", err)
}
defer conn.Close()

// Serve WebSocket connection
err = sessMgr.Serve(conn)
if err != nil {
    log.Println("Error serving websocket: ", err)
}
```

### Defining a Message Handler
```go
type pingHandler struct{}

func (p *pingHandler) WSHandle(conn wssession.Writer, msg json.RawMessage) {
    if err := conn.SendStr("pong", "Pong!"); err != nil {
        slog.Error("Error writing pong response", "err", err)
    }
}
```

## Roadmap
- Configurable caching behavior and support for multiple cache types.

## License
MIT License
