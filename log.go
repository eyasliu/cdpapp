package cdpapp

import (
  lg "log"
)

type logger struct {
  isDebug bool
}

var log = &logger{}

func (l logger) Log(v ...interface{}) {
  if l.isDebug {
    lg.Print(v...)
  }
}

func (l logger) Logf(f string, v ...interface{}) {
  if l.isDebug {
    lg.Printf(f, v...)
  }
}

