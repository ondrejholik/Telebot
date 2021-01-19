package main

import (
	//"context"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	badger "github.com/dgraph-io/badger"
	tele "gopkg.in/tucnak/telebot.v2"
)

// FileConfig for storing telegram token
type FileConfig struct {
	Token string
}



func init() {
	var conf FileConfig
	db, err := badger.Open(badger.DefaultOptions("badger/"))
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	b, err := tele.NewBot(tele.Settings{

		Token:  conf.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	// Add variables to context

}

func main() {

	// Routes
	b.Handle("/start", func(m *tele.Message) {
		start(m *tele.Message, )
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
		b.Send(m.Sender, wc(m.Text, true))
	})

	// QR generator
	b.Handle("/qr", func(m *tele.Message) {
		//b.Send(m.Sender, qr(m.Text))
	})

	// Time ( clock svg/png )

	// Other
	b.Handle(tele.OnText, func(m *tele.Message) {
	})

	b.Start()
}
