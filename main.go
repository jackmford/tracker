package main

import (
  "encoding/json"
  "fmt"
  "github.com/getlantern/systray"
  "os"
  "os/user"
  "net/http"
  "time"
)

type Session struct {

  Category string `json:"category"`
  StartTime time.Time `json:"start_time"`
  EndTime time.Time `json:"end_time"`
  Duration time.Duration `json:"duration"`

}

type ActiveSession struct {
  Category string `json:"category"`
  StartTime time.Time `json:"start_time"`
}

func getFilePath(fileName string) string {
    usr, _ := user.Current()
    // This puts the files in your home folder (e.g., /Users/jack/history.json)
    return usr.HomeDir + "/" + fileName
}

func startSession(category string) error {
  if _, err := os.Stat(getFilePath("active.json")); err == nil {
    fmt.Println("A session is already running. Stop it before starting a new one.")
    return nil
  }
  active := ActiveSession{
    Category: category,
    StartTime: time.Now(),
  }
  data, err := json.MarshalIndent(active, "", " ")

  if err != nil {
    return err
  }

  return os.WriteFile(getFilePath("active.json"), data, 0644)
}

func stopSession() error {
  data, err := os.ReadFile(getFilePath("active.json"))
  if err != nil {
    return err
  }

  var active ActiveSession

  err = json.Unmarshal(data, &active)
  if err != nil {
    return err
  }

  session := Session{
    Category: active.Category,
    StartTime: active.StartTime,
    EndTime: time.Now(),
    Duration: time.Since(active.StartTime),
  }

  var history []Session
  historyData, err := os.ReadFile(getFilePath("history.json"))

  if err == nil {
    err = json.Unmarshal(historyData, &history)
    if err != nil {
      fmt.Println("Warning: Could not parse history.json as a list. Starting fresh list.")
      history = []Session{} 
    }
  }

  history = append(history, session)

  finalData, err := json.MarshalIndent(history, "", " ")
  if err != nil {
    return err
  }

  err = os.WriteFile(getFilePath("history.json"), finalData, 0644)
  if err != nil {
    return err
  }

  return os.Remove(getFilePath("active.json"))

}

func checkStatus() (string, error) {
  data, err := os.ReadFile(getFilePath("active.json"))
  if err != nil {
    fmt.Println("No active session running.")
    return "", err
  }

  var active ActiveSession
  json.Unmarshal(data, &active)

  elapsed := time.Since(active.StartTime).Round(time.Second)
  fmt.Printf("Currently tracking: %s\n", active.Category)
  fmt.Printf("Time elapsed: %s\n", elapsed)

  return active.Category, nil
}

func startServer() {
  fs := http.FileServer(http.Dir("."))
  http.Handle("/", fs)

  fmt.Println("📊 Dashboard available at: http://localhost:8080")
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    fmt.Printf("Error starting server: %v\n", err)
  }
}

func handleCli() {
  if len(os.Args) < 2 {
    fmt.Println("Usage tracker [start <category> | stop | status]")
    return
  }

  command := os.Args[1]

  switch command {
  case "start":
    if len(os.Args) < 3 {
      fmt.Println("Please provide a category (e.g. start programming)")
    }
    err := startSession(os.Args[2])
    if err != nil {
      fmt.Printf("Error starting session %v\n", err)
    } else {
      fmt.Printf("Started tracking: %s\n", os.Args[2])
    }
  case "stop": 
    err := stopSession()

    if err != nil {
      fmt.Printf("Error stopping session: %v\n", err)
    } else {
      fmt.Println("Session stopped and saved to history.")
    }
  case "status":
    checkStatus()
  case "serve":
    startServer()
  default:
    fmt.Println("Unknown command. Use: start, stop, status.")
  }
}

func onReady() {
  _, err := os.Stat(getFilePath("active.json"))
  if err == nil {
    cat, _ := checkStatus()
    if cat == "" {cat = "Active"}
    title := fmt.Sprintf("🚀 %s", cat)
    systray.SetTitle(title)
  } else{
    systray.SetTitle("🕒")
    systray.SetTooltip("Time Tracker")
  }

  mProgramming := systray.AddMenuItem("Start: Programming", "Start tracking programming")
	mReading := systray.AddMenuItem("Start: Reading", "Start tracking reading")
	systray.AddSeparator()
	mStop := systray.AddMenuItem("Stop Tracking", "Stop the current session")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

  go func() {
		for {
			select {
			case <-mProgramming.ClickedCh:
				startSession("programming")
				systray.SetTitle("🚀 programming")
			case <-mReading.ClickedCh:
				startSession("reading")
				systray.SetTitle("🚀 reading")
			case <-mStop.ClickedCh:
				stopSession()
				systray.SetTitle("🕒 Idle")
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {}

func main() {
  if len(os.Args) > 1 {
		handleCli() // Move your old switch statement here
		return
	}

  systray.Run(onReady, onExit)
}
