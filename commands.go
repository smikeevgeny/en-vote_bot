package main

import (
	"encoding/json"
	"strings"

	"fmt"

	"unsafe"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func VoiceWrapper(bot *tgbotapi.BotAPI, Msg *tgbotapi.Message) {
	//r, _ := bot.GetFileDirectURL(Msg.Voice.FileID)
	switch cachedAccounts[Msg.From.ID].GetCtx() {
	case 127:
		MsgConf := tgbotapi.NewVoiceShare(Msg.Chat.ID, Msg.Voice.FileID)
		bot.Send(MsgConf)
		cachedAccounts[Msg.From.ID].SetCtx("0")

	default:
		MsgConf := tgbotapi.NewMessage(Msg.Chat.ID, "Bye!")
		bot.Send(MsgConf)
	}

}
func cmdWrapper(Msg *tgbotapi.Message) (MsgConf tgbotapi.MessageConfig) {
	switch Msg.Command() {
	case "auth":
		if len(strings.Split(Msg.CommandArguments(), " ")) < 2 {
			MsgConf = tgbotapi.NewMessage(Msg.Chat.ID, "Нах!")
			return MsgConf
		}
		creds := strings.Split(Msg.CommandArguments(), " ")
		MsgConf = authenticateUser(creds, Msg.From, Msg)

	case "hello":
		MsgConf = tgbotapi.NewMessage(Msg.Chat.ID, "World")
	case "rem":
		MsgConf = tgbotapi.NewMessage(Msg.Chat.ID, "World")
		cachedAccounts[Msg.From.ID].SetCtx("127")
	case "lens":
		r := fmt.Sprintf("%d", unsafe.Sizeof(cachedAccounts))
		MsgConf = tgbotapi.NewMessage(Msg.Chat.ID, r)
	case "me":
		r, _ := json.MarshalIndent(cachedAccounts[Msg.From.ID], "", "\t")
		MsgConf = tgbotapi.NewMessage(Msg.Chat.ID, string(r))
	case "setctx":
		cachedAccounts[Msg.From.ID].SetCtx(Msg.CommandArguments())
		r, _ := json.MarshalIndent(cachedAccounts[Msg.From.ID], "", "\t")
		MsgConf = tgbotapi.NewMessage(Msg.Chat.ID, string(r))

	default:
		MsgConf = tgbotapi.NewMessage(Msg.Chat.ID, "NIE")
	}

	return
}

func authenticateUser(creds []string, u *tgbotapi.User, rMsg *tgbotapi.Message) (msg tgbotapi.MessageConfig) {
	a := Account{Login: creds[0], Password: creds[1]}
	has, err := engine.Get(&a)
	errh(err)
	if has {
		a.Login = " "
		a.Password = " "
		a.TelegramID = u.ID
		a.Name = u.FirstName + " " + u.LastName
		a.UserName = u.UserName
		_, err = engine.Update(&a)
		errh(err)
		msg = tgbotapi.NewMessage(rMsg.Chat.ID, "User successfuly authenticated")
	} else {
		msg = tgbotapi.NewMessage(rMsg.Chat.ID, rMsg.Text)
		msg.ReplyToMessageID = rMsg.MessageID
	}
	return msg
}
