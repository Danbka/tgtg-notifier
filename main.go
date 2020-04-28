package main

import (
	"fmt"
	"github.com/Danbka/tgtg_notifier/tgtg"
	"github.com/joho/godotenv"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main()  {
	godotenv.Load(".env")

	stores := strings.Split(os.Getenv("TGTG_STORES"), ",")

	var tgtgChecker tgtg.Client

	tgtgChecker.SetUserId(os.Getenv("TGTG_USER_ID"))
	tgtgChecker.SetAccessToken(os.Getenv("TGTG_ACCESS_TOKEN"))
	tgtgChecker.SetRefreshToken(os.Getenv("TGTG_REFRESH_TOKEN"))

	for _, storeId := range stores {
		res, err := tgtgChecker.HasAvailableItem(storeId)

		if err != nil {
			fmt.Println(err)
		}

		if res {
			message := "Available items in shop " + storeId
			http.Get("https://api.telegram.org/bot"+os.Getenv("TG_BOT")+"/sendMessage?chat_id="+os.Getenv("CHAT_ID")+"&text="+url.QueryEscape(message))
		}
	}
}
