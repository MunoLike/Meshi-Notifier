package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var DiscordSession *discordgo.Session

func notifyMeshi(c echo.Context) error {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("%v", err)
	}
	now := time.Now().In(jst)

	hour := now.Hour()
	if 3 <= hour && hour <= 10 {
		// 朝
		DiscordSession.ChannelMessageSend(os.Getenv("CHANNEL_ID"), "朝飯の時間よ～～🍚")
	} else if 11 <= hour && hour <= 15 {
		// 昼
		DiscordSession.ChannelMessageSend(os.Getenv("CHANNEL_ID"), "昼飯～～\u2600")
	} else {
		// 夜
		DiscordSession.ChannelMessageSend(os.Getenv("CHANNEL_ID"), "晩飯だからすぐに来い😁")
	}

	return c.String(http.StatusOK, "Hello, World")
}

func main() {
	if os.Getenv("ISDBG") == "True" {
		err := godotenv.Load("./.env")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Env Loaded from the file.")
	}

	// Init discord api
	ds, err := discordgo.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("ログインに失敗しました．")
		fmt.Println(err)
		return
	}

	ds.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		DiscordSession = s
		DiscordSession.UpdateGameStatus(0, "Minecraft")
		fmt.Println("Bot is ready")
	})

	ds.Identify.Intents = discordgo.IntentGuildMessages

	err = ds.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer ds.Close()

	// Start REST API
	e := echo.New()

	e.GET("/call", notifyMeshi)

	e.Logger.Fatal(e.Start(":" + os.Getenv(("PORT"))))
}
