package main

import (
	"encoding/json"
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

func main() {

}
