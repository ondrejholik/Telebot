package telebot

import (
	"encoding/json"
	"log"
	"strconv"

	badger "github.com/dgraph-io/badger"
	tele "gopkg.in/tucnak/telebot.v2"
)

// Point - Gps point(Latitude, Longtitude)
type Point struct {
	Lat float64
	Lon float64
}

// UserSettings structure which will be saved as json to db
type UserSettings struct {
	Username           string
	LastLocation       Point
	LastWeatherRequest int32
}

// Start initialize user(put user data to BadgerDB)
func Start(db *badger.DB, m *tele.Message) string {
	// Does user exists in database?
	userid := strconv.Itoa(m.Sender.ID)
	userexists := UserExists(db, userid)
	var msg string

	// If user exists, then print something and end
	if userexists {
		msg = "You are already in our database. Have a nice day!"
	} else {
		CreateNewUser(userid, db, m)
		msg = "You are now successfuly in our database."
	}
	return msg
}

// CreateNewUser -> add user with key userid to badgerdB
func CreateNewUser(userid string, db *badger.DB, m *tele.Message) {
	err := db.Update(func(txn *badger.Txn) error {
		userset := &UserSettings{
			Username: m.Sender.Username,
			LastLocation: Point{
				Lat: 0.0,
				Lon: 0.0,
			},
			LastWeatherRequest: 0,
		}
		encoded, err := json.Marshal(userset)
		if err != nil {
			log.Println("Error when encoding json")
			log.Println(err)
			return err
		}

		e := badger.NewEntry([]byte(userid), []byte(encoded))
		err = txn.SetEntry(e)
		if err != nil {
			log.Println("Error, when set json value to user")
			log.Println(err)
		}

		return err
	})
	if err != nil {
		log.Println(err)
	}
}

// UserExists check if user exists in badgerDB
func UserExists(db *badger.DB, userid string) bool {
	var found = true
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(userid))
		if err != nil {
			found = false
		}
		return nil
	})
	if err != nil {
		return false
	}

	return found
}
