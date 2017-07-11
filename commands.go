package main

import (
	"encoding/json"
	"strings"

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

func DocumentWrapper(bot *tgbotapi.BotAPI, Msg *tgbotapi.Message) {
	//r, _ := bot.GetFileDirectURL(Msg.Voice.FileID)
	switch cachedAccounts[Msg.From.ID].GetCtx() {
	case 126:
		fileDS := strings.Split(Msg.Document.FileName, ".")
		if allowedVideoExt[fileDS[len(fileDS)-1]] {
			MsgConf := tgbotapi.NewMessage(Msg.Chat.ID, "Yooopta!")
			bot.Send(MsgConf)
		} else {
			//aExt, err := json.MarshalIndent(allowedVideoExt, "", "\t")
			//errh(err)
			//здеся блядь паника, сука, ебаный мэп!
			MsgConf := tgbotapi.NewMessage(Msg.Chat.ID, "Это чо бля?!/n")
			bot.Send(MsgConf)
		}
		//cachedAccounts[Msg.From.ID].SetCtx("0")
		return

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
	case "help":
		MsgConf = tgbotapi.NewMessage(Msg.Chat.ID, "World")
	case "me":
		r, _ := json.MarshalIndent(cachedAccounts[Msg.From.ID], "", "\t")
		MsgConf = tgbotapi.NewMessage(Msg.Chat.ID, string(r))
	case "uplv":
		cachedAccounts[Msg.From.ID].SetCtx("126")
		r := "Грузи видосик, епта!"
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
