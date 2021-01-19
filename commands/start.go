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

  log.Println("Someone trigger start")
	//db := ctx.Value("db")
	var status string
	// Does user exists in database?
	err := db.View(func(txn *badger.Txn) error {
		userid := strconv.Itoa(m.Sender.ID)
		_, err := txn.Get([]byte(userid))
		if err != nil {
			// User doesnt exist
      log.Println("User does not exists")
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
				status = "Success, your are now initialized in our database"
				return err
			})
			if err != nil {
				log.Println("User does not exists, creating new user")
			}
		} else {
		  status = "Your are already in our database"
    }
		return nil

	})
	if err != nil {
		log.Println("Error view db")
		log.Println(err)
	}

	return status

}
