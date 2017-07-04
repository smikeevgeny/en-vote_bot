package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

type Configuration struct {
	TgAPIToken string `json:"TgAPIToken"`
	repo_token string `json:"repo_token"`
}
type Account struct {
	ID         int64  `xorm:"Int pk autoincr not null unique 'id'"`
	Login      string `xorm:"Varchar(15) not null unique 'login'"`
	Password   string `xorm:"Varchar(15) not null 'password'"`
	Role       int8   `xorm:"Int not null default 0 'Role'"`
	TelegramID int    `xorm:"Int unique 'telegram_id'"`
	Name       string `xorm:"Varchar(255) 'name'"`
	UserName   string `xorm:"Varchar(255) unique 'username'"`
	Verified   bool   `xorm:"Bool not null default 0 'verified'"`
}

type CachedAccount struct {
	Role     int8
	Name     string
	UserName string
	Context  int8
}

func (c *CachedAccount) SetCtx(s string) {
	cmdArgs := strings.Split(s, " ")
	ctxNum, _ := strconv.ParseInt(cmdArgs[0], 10, 64)
	c.Context = int8(ctxNum)
	return
}
func (c *CachedAccount) GetCtx() int8 {
	return c.Context
}

var (
	engine         *xorm.Engine
	cachedAccounts = make(map[int]*CachedAccount)
	configuration  Configuration
)

func errh(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func dbInit() *xorm.Engine {
	engine, err := xorm.NewEngine("sqlite3", "./en-vote_bot.db")
	errh(err)
	err = engine.Sync(new(Account))
	errh(err)
	return engine
}

func checkSender(u *tgbotapi.User) bool {
	if _, has := cachedAccounts[u.ID]; has {
		return has
	}
	a := Account{TelegramID: u.ID}
	has, err := engine.Get(&a)
	errh(err)
	if has {
		cachedAccounts[a.TelegramID] = &CachedAccount{
			Role:     a.Role,
			Name:     a.Name,
			UserName: a.UserName,
			Context:  0,
		}
	}
	return has
}

func main() {
	file, _ := os.Open("./conf.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	errh(err)
	engine = dbInit()
	bot, err := tgbotapi.NewBotAPI(configuration.TgAPIToken)
	errh(err)

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	errh(err)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		from := update.Message.From
		log.Printf("ID[%d] Nick:%s %s %s  %s", from.ID, from.UserName, from.FirstName, from.LastName, update.Message.Text)
		tgbotapi.NewChatAction(update.Message.Chat.ID, "act")

		if !checkSender(from) && update.Message.Command() != "auth" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		} else {
			if update.Message.IsCommand() {
				msg := cmdWrapper(update.Message)
				bot.Send(msg)
			} else if update.Message.Voice != nil {
				VoiceWrapper(bot, update.Message)

			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "неправильный ввод, я потом запилю командлист")
				bot.Send(msg)
			}

		}
	}
}
