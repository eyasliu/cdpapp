package cdpapp

type Flags map[string]string

var DefaultExecFlags = Flags{
  "disable-background-networking":          "",
  "enable-features":                        "NetworkService,NetworkServiceInProcess",
  "disable-background-timer-throttling":    "",
  "disable-backgrounding-occluded-windows": "",
  "disable-breakpad":                       "",
  "disable-client-side-phishing-detection": "",
  "disable-default-apps":                   "",
  "disable-dev-shm-usage":                  "",
  "disable-extensions":                     "",
  "disable-features":                       "site-per-process,TranslateUI,BlinkGenPropertyTrees",
  "disable-hang-monitor":                   "",
  "disable-ipc-flooding-protection":        "",
  "disable-popup-blocking":                 "",
  "disable-prompt-on-repost":               "",
  "disable-renderer-backgrounding":         "",
  "disable-sync":                           "",
  "force-color-profile":                    "srgb",
  "metrics-recording-only":                 "",
  "safebrowsing-disable-auto-update":       "",
  "enable-automation":                      "",
  "password-store":                         "basic",
  "use-mock-keychain":                      "",
  "no-sandbox":                             "",
  "remote-debugging-port":                  "0",
}

func (e Flags) Args() []string {
  s := make([]string, len(e))
  i := 0
  for k, v := range e {
    s[i] = "--" + k
    if v != "" {
      s[i] += "=" + v
    }
    i++
  }
  return s
}

func (e Flags) Assign(f Flags) {
  for k, v := range f {
    e[k] = v
  }
}

func (e Flags) Set(k string, v ...string) {
  if len(v) == 0 {
    e[k] = ""
  } else {
    e[k] = v[0]
  }
}
