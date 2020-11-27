package cdpapp

import (
  "context"
  "github.com/mafredri/cdp"
  "github.com/mafredri/cdp/protocol/page"
  "github.com/mafredri/cdp/protocol/target"
  "github.com/mafredri/cdp/rpcc"
)

type Window struct {
  ctx context.Context
  conn *rpcc.Conn
  client *cdp.Client
  Config *WinConfig
  targetID target.ID
}

func newWindow(targetID target.ID, conn *rpcc.Conn, client *cdp.Client, conf *WinConfig) (*Window) {
  return &Window{
    conn:     conn,
    client:   client,
    Config:   conf,
    targetID: targetID,
    ctx: context.Background(),
  }
}

func (w *Window) Load(url string) error {
  _, err := w.client.Page.Navigate(w.ctx, &page.NavigateArgs{URL: url})
  return err
}

func (w *Window) Fullscreen() {

}
func (w *Window) Maximize() {}
func (w *Window) Minimize() {}
func (w *Window) SetIcon() {}
func (w *Window) SetTitle() {}