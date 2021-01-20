package telebot

import "github.com/dgraph-io/badger"

// Village name, gps point lat, lon
type Village struct {
	name  string
	point Point
}

// AreVillagesInDB -- check if villages has been initialized(if exists)
func AreVillagesInDB(db *badger.DB) bool {
	var found = true
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("villages"))
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

// ClosestVillage find closest village to point
func ClosestVillage(db *badger.DB, p Point) string {
	return "a"
}

func LoadVillages(db *badger.DB) {

}
