package telebot

import (
	"encoding/json"
	"log"

	badger "github.com/dgraph-io/badger"
	tele "gopkg.in/tucnak/telebot.v2"
)

// Point - Gps point(Latitude, Longtitude)
type Point struct {
	Lat float32
	Lon float32
}

// UserSettings structure which will be saved as json to db
type UserSettings struct {
	Username           string
	LastLocation       Point
	LastWeatherRequest int32
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

// UpdateUser settings into badgerDB
func UpdateUser(userid string, userset UserSettings, db *badger.DB, m *tele.Message) {
	err := db.Update(func(txn *badger.Txn) error {
		encoded, err := json.Marshal(userset)
		if err != nil {
			log.Println("Error when encoding json")
			log.Println(err)
			return err
		}

		err = txn.Set([]byte(userid), []byte(encoded))
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

// GetUserSettings TODO
func GetUserSettings(db *badger.DB, userid string) UserSettings {
	var userset UserSettings
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userid))
		if err != nil {
			log.Println(err)
		}

		var valCopy []byte
		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			log.Println(err)
		}

		// Alternatively, you could also use item.ValueCopy().
		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			log.Println(err)
		}

		// From json to struct
		err = json.Unmarshal(valCopy, &userset)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return userset
}
