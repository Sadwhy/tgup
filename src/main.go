package main

import (
  "fmt"
  "log"
  "os"
  "encoding/json"
  "strconv"
  "os/exec"
  "time"
  "github.com/joho/godotenv"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
  "github.com/Sadwhy/tgup/src/utils"
  "github.com/Sadwhy/tgup/src/media"
  "github.com/Sadwhy/tgup/src/dif"
)

var ( 
  bot *tgbotapi.BotAPI
  globalLink string
  failedList []string
  )

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	botId := os.Getenv("API_ID")
	if botId == "" {
		log.Fatal("API id not found in env")
	}

	botHash := os.Getenv("API_HASH")
	if botHash == "" {
		log.Fatal("API hash not found in env")
	}
	
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Bot token not found in env")
	}

  resultChan := make(chan *exec.Cmd)
  
  go func() {
      resultChan <- startBotAPI(botId, botHash)
  }()
  
  cmd := <-resultChan
  if cmd == nil {
    log.Fatal("Failed to start Telegram Bot API server.")
  }
  defer cleanUp(cmd)

  apiEndpoint := "http://localhost:8081/bot%s/%s"
  
	dif.Bot, err = tgbotapi.NewBotAPIWithAPIEndpoint(botToken, apiEndpoint)
	if err != nil {
		log.Panic(err)
	}

  log.Printf("Authorized on account %s", dif.Bot.Self.UserName)

	chatIdStr := utils.GetInput("Enter Chat ID: ")

  dif.ChatId, err = strconv.ParseInt(chatIdStr, 10, 64)
  if err != nil {
    log.Fatalf("Invalid chat Id: %s", err)
  }

  ForumIdStr := utils.GetInput("Enter Topic Id (optional): ")
  
  if ForumIdStr != "" {
    ForumId64, err := strconv.ParseInt(ForumIdStr, 10, 64)
    if err != nil {
      log.Fatalf("Invalid Topic Id: %s", err)
    }
    dif.ForumId = int(ForumId64)
  }

  links := getLinks()
  
  for i, link := range links {
    globalLink = link
    err := media.Download(link)
    if err != nil {
      fmt.Printf("Failed to process %s: %s", i, err)
      os.RemoveAll("download")
      failed()
      continue
    }
  }

	log.Println("Job completed!")
}


func getLinks() []string {
  if _, err := os.Stat("list.json"); err == nil {
    fmt.Println("List found. Downloading from it instead")
    data, err := os.ReadFile("list.json")
    if err != nil {
      link := utils.GetInput("Invalid list. Enter Url instead: ")
      return []string{link}
    }
    var links []string
    if err := json.Unmarshal(data, &links); err != nil {
      link := utils.GetInput("Invalid list. Enter Url instead: ")
      return []string{link}
    }
    return links
  }
  link := utils.GetInput("Enter file URL: ")
  return []string{link}
}

func startBotAPI(tgId string, tgHash string) *exec.Cmd {
  apiId := fmt.Sprintf("--api-id=%s", tgId)
  
  apiHash := fmt.Sprintf("--api-hash=%s", tgHash)
  
	cmd := exec.Command("./telegram-bot-api",
		apiId,
		apiHash,
		"--http-port=8081",
		"--local",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Starting Telegram Bot API server...")
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start Telegram Bot API server: %v", err)
		return nil
	}

	time.Sleep(3 * time.Second)

	log.Println("Telegram Bot API server started successfully!")
	return cmd
}

func cleanUp(cmd *exec.Cmd) {
  os.RemoveAll("download")
  
  if len(failedList) > 0 {
      filePath := "failed_list.json"
      data, _ := json.MarshalIndent(failedList, "", "  ")
      err := os.WriteFile(filePath, data, 0644)
      if err != nil {
          fmt.Println("Error writing file:", err)
      } else {
          fmt.Println("Saved failed links to", filePath)
          media.Upload(filePath)
      }
  }

  if cmd != nil {
      err := cmd.Process.Kill()
      if err != nil {
          fmt.Println("Failed to kill API server:", err)
      } else {
          fmt.Println("API server killed.")
      }
  }
}

func failed() {
  failedList = append(failedList, globalLink)
  failedText := fmt.Sprintf("Failed: %s", globalLink)
  msg := tgbotapi.NewMessage(dif.ChatId, failedText)
  if dif.ForumId != 0 {
    msg.MessageThreadID = dif.ForumId
  }
  _, err := dif.Bot.Send(msg)
  if err != nil {
    fmt.Println("m Failed to send error.")
  }
}
