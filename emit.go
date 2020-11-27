package cdpapp

import (
  "context"
  "encoding/json"
  "github.com/mafredri/cdp/rpcc"
  "reflect"
  "sync"
)

type eventHandler = func(message json.RawMessage)

type eventIns struct {
  ctx context.Context
  cancel context.CancelFunc
  stream rpcc.Stream
  handler []eventHandler
}

type Emitter struct {
  onMu sync.Mutex
  ctx context.Context
  conn *rpcc.Conn
  streams map[string]*eventIns
}

func NewEmitter(ctx context.Context, conn *rpcc.Conn) *Emitter {
  return &Emitter{
    ctx: ctx,
    conn: conn,
    streams: map[string]*eventIns{},
  }
}

func (e *Emitter) On(method string, handle eventHandler) (err error) {
  e.onMu.Lock()
  defer e.onMu.Unlock()

  ins, ok := e.streams[method]
  if !ok {
    stream, err := rpcc.NewStream(e.ctx, method, e.conn)

    if err != nil {
      return err
    }
    ins = &eventIns{
      ctx:     context.Background(),
      stream:  stream,
      handler: []eventHandler{handle},
    }
    ins.ctx, ins.cancel = context.WithCancel(ins.ctx)

    go func(ins *eventIns) {
      <- stream.Ready()
      // stream close
      defer stream.Close()
      defer delete(e.streams, method)

      for  {
        var raw []byte
        err := stream.RecvMsg(&raw)
        if err != nil {
          break
        }
        for _, h := range ins.handler {
          h(raw)
        }
      }
    }(ins)
  } else {
    ins.handler = append(ins.handler, handle)
  }

  return nil
}

func (e *Emitter) Off(method string, hds ...eventHandler) {
  e.onMu.Lock()
  defer e.onMu.Lock()

  ins, ok := e.streams[method]
  if !ok {
    return
  }
  if len(hds) == 0 {
    ins.cancel()
    return
  }

  nextHandlers := []eventHandler{}
  for _, existH := range ins.handler {
    rm := false
    for _, offH := range hds {
      if reflect.ValueOf(existH).Pointer() == reflect.ValueOf(offH).Pointer() {
        rm = true
      }
    }
    if !rm {
      nextHandlers = append(nextHandlers, existH)
    }
  }
  ins.handler = nextHandlers
}
