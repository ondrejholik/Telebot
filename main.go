package telebot

import (
	"log"
	"time"

  cmd "github.com/ondrejholik/telebot/commands"
	"github.com/BurntSushi/toml"
	badger "github.com/dgraph-io/badger"
	tele "gopkg.in/tucnak/telebot.v2"
)

// Get rid off global variable
var db *badger.DB

// FileConfig for storing telegram token
type FileConfig struct {
	Token string
}

func init() {
	var err error
	db, err = badger.Open(badger.DefaultOptions("badger/"))

	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	//c := context.Background()
	//con := context.WithValue(c, "db", db)
}

func main() {

	var conf FileConfig
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
	b.Handle("/start", func(m *tele.Message) {
		msg := cmd.Start(db, m)
    b.Send(m.Sender, msg)
	})

	// SplitWise

	// Weather with GUI
	b.Handle("/w", func(m *tele.Message) {
		if cmd.WeatherGui() {

		} else {
			b.Send(m.Sender, cmd.Weather())
		}
	})

	// Word Count
	b.Handle("/wc", func(m *tele.Message) {
		b.Send(m.Sender, cmd.Wc(m.Text, true))
	})

	// QR generator

	// Time ( clock svg/png )

	// Other
	b.Handle(tele.OnText, func(m *tele.Message) {

	})

	b.Start()
}
