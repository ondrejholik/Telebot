package telebot

import (
	"strconv"

	badger "github.com/dgraph-io/badger"
	tele "gopkg.in/tucnak/telebot.v2"
)

// LocationHandler handle incoming location
func LocationHandler(db *badger.DB, m *tele.Message) string {
	userid := strconv.Itoa(m.Sender.ID)
	if !UserExists(db, userid) {
		CreateNewUser(userid, db, m)
	}
	msg := UpdateLocationDB(db, m, userid, Point{Lat: m.Location.Lat, Lon: m.Location.Lng})
	return msg
}

// UpdateLocationDB change coordinates
func UpdateLocationDB(db *badger.DB, m *tele.Message, userid string, p Point) string {
	userset := GetUserSettings(db, userid)
	userset.LastLocation = p
	UpdateUser(userid, userset, db, m)
	return LocationInfo(db, p)
}

// LocationInfo info about city you are currently in
func LocationInfo(db *badger.DB, p Point) string {
	return "City: "
}

/*
func getCities() {

}
*/
