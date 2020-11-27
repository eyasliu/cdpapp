package cdpapp

import (
  "context"
  "github.com/mafredri/cdp"
  "github.com/mafredri/cdp/protocol/target"
  "github.com/mafredri/cdp/rpcc"
  "github.com/mafredri/cdp/session"
)

type CDPSession struct {
  conf *Config
  ctx context.Context
  client *cdp.Client
  manager *session.Manager
  conn *rpcc.Conn
  targets map[string]*target.GetTargetsReply
  pageSession map[target.ID]*Window
  events *Emitter
}

func NewSession(ctx context.Context, client *cdp.Client, conn *rpcc.Conn) (*CDPSession, error) {
  manager, err := session.NewManager(client)
  if err != nil {
    return nil, err
  }

  return &CDPSession{
    ctx:         ctx,
    client:         client,
    manager:     manager,
    conn:        conn,
    targets:     make(map[string]*target.GetTargetsReply),
    pageSession: make(map[target.ID]*Window),
    events:      NewEmitter(ctx, conn),
  }, nil
}

func (s *CDPSession) targetCreated(targetId target.ID) (*cdp.Client, error) {
  conn, err := s.manager.Dial(s.ctx, targetId)
  if err != nil {
    return nil, err
  }
  client := cdp.NewClient(conn)

  s.pageSession[targetId] = newWindow(targetId, conn, client, &WinConfig{})

  return client, nil
}

func (s *CDPSession) GetMainPage() *cdp.Client {
  for _, s := range s.pageSession {
    return s.client
  }
  return nil
}

func (s *CDPSession) Get(targetId target.ID) (*cdp.Client, bool) {
  c, ok := s.pageSession[targetId]
  return c.client, ok
}

func (s *CDPSession) GetConn(targetId target.ID) (*rpcc.Conn, bool) {
  c, ok := s.pageSession[targetId]
  return c.conn, ok
}