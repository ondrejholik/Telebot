package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	badger "github.com/dgraph-io/badger"
	cmd "github.com/ondrejholik/telebot/commands"
	misc "github.com/ondrejholik/telebot/misc"
	tele "gopkg.in/tucnak/telebot.v2"
)

// Get rid off global variable
var db *badger.DB

// FileConfig for storing telegram token
type FileConfig struct {
	Token string
}

func main() {
	var err error
	db, err = badger.Open(badger.DefaultOptions("badger/"))

	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	if !misc.AreVillagesInDB(db) {
		log.Println("Villages not int DB")
		log.Println("Adding villages to DB")
		misc.LoadVillagesToDB(db)
	}

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
	/////////

	// Start
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

	// Birthdays reminder
	b.Handle("/bd", func(m *tele.Message) {
		b.Send(m.Sender, cmd.Bd())
	})

	// Youtube download
	b.Handle("/yt", func(m *tele.Message) {
		path, err := cmd.YtDownload(m.Text)
		if err != nil {
			b.Send(m.Sender, fmt.Sprintf("Not valid link(error)\n"))
		}
		p := &tele.Video{File: tele.FromDisk(path)}
		_, err = b.Send(m.Sender, p)
		if err != nil {
			b.Send(m.Sender, fmt.Sprint("Something break when sending:\n"))
		} else {
			go os.Remove(path)
		}
	})

	// QR generator

	// Time ( clock svg/png )

	// Other
	b.Handle(tele.OnText, func(m *tele.Message) {

	})

	// On Location
	b.Handle(tele.OnLocation, func(m *tele.Message) {
		b.Send(m.Sender, cmd.LocationHandler(db, m))
	})

	b.Start()
}
