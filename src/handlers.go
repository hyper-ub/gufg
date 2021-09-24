package main

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

import (
	"codeberg.org/massivebox/botsarchive-api"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strconv"
	"time"
)

func (pt passthrough) handleStartMessage(update tgbotapi.Update) {

	data, err := databaseGetUserData(pt.Db, update.Message.Chat.ID)
	if err != nil {
		pt.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "We had a problem registering you to the database.\nPlease try again -> /start."))
		fmt.Println(err.Error())
		return
	}
	if data.ID == 0 {
		databaseRegisterUser(pt.Db, userData{
			UserID:           update.Message.Chat.ID,
			Username:         update.Message.Chat.UserName,
			RegistrationTime: time.Now().Unix(),
		})
	}
	if data.Frames != "" {
		pt.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "🗑 GIF discarded."))
		data.Frames = ""
		databaseUpdateUserData(pt.Db, data)
	}

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Create GIF ➕", "/create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⭐️ Rate", "/vote"),
			tgbotapi.NewInlineKeyboardButtonURL("Source 🖥", pt.Env.Source),

		),
	)

	var messageText = "🖼 <b>Welcome to GIF Creator,</b> the bot to easily create GIFs from images.\nClick \"➕ Create GIF ➕\" to start.\n\n🤖 <a href=\""+pt.Env.DevLink+"\">Developer</a>\n📣 @"+pt.Env.DevChannel+"\n📎 <a href=\""+pt.Env.Privacy+"\">Privacy</a> - <a href=\""+pt.Env.Terms+"\">ToS</a>"

	resp := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	resp.DisableWebPagePreview = true
	resp.ParseMode = "HTML"
	resp.ReplyMarkup = keyboard
	pt.Bot.Send(resp)

	pt.matomoAnalytics("start", update.Message.Chat.ID)

}

func (pt passthrough) handleIncomingPhoto(update tgbotapi.Update) {

	data, err := databaseGetUserData(pt.Db, update.Message.Chat.ID)
	if data.Page != "create" {
		pt.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You're not creating a GIF at the moment. Click the \"➕ Create GIF ➕\" in the /start menu to create a GIF."))
		return
	}

	var length int
	if update.Message.Caption == "" {
		length = 100
	}else{
		length, err = strconv.Atoi(update.Message.Caption)
		if err != nil {
			pt.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid duration value! Only include a single number in the photo caption, representing how much time it should stay on screen, in 100ths of a second."))
			return
		}
	}

	if len(update.Message.Photo) < 1 {
		pt.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error processing the image size. Please try again."))
		return
	}

	err = savePhoto(update.Message.Photo[len(update.Message.Photo)-1].FileID, length, update.Message.Chat.ID, pt.Db)
	if err != nil {
		if err.Error() == "free" {
			pt.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are exceeding the limit of 10 seconds. You can lift this limit by following the instructions in the \"⭐️ Rate\" section of the /start menu."))
			return
		}
		if err.Error() == "limits" {
			pt.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are exceeding the limit of 150 different frames per GIF. This limit can not be lifted."))
			return
		}
		pt.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error saving the image. Please try again."))
		return
	}

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎥 Render", "/render"),
			tgbotapi.NewInlineKeyboardButtonData("🗑 Discard", "/start"),

		),
	)

	resp := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ <b>Photo added.</b> You can send more photos to add them to the GIF, or click one of the buttons below.")
	resp.ParseMode = "HTML"
	resp.ReplyMarkup = keyboard
	pt.Bot.Send(resp)

	pt.matomoAnalytics("incoming_photo", update.Message.Chat.ID)

}

func (pt passthrough) handleUnhandledMessage(update tgbotapi.Update) {
	pt.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unrecognized command. Press /start to continue."))
	pt.matomoAnalytics("unhandled", update.Message.Chat.ID)
}


func (pt passthrough) handleBackCallback(update tgbotapi.Update) {

	data, err := databaseGetUserData(pt.Db, update.CallbackQuery.Message.Chat.ID)
	if err == nil {
		if data.Frames != "" {
			pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "🗑 GIF discarded."))
			data.Frames = ""
			databaseUpdateUserData(pt.Db, data)
		}
	}

	pt.Bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Create GIF ➕", "/create"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⭐️ Rate", "/vote"),
			tgbotapi.NewInlineKeyboardButtonURL("Source 🖥", "https://codeberg.org/massivebox/gifcreator"),
		),
	)

	var messageText = "🖼 <b>Welcome to GIF Creator,</b> the bot to easily create GIFs from images.\nClick \"➕ Create GIF ➕\" to start.\n\n🤖 <a href=\""+pt.Env.DevLink+"\">Developer</a>\n📣 @"+pt.Env.DevChannel+"\n📎 <a href=\""+pt.Env.Privacy+"\">Privacy</a> - <a href=\""+pt.Env.Terms+"\">ToS</a>"
	resp := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, messageText)
	resp.DisableWebPagePreview = true
	resp.ParseMode = "HTML"
	resp.ReplyMarkup = &keyboard
	pt.Bot.Send(resp)

	pt.matomoAnalytics("back", update.CallbackQuery.Message.Chat.ID)

}

func (pt passthrough) handleCreateCallback(update tgbotapi.Update) {

	pt.Bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

	data, err := databaseGetUserData(pt.Db, update.CallbackQuery.Message.Chat.ID)
	if err == nil {
		data.Page = "create"
		databaseUpdateUserData(pt.Db, data)
	}

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "/start"),
		),
	)

	var messageText = "☑️ <b>You're now creating a GIF.</b>\nSend me any image, <i>with the duration it should have (in 100ths of a second) written in the caption.</i>\nFor example, if you want an image to stay on screen for 2 seconds, write <code>200</code> in its caption. <i>If the caption is empty, it will stay on screen one second.</i>"
	resp := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, messageText)
	resp.DisableWebPagePreview = true
	resp.ParseMode = "HTML"
	resp.ReplyMarkup = &keyboard
	pt.Bot.Send(resp)

	pt.matomoAnalytics("create", update.CallbackQuery.Message.Chat.ID)

}

func (pt passthrough) handleVoteCallback(update tgbotapi.Update) {

	data, err := databaseGetUserData(pt.Db, update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Error, please try again"))
		return
	}
	if data.HasVoted {
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "✅ You've already voted!"))
		return
	}

	botInfo, err := baapi.GetBotInfo(pt.Bot.Self.UserName)
	if err != nil {
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Error getting BotsArchive data, please try again"))
		return
	}

	pt.Bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("⭐️ Rate", botInfo.ChannelLink),
			tgbotapi.NewInlineKeyboardButtonData("☑️ Check Rating", "/check_vote"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Back", "/start"),
		),
	)

	var messageText = "➕ Do you want to <b>make GIFs that are longer than 10 seconds?</b>\nYou can unlock that feature by rating our bot 5 stars on BotsArchive."
	resp := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, messageText)
	resp.ParseMode = "HTML"
	resp.ReplyMarkup = &keyboard
	pt.Bot.Send(resp)

	pt.matomoAnalytics("vote", update.CallbackQuery.Message.Chat.ID)

}

func (pt passthrough) handleCheckVoteCallback(update tgbotapi.Update) {

	hasVoted, rating, err := baapi.GetUserRatingByBotUsername(pt.Bot.Self.UserName, update.CallbackQuery.From.ID)
	if err != nil {
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Error getting BotsArchive data, please try again"))
		return
	}
	if !hasVoted {
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "You haven't voted yet, or you have vote privacy on."))
		return
	}
	if rating != 5 {
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "You haven't selected 5 stars."))
		return
	}

	pt.Bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

	data, err := databaseGetUserData(pt.Db, update.CallbackQuery.Message.Chat.ID)
	if err == nil {
		data.HasVoted = true
		databaseUpdateUserData(pt.Db, data)
	}

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Home", "/start"),
		),
	)

	var messageText = "✅ <b>Rating received!</b> Your limits have been lifted.\n<b>Don't change your rating!</b> You'll be banned from using the bot if you do so."
	resp := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, messageText)
	resp.ParseMode = "HTML"
	resp.ReplyMarkup = &keyboard
	pt.Bot.Send(resp)

	pt.matomoAnalytics("has_voted", update.CallbackQuery.Message.Chat.ID)

}

func (pt passthrough) handleRenderCallback(update tgbotapi.Update) {

	data, err := databaseGetUserData(pt.Db, update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Error getting database data, please try again"))
		return
	}
	if data.Page == "rendering" {
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "You already have a render in progress! Wait for it to finish."))
		return
	}
	if data.Frames == "" {
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "You don't have any frame saved."))
		return
	}

	pt.Bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

	var messageText = "🚛 <b>We're rendering your GIF</b>, please be patient, it might take some time."
	resp := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, messageText)
	resp.ParseMode = "HTML"
	pt.Bot.Send(resp)

	data.Page = "rendering"
	databaseUpdateUserData(pt.Db, data)

	err = pt.renderGIF(data)
	if err != nil {
		data.Page = ""
		databaseUpdateUserData(pt.Db, data)
		pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Error rendering GIF: " + err.Error()))
		return
	}
	data.Frames = ""
	data.Page = ""
	databaseUpdateUserData(pt.Db, data)

	messageText = "🚛 <b>Gif rendered successfully!</b> we will send it to you soon, please wait."
	resp = tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, messageText)
	resp.ParseMode = "HTML"
	pt.Bot.Send(resp)

	gifPath := fmt.Sprintf("./tmp/%d.gif", update.CallbackQuery.Message.Chat.ID)

	resp2 := tgbotapi.NewAnimationUpload(update.CallbackQuery.Message.Chat.ID, gifPath)
	resp2.Caption = "🎥 Made with @" + pt.Bot.Self.UserName
	_, err = pt.Bot.Send(resp2)
	if err != nil {
		messageText = "🛑 Error sending GIF, try again later.\nError: " + err.Error()
		resp2 := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messageText)
		resp2.ParseMode = "HTML"
		pt.Bot.Send(resp2)
	}

	os.Remove(gifPath)

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Home", "/start"),
		),
	)

	messageText = "⬆️ <b>GIF uploaded.</b> check the message above.\nIf it didn't work, contact the developer or open an issue on the repository."
	resp3 := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, messageText)
	resp3.ParseMode = "HTML"
	resp3.ReplyMarkup = &keyboard
	pt.Bot.Send(resp3)

	pt.matomoAnalytics("rendered", update.CallbackQuery.Message.Chat.ID)

}

func (pt passthrough) handleUnhandledCallback(update tgbotapi.Update) {
	pt.Bot.Send(tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "Internal error. Please try again."))
	pt.matomoAnalytics("unhandled_callback", update.CallbackQuery.Message.Chat.ID)
}