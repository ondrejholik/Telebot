package telebot

import (
	"strconv"

	badger "github.com/dgraph-io/badger"
	tele "gopkg.in/tucnak/telebot.v2"
)

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
