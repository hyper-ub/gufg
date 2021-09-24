package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/*

	GifCreator - Telegram Bot to create GIFs from a series of images.
	Copyright (C) 2021  MassiveBox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

type passthrough struct {
	Bot *tgbotapi.BotAPI
	Db *sql.DB
	Env env
}

type env struct {
	DevLink string
	DevChannel string
	Source     string
	Privacy    string
	Terms string
	MatomoEnabled bool
	MatomoSiteID string
	MatomoHost string
}

func main() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Authorized on", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	db, err := databaseConnect()
	if err != nil {
		panic(err)
	}
	createTable(db)
	pt := passthrough{
		Bot: bot,
		Db:  db,
		Env: env{
			DevLink:       os.Getenv("DEV_LINK"),
			DevChannel:    os.Getenv("DEV_CHANNEL"),
			Source:        os.Getenv("SOURCE"),
			Privacy:       os.Getenv("PRIVACY"),
			Terms:         os.Getenv("TERMS"),
			MatomoEnabled: os.Getenv("MATOMO_ENABLED") == "1",
			MatomoSiteID:  os.Getenv("MATOMO_SITEID"),
			MatomoHost:    os.Getenv("MATOMO_HOST"),
		},
	}

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		go pt.dispatchUpdate(update)
	}

}

func (pt passthrough) dispatchUpdate(update tgbotapi.Update) {

	if update.Message != nil {
		if update.Message.Photo != nil {
			pt.handleIncomingPhoto(update)
		}else if strings.Contains(update.Message.Text, "/start") {
			pt.handleStartMessage(update)
		}else{
			pt.handleUnhandledMessage(update)
		}
	}
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Data == "/start" {
			pt.handleBackCallback(update)
		}else if update.CallbackQuery.Data == "/create" {
			pt.handleCreateCallback(update)
		}else if update.CallbackQuery.Data == "/vote" {
			pt.handleVoteCallback(update)
		}else if update.CallbackQuery.Data == "/check_vote" {
			pt.handleCheckVoteCallback(update)
		}else if update.CallbackQuery.Data == "/render" {
			pt.handleRenderCallback(update)
		}else{
			pt.handleUnhandledCallback(update)
		}
	}

}

func (pt passthrough) matomoAnalytics(actionName string, uID int64) {

	if pt.Env.MatomoEnabled {

		client := &http.Client{}
		payload := url.Values{}
		payload.Add("idsite", pt.Env.MatomoSiteID)
		payload.Add("rec", pt.Env.MatomoSiteID)
		payload.Add("url", "https://bots.com/gifcreator/"+actionName)
		payload.Add("uid", fmt.Sprintf("%d", uID))

		r, err := http.NewRequest("GET", "https://"+pt.Env.MatomoHost+"/matomo.php?"+payload.Encode(), nil)
		if err != nil {
			return
		}
		client.Do(r)

	}

}