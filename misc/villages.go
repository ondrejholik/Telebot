package telebot

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
)

// Gps struct with Latitude, Longtitude
type Gps struct {
	Lat float32
	Lon float32
}

// Village name, gps point lat, lon
type Village struct {
	Name  string
	Point Gps
}

// VillageArr is there only, because json.Marshal
type VillageArr struct {
	Arr []Village
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
func ClosestVillage(db *badger.DB, lat float32, lon float32) string {
	return "a"
}

// ParseToVillage take slice of string splitted by delimiter and convert to Village struct
func ParseToVillage(row []string) Village {
	var lat, lon float32
	val, err := strconv.ParseFloat(row[1], 32)
	if err != nil {
		log.Println(err)
	}
	lat = float32(val)
	val, err = strconv.ParseFloat(row[2], 32)
	if err != nil {
		log.Println(err)
	}
	lon = float32(val)
	return Village{Name: row[0], Point: Gps{Lat: lat, Lon: lon}}
}

// LoadVillagesToDB load to badgerDB
func LoadVillagesToDB(db *badger.DB) {
	var villages []Village

	dat, err := ioutil.ReadFile("assets/villages.csv")
	if err != nil {
		log.Panic(err)
	}
	r := csv.NewReader(strings.NewReader(string(dat)))
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

	err = db.Update(func(txn *badger.Txn) error {
		vlgsStruct := &VillageArr{
			Arr: villages,
		}
		encoded, err := json.Marshal(vlgsStruct)
		if err != nil {
			log.Println(err)
		}
		err = txn.Set([]byte("villages"), []byte(encoded))
		return err
	})
	if err != nil {
		log.Panic(err)
	}
}
