package cdpapp

import (
  "context"
  "encoding/json"
  "errors"
  "github.com/mafredri/cdp"
  "github.com/mafredri/cdp/protocol/target"
  "github.com/mafredri/cdp/rpcc"
  "os"
  "os/exec"
  "time"
)

func (a *App) run() error {
  a.cmd = exec.CommandContext(a.ctx, a.Config.ExecutablePath, a.flags.Args()...)

  // set env
  a.cmd.Env = os.Environ()
  if len(a.Config.Envs) > 0 {
    a.cmd.Env = append(a.cmd.Env, a.Config.Envs...)
  }

  // ws url will output to stderr
  stderr, err := a.cmd.StderrPipe()
  if err != nil {
    return err
  }

  if err := a.cmd.Start(); err != nil {
    return err
  }
  a.wsURL, err = getWsURLFromPipeWithTimeout(stderr, wsURLReadTimeout)
  if err != nil {
    return err
  }
  err = a.setupWsConnection()
  if err != nil {
    return err
  }

  err = a.setupCdpSession()
  if err != nil {
    return err
  }
  return nil
}

func (a *App) setupWsConnection() error {
  initConn, err := rpcc.DialContext(a.ctx, a.wsURL)
  if err != nil {
    return err
  }

  a.initClient = cdp.NewClient(initConn)


  session, err := NewSession(a.ctx, a.initClient, initConn)
  if err != nil {
    return err
  }
  a.session = session

  return nil
}

func (a *App) setupCdpSession() error {
  ctx, _ := context.WithTimeout(a.ctx, 10 * time.Second)
  errCh := make(chan error, 0)
  err := a.session.events.On("Target.targetCreated", func(msg json.RawMessage) {
    data := &target.CreatedReply{}
    _ = json.Unmarshal(msg, data)
    if data.TargetInfo.Type == "page" {
      client, err := a.session.targetCreated(data.TargetInfo.TargetID)
      if err != nil {
        errCh <- err
        return
      }
      err = client.Page.Enable(a.ctx)
      if err != nil {
        errCh <- err
        return
      }
      err = client.Runtime.Enable(a.ctx)
      if err != nil {
        errCh <- err
        return
      }

      err = client.DOM.Enable(a.ctx)
      if err != nil {
        errCh <- err
        return
      }

      errCh <- nil

    }
  })
  if err != nil {
    return err
  }

  err = a.initClient.Target.SetDiscoverTargets(a.ctx, &target.SetDiscoverTargetsArgs{Discover: true})
  if err != nil {
    return err
  }

  select {
  case <- ctx.Done():
    return errors.New("connect session timeout")
  case err := <- errCh:
    if err != nil {
      return err
    }
  }

  return nil
}
