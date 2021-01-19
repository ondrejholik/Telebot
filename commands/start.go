package telebot

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	badger "github.com/dgraph-io/badger"
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
func Start() {

	// Does user exists in database?
	err := db.View(func(txn *badger.Txn) error {
		userid := strconv.Itoa(m.Sender.ID)
		item, err := txn.Get([]byte(userid))
		if err != nil {
			// User doesnt exist
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
					return err
				}

				err = txn.Set([]byte(userid), []byte(encoded))
				return err
			})
			if err != nil {
				log.Println("Error when getting user, creating new user")
			}
		}

		var valCopy []byte
		err = item.Value(func(val []byte) error {
			// Accessing val here is valid.
			valCopy = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			log.Println("Error reading iteam")
			log.Println(err)
		}

		fmt.Printf("The answer is: %s\n", valCopy)
		b.Send(m.Sender, "valCopy")
		return nil
	})
	if err != nil {
		log.Println("Error view db")
		log.Println(err)
	}

}
