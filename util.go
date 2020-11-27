package cdpapp

import (
  "bufio"
  "bytes"
  "crypto/md5"
  "errors"
  "fmt"
  "io"
  "math/rand"
  "os/exec"
  "time"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStr 生成随机字符串
func RandStr(n int) string {
  b := make([]rune, n)
  for i := range b {
    b[i] = letterRunes[rand.Intn(len(letterRunes))]
  }
  return string(b)
}


// getWsURLFromOutput read ws url from io.pipe
func getWsURLFromOutput(rc io.ReadCloser) (wsURL string, _ error) {

  prefix := []byte("DevTools listening on")
  var accumulated bytes.Buffer
  bufr := bufio.NewReader(rc)
readLoop:
  for {
    line, err := bufr.ReadBytes('\n')
    if err != nil {
      return "", fmt.Errorf("chrome failed to start:\n%s",
        accumulated.Bytes())
    }

    if bytes.HasPrefix(line, prefix) {
      line = line[len(prefix):]
      // use TrimSpace, to also remove \r on Windows
      line = bytes.TrimSpace(line)
      wsURL = string(line)
      break readLoop
    }
    accumulated.Write(line)
  }
  return wsURL, nil
}

func getWsURLFromPipeWithTimeout(rc io.ReadCloser, to time.Duration) (string, error) {
  wsUrlChan := make(chan error, 0)
  var wsURL string
  var err error
  go func() {
    wsURL, err = getWsURLFromOutput(rc)
    wsUrlChan <- err
  }()
  select {
  case <-wsUrlChan:
  case <-time.After(wsURLReadTimeout):
    err = errors.New("connect chrome devtool protocol websocket timeout")
  }
  if err != nil {
    return "", err
  }

  return wsURL, nil

}

func md5Sign(src []byte) string {
  has := md5.Sum(src)
  return fmt.Sprintf("%x", has)
}

var preDefinedExecpath = []string{
  // Unix-like
  "headless_shell",
  "headless-shell",
  "chromium",
  "chromium-browser",
  "google-chrome",
  "google-chrome-stable",
  "google-chrome-beta",
  "google-chrome-unstable",
  "/usr/bin/google-chrome",

  // Windows
  "chrome",
  "chrome.exe", // in case PATHEXT is misconfigured
  `C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
  `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
  `D:\Program Files (x86)\360Chrome\Chrome\Application\360chrome.exe`,
  `D:\Program Files (x86)\360Chrome\Chrome\Application\360se.exe`,
  `C:\Program Files\Google\Chrome\Application\chrome.exe`,

  // Mac
  "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
}

// findChromeExecPath 查找当前机器的chrome安装位置
func findChromeExecPath(execPath string) string {
  checkPaths := make([]string, 0, len(preDefinedExecpath) + 1)
  checkPaths = append(checkPaths, execPath)
  checkPaths = append(checkPaths, preDefinedExecpath...)

  for _, path := range checkPaths {
    if path == "" {
      continue
    }
    found, err := exec.LookPath(path)
    if err == nil {
      return found
    }
  }
  return ""
}
