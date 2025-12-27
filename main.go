package main

import (
	"encoding/json"
  "fmt"
	"os"
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

func startSession(category string) error {
  if _, err := os.Stat("active.json"); err == nil {
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

  return os.WriteFile("active.json", data, 0644)
}

func stopSession() error {
  data, err := os.ReadFile("active.json")
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
  historyData, err := os.ReadFile("history.json")

  if err == nil {
    json.Unmarshal(historyData, &history)

  }

  history = append(history, session)

  finalData, err := json.MarshalIndent(session, "", " ")
  if err != nil {
    return err
  }

  err = os.WriteFile("history.json", finalData, 0644)
  if err != nil {
    return err
  }

  return os.Remove("active.json")

}

func checkStatus() {
	data, err := os.ReadFile("active.json")
	if err != nil {
		fmt.Println("No active session running.")
		return
	}

	var active ActiveSession
	json.Unmarshal(data, &active)

	elapsed := time.Since(active.StartTime).Round(time.Second)
	fmt.Printf("Currently tracking: %s\n", active.Category)
	fmt.Printf("Time elapsed: %s\n", elapsed)
}

func main() {
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
  default:
    fmt.Println("Unknown command. Use: start, stop, status.")
  }


}
