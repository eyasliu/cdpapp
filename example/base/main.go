package main

import (
  "cdpapp"
)

func main() {
  app, err := cdpapp.NewApp(&cdpapp.Config{
    Debug: true,
    Title:  "My Test Gui app",
    ExecutablePath: "",
    Flags: cdpapp.Flags{},
    // BackgroundColor: "#66ccff",
    Width:  800,
    Height: 600,
    Top: 0,
    Left: 0,
    URL: "http://www.baidu.com",
  })
  if err != nil {
    panic(err)
  }

  _ = app.Load("https://github.com")

  
  win, err := app.CreateWindow(&cdpapp.WinConfig{
    Title:           "",
    URL:             "https://www.baidu.com",
    BackgroundColor: "",
    Width:           1024,
    Height:          768,
    Top:             0,
    Left:            0,
  })
  win.Maximize()
  // win.Load("https://www.baidu.com")

  if err := app.Wait(); err != nil {
    panic(err)
  }
}
