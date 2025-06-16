package bot

import (
	"github.com/mymmrac/telego"
	"github.com/nktauserum/anonymous-messages/config"
	"sync"
)

var (
	once sync.Once
	bot  *telego.Bot
)

func LoadBot() (*telego.Bot, error) {
	var global_err error
	once.Do(func() {
		c := config.MustLoadConfig()

		tgbot, err := telego.NewBot(c.Telegram.Token, telego.WithDefaultDebugLogger())
		if err != nil {
			global_err = err
			return
		}

		bot = tgbot
	})

	return bot, global_err
}
