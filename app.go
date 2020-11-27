package cdpapp

import (
  "context"
  "fmt"
  "github.com/mafredri/cdp"
  "github.com/mafredri/cdp/protocol/page"
  "github.com/mafredri/cdp/protocol/runtime"
  "net/url"
  "os/exec"
)

type App struct {
  Config *Config
  flags Flags
  ctx    context.Context
  cancel context.CancelFunc
  cmd *exec.Cmd
  wsURL string
  session *CDPSession
  initClient *cdp.Client
  pendingWindow map[string]*WinConfig // creating windows
}

func NewApp(conf *Config) (*App, error) {
  conf.setDefault()
  if conf.Debug {
    log.isDebug = true
  }

  flags := DefaultExecFlags
  flags.Assign(conf.Flags)
  bootstrapHTML := fmt.Sprintf("data:text/html," +
    "<title>%s</title><style>html{background:%s;}</style>" +
    "<h1>Welcome To %s</h1>", url.PathEscape(conf.Title), url.PathEscape(conf.BackgroundColor), url.PathEscape(conf.Title))
  flags.Set("app", bootstrapHTML)
  flags.Set("window-size", fmt.Sprintf("%d,%d", conf.Width, conf.Height))
  if conf.Top != 0 || conf.Left != 0 {
    flags.Set("window-position", fmt.Sprintf("%d,%d", conf.Top, conf.Left))
  }
  flags.Set("user-data-dir", conf.UserDataDir)

  log.Logf("flags: %+v", flags)
  ctx, cancel := context.WithCancel(context.Background())

  app := &App{
    ctx:    ctx,
    cancel: cancel,
    flags: flags,
    Config: conf,
    pendingWindow: make(map[string]*WinConfig),
  }

  err := app.run()

  return app, err
}

func (a *App) Wait() error {
  defer a.Destroy()
  return a.cmd.Wait()
}

// Destroy TODO
func (a *App) Destroy() {
  a.cmd.Process.Kill()
}

func (a *App) MainWindow() *cdp.Client {
  return a.session.GetMainPage()
}

// TODO wait for window
func (a *App) CreateWindow(conf *WinConfig) (*Window, error) {
  seqno := RandStr(6)
  _, err := a.MainWindow().Runtime.Evaluate(a.ctx, runtime.NewEvaluateArgs(fmt.Sprintf(`window.open('%s', '', 'width=%d,height=%d,name=%s')`, conf.URL, conf.Width, conf.Height, seqno)))
  return nil, err
}

func (a *App) Load(url string) error {
  _, err := a.MainWindow().Page.Navigate(a.ctx, &page.NavigateArgs{URL: url})
  return err
}