package telebot

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/slice"
)

// Date --
type Date struct {
	Countable bool
	Name      string
	Day       int
	Month     int
	Year      int
	FromToday int
	YearDiff  int
}

func isCountable(sign string) bool {
	if sign == "1" {
		return true
	}
	return false
}

func parseToDate(record []string) Date {
	fulldate := strings.Split(record[2], ".")

	var day, month, year int
	var err error

	day, err = strconv.Atoi(fulldate[0])
	if err != nil {
		log.Panic(err)
	}
	month, err = strconv.Atoi(fulldate[1])
	if err != nil {
		log.Panic(err)
	}
	year, err = strconv.Atoi(fulldate[2])
	if err != nil {
		log.Panic(err)
	}
	countable := isCountable(record[0])
	yeardiff := 0
	if countable {
		yeardiff = yearDiff(year, month, day)
	}

	return Date{Countable: countable, Name: record[1], Day: day, Month: month, Year: year, FromToday: fromToday(day, month, year), YearDiff: yeardiff}
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func sign(m1, d1, m2, d2 int) int {
	if m1 != m2 {
		if m1 < m2 {
			return 1
		}
		return -1
	}

	if d1 < d2 {
		return -1
	}
	return 1
}

func yearDiff(y1, m1, d1 int) int {

	t := time.Now()
	add := 0
	y2 := int(t.Year())
	m2 := int(t.Month())
	d2 := int(t.Day())
	if m1 < m2 {
		add = 1
	}
	if m1 == m2 && d1 < d2 {
		add = 1
	}
	return int(y2-y1) + add

}

func fromToday(day, month, year int) int {
	today := time.Now()
	date := date(year, month, day)
	m2 := today.Month()
	d2 := today.Day()
	sgn := sign(month, day, int(m2), d2)
	sub := sgn * (int((today.Sub(date).Hours() / 24)) % 365)
	if sgn == -1 {
		sub = 365 + sub
	}
	return sub
}

// Bd -- birthday handler
func Bd() string {
	var result string
	var dates []Date
	dat, err := ioutil.ReadFile("assets/birthdays.csv")
	if err != nil {
		log.Panic(err)
	}
	r := csv.NewReader(strings.NewReader(string(dat)))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		date := parseToDate(record)
		dates = append(dates, date)
	}

	slice.Sort(dates[:], func(i, j int) bool {
		return dates[i].FromToday < dates[j].FromToday
	})

	for _, x := range dates {
		if x.Countable {
			result += fmt.Sprintf("%03d -- %s(%d)\n", x.FromToday, x.Name, x.YearDiff)
		} else {
			result += fmt.Sprintf("%03d -- %s\n", x.FromToday, x.Name)
		}
	}

	return result

}
