package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	badger "github.com/dgraph-io/badger"
	tele "gopkg.in/tucnak/telebot.v2"
)

type Config struct {
	Token string
}

type Point struct {
	Lat float64
	Lon float64
}

type UserSettings struct {
	Username           string
	UserId             string
	LastLocation       Point
	LastWeatherRequest int32
}

func main() {
	var conf Config
	db, err := badger.Open(badger.DefaultOptions("./badger/"))
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

	// Routes
	b.Handle("/start", func(m *tele.Message) {
		// Does user exists in database?
		err := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(m.InlineID))

			var valCopy []byte
			err = item.Value(func(val []byte) error {
				// Accessing val here is valid.
				fmt.Printf("The answer is: %s\n", val)

				// Copying or parsing val is valid.
				valCopy = append([]byte{}, val...)

				return nil
			})
			if err != nil {
				log.Println(err)
			}

			fmt.Printf("The answer is: %s\n", valCopy)
			b.Send(m.Sender, valCopy)
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		// Create user in DB if not exist already
		/*
		   err := db.Update(func(txn *badger.Txn) error {
		     if err != nil {
		       log.Panic(err)
		     }


		     return nil
		   })
		*/

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
