package telebot

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
	tele "gopkg.in/tucnak/telebot.v2"
  "github.com/ondrejholik/telebot/commands"
)

type Config struct {
	Token string
}

func main() {
	var conf Config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}

	b, err := tele.NewBot(tele.Settings{

		Token:  conf.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	// Routes

	b.Handle("/hello", func(m *tele.Message) {
		b.Send(m.Sender, "Hello World!")
	})

	// SplitWise

	// Weather with GUI
	b.Handle("/w", func(m *tele.Message) {
		if weatherGui() {

		} else {
			b.Send(m.Sender, weather())
		}
	})

	// Word Count
	b.Handle("/wc", func(m *tele.Message) {
		b.Send(m.Sender, Wc(m.Text, true))
	})

	// QR generator

	// Time ( clock svg/png )

	// Other
	b.Handle(tele.OnText, func(m *tele.Message) {

	})

	b.Start()
}
