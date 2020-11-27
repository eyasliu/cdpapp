package cdpapp

import (
  "os"
  "time"
)

var (
  DefaultWidth           = 800
  DefaultHeight          = 600
  DefaultBackgroundColor = "#2b2e3b"

  wsURLReadTimeout = 20 * time.Second
)

type empty struct{}

type Config struct {
  Title           string
  URL             string
  BackgroundColor string
  Width           int
  Height          int
  Top             int
  Left            int
  Envs []string
  Bootstrap string
  UserDataDir string
  ExecutablePath string
  Flags       Flags
  Debug bool
}

func (conf *Config) setDefault() {
  if conf.Title == "" {
    conf.Title = "GCDP App"
  }
  if conf.Width == 0 {
    conf.Width = DefaultWidth
  }
  if conf.Height == 0 {
    conf.Height = DefaultHeight
  }
  if conf.BackgroundColor == "" {
    conf.BackgroundColor = DefaultBackgroundColor
  }
  if conf.UserDataDir == "" {
    dir, _ := os.UserHomeDir()
    conf.UserDataDir = dir + "/.gcdp/apps/" + md5Sign([]byte(conf.Title))
  }
  if conf.URL == "" {
    conf.URL = "about:blank"
  }
  conf.ExecutablePath = findChromeExecPath(conf.ExecutablePath)
}

type WinConfig struct {
  Title           string
  URL             string
  BackgroundColor string
  Width           int
  Height          int
  Top             int
  Left            int
}
