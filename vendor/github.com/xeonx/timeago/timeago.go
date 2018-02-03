// Copyright 2013 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//timeago allows the formatting of time in terms of fuzzy timestamps.
//For example:
//	one minute ago
//	3 years ago
//	in 2 minutes
package timeago

import (
	"fmt"
	"strings"
	"time"
)

//Config allows the customization of timeago.
//You may configure string items (cfguage, plurals, ...) and
//maximum allowed duration value for fuzzy formatting.
type Config struct {
	PastPrefix    string
	PastSuffix    string
	FuturePrefix  string
	FutureSuffix  string
	Second        string
	Seconds       string
	Minute        string
	Minutes       string
	Hour          string
	Hours         string
	Day           string
	Days          string
	Month         string
	Months        string
	Year          string
	Years         string
	Max           time.Duration //Maximum duration for using the special formatting.
	DefaultLayout string        //Layout to use if delta is greater than Max
}

//Predefined english configuration
var English = Config{
	PastPrefix:    "",
	PastSuffix:    " ago",
	FuturePrefix:  "in ",
	FutureSuffix:  "",
	Second:        "about a second",
	Seconds:       "less than a minute",
	Minute:        "about a minute",
	Minutes:       "%d minutes",
	Hour:          "about an hour",
	Hours:         "%d hours",
	Day:           "one day",
	Days:          "%d days",
	Month:         "one month",
	Months:        "%d months",
	Year:          "one year",
	Years:         "%d years",
	Max:           73 * time.Hour,
	DefaultLayout: "2006-01-02",
}

//Predefined french configuration
var French = Config{
	PastPrefix:    "il y a ",
	PastSuffix:    "",
	FuturePrefix:  "dans ",
	FutureSuffix:  "",
	Second:        "environ une seconde",
	Seconds:       "moins d'une minute",
	Minute:        "environ une minute",
	Minutes:       "%d minutes",
	Hour:          "environ une heure",
	Hours:         "%d heures",
	Day:           "un jour",
	Days:          "%d jours",
	Month:         "un mois",
	Months:        "%d mois",
	Year:          "un an",
	Years:         "%d ans",
	Max:           73 * time.Hour,
	DefaultLayout: "02/01/2006",
}

//Format returns a textual representation of the time value formatted according to the layout
//defined in the Config. The time is compared to time.Now() and is then formatted as a fuzzy
//timestamp (eg. "4 days ago")
func (cfg Config) Format(t time.Time) string {
	return cfg.FormatReference(t, time.Now())
}

//FormatReference is the same as Format, but the reference has to be defined by the caller
func (cfg Config) FormatReference(t time.Time, reference time.Time) string {

	d := reference.Sub(t)

	if (d >= 0 && d >= cfg.Max) || (d < 0 && -d >= cfg.Max) {
		return t.Format(cfg.DefaultLayout)
	}

	return cfg.FormatDuration(d)
}

//FormatReference is the same as Format, but for time.Duration.
//Config.Max is not used in this function, as there is no other alternative.
func (cfg Config) FormatDuration(duration time.Duration) string {

	isPast := duration >= 0

	if duration < 0 {
		duration = -duration
	}

	s := cfg.getTimeText(duration)

	if isPast {
		return strings.Join([]string{cfg.PastPrefix, s, cfg.PastSuffix}, "")
	} else {
		return strings.Join([]string{cfg.FuturePrefix, s, cfg.FutureSuffix}, "")
	}
}

//Count the number of parameters in a format string
func nbParamInFormat(f string) int {
	return strings.Count(f, "%") - 2*strings.Count(f, "%%")
}

//Round the duratiion d in termes of step.
func round(d time.Duration, step time.Duration) int64 {
	return int64(float64(d)/float64(step) + 0.5)
}

//Convert a duration to a text, based on the current cfguage
func (cfg Config) getTimeText(d time.Duration) string {

	//Less than 1.5 second
	if d < 1500*time.Millisecond {
		return cfg.Second
	}

	//Less than 1 minute
	if d < 60*time.Second {

		switch nbParamInFormat(cfg.Seconds) {
		case 1:
			return fmt.Sprintf(cfg.Seconds, round(d, time.Second))
		}
		return cfg.Seconds
	}

	//Less than 1.5 minute
	if d < 90*time.Second {
		return cfg.Minute
	}

	//Less than 1 hour
	if d < 59*time.Minute+30*time.Second {

		switch nbParamInFormat(cfg.Minutes) {
		case 1:
			return fmt.Sprintf(cfg.Minutes, round(d, time.Minute))
		}
		return cfg.Minutes
	}

	//Less than 1.5 hour
	if d < 90*time.Minute {
		return cfg.Hour
	}

	//Less than 1 day
	if d < 23*time.Hour+30*time.Minute {

		switch nbParamInFormat(cfg.Hours) {
		case 1:
			return fmt.Sprintf(cfg.Hours, round(d, time.Hour))
		}
		return cfg.Hours
	}

	//Less than 1.5 day
	if d < 36*time.Hour {
		return cfg.Day
	}

	//Less than 30 days
	if d < 30*24*time.Hour {

		switch nbParamInFormat(cfg.Days) {
		case 1:
			return fmt.Sprintf(cfg.Days, round(d, time.Hour*24))
		}
		return cfg.Days
	}

	//Less than 1.5 month
	if d < 45*24*time.Hour {
		return cfg.Month
	}

	//Less than 1 year
	if d < 365*24*time.Hour {

		switch nbParamInFormat(cfg.Months) {
		case 1:
			return fmt.Sprintf(cfg.Months, round(d, time.Hour*30*24))
		}

		return cfg.Months
	}

	//Less than 1.5 year
	if d < 548*24*time.Hour { //548 days = 1.5 years
		return cfg.Year
	}

	switch nbParamInFormat(cfg.Years) {
	case 1:
		return fmt.Sprintf(cfg.Years, round(d, time.Hour*24*365))
	}
	return cfg.Years

}

//NoMax creates an new config without a maximum
func NoMax(cfg Config) Config {
	return WithMax(cfg, 9223372036854775807, "")
}

//WithMax creates an new config with special formatting limited to durations less than max.
//Values greater than max will be formatted by the standard time package using the defaultLayout.
func WithMax(cfg Config, max time.Duration, defaultLayout string) Config {
	n := cfg
	n.Max = max
	n.DefaultLayout = defaultLayout
	return n
}
