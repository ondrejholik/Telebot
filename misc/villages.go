package telebot

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
)

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

func ParseToVillage(row []string) Village {
	return Village{row[0], Point{Lat: strconv.Itoa(row[1]), Lon: strconv.Itoa(row[2])}}
}

func LoadVillagesToDB(db *badger.DB) {
	var villages []Village

	r := csv.NewReader(strings.NewReader(in))
	r.Comma = ';'

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		village := ParseToVillage(record)
		villages = append(villages, village)
	}

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("villages"), []byte(json.Marshall(villages)))
		return err
	})
	if err != nil {
		log.Panic(err)
	}
}
