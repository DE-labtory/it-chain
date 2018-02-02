// Copyright 2013 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package timeago

import (
	"testing"
	"time"
)

// Base time for testing
var tBase = time.Date(2013, 8, 30, 12, 0, 0, 0, time.UTC)

// Test data for TestFormatReference
var formatReferenceTests = []struct {
	t        time.Time // input time
	ref      time.Time // input reference
	cfg      Config    // input cfguage
	expected string    // expected result
}{
	// Lang
	{tBase, tBase, NoMax(English), "about a second ago"},
	{tBase, tBase, NoMax(French), "il y a environ une seconde"},
	{tBase, tBase, NoMax(Chinese), "1 秒前"},
	{tBase, tBase, NoMax(Portuguese), "há menos de um segundo"},
	{tBase, tBase, NoMax(German), "vor einer Sekunde"},

	// Thresholds
	{tBase, tBase.Add(1*time.Second + 500000000).Add(-1), NoMax(English), "about a second ago"},
	{tBase, tBase.Add(1*time.Second + 500000000), NoMax(English), "2 seconds ago"},
	{tBase, tBase.Add(1 * time.Minute), NoMax(English), "about a minute ago"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second).Add(-1), NoMax(English), "about a minute ago"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second), NoMax(English), "2 minutes ago"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second), NoMax(English), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute), NoMax(English), "2 hours ago"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute).Add(-1), NoMax(English), "23 hours ago"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute), NoMax(English), "one day ago"},
	{tBase, tBase.Add(36 * time.Hour), NoMax(English), "2 days ago"},
	{tBase, tBase.Add(30 * 24 * time.Hour), NoMax(English), "one month ago"},
	{tBase, tBase.Add(Year).Add(-30 * Day), NoMax(English), "11 months ago"},
	{tBase, tBase.Add(Year), NoMax(English), "one year ago"},
	{tBase, tBase.Add(547 * Day), NoMax(English), "one year ago"},
	{tBase, tBase.Add(548 * Day), NoMax(English), "2 years ago"},
	{tBase, tBase.Add(10 * Year), NoMax(English), "10 years ago"},

	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(Portuguese), "há uma hora"},
	{tBase, tBase.Add(45 * 24 * time.Hour).Add(-1), NoMax(Portuguese), "há um mês"},
	{tBase, tBase.Add(36 * time.Hour).Add(-1), NoMax(Portuguese), "há um dia"},
	{tBase, tBase.Add(1 * time.Minute).Add(-500000001), NoMax(Portuguese), "há 59 segundos"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second).Add(-1), NoMax(Portuguese), "há 59 minutos"},
	{tBase, tBase.Add(30 * 24 * time.Hour).Add(-12*time.Hour - 1), NoMax(Portuguese), "há 29 dias"},
	{tBase, tBase.Add(45 * 24 * time.Hour), NoMax(Portuguese), "há 2 meses"},
	{tBase, tBase.Add(10 * Year), NoMax(Portuguese), "há 10 anos"},

	{tBase, tBase.Add(1*time.Second + 500000000).Add(-1), NoMax(German), "vor einer Sekunde"},
	{tBase, tBase.Add(1*time.Second + 500000000), NoMax(German), "vor 2 Sekunden"},
	{tBase, tBase.Add(1 * time.Minute), NoMax(German), "vor einer Minute"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second).Add(-1), NoMax(German), "vor einer Minute"},
	{tBase, tBase.Add(1*time.Minute + 30*time.Second), NoMax(German), "vor 2 Minuten"},
	{tBase, tBase.Add(59*time.Minute + 30*time.Second), NoMax(German), "vor einer Stunde"},
	{tBase, tBase.Add(90 * time.Minute), NoMax(German), "vor 2 Stunden"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute).Add(-1), NoMax(German), "vor 23 Stunden"},
	{tBase, tBase.Add(23*time.Hour + 30*time.Minute), NoMax(German), "vor einem Tag"},
	{tBase, tBase.Add(36 * time.Hour), NoMax(German), "vor 2 Tagen"},
	{tBase, tBase.Add(30 * 24 * time.Hour), NoMax(German), "vor einem Monat"},
	{tBase, tBase.Add(Year).Add(-30 * Day), NoMax(German), "vor 11 Monaten"},
	{tBase, tBase.Add(Year), NoMax(German), "vor einem Jahr"},
	{tBase, tBase.Add(547 * Day), NoMax(German), "vor einem Jahr"},
	{tBase, tBase.Add(548 * Day), NoMax(German), "vor 2 Jahren"},
	{tBase, tBase.Add(10 * Year), NoMax(German), "vor 10 Jahren"},

	// Max
	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(English), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), WithMax(English, 90*time.Minute, ""), "about an hour ago"},
	{tBase, tBase.Add(90 * time.Minute), WithMax(English, 90*time.Minute, "2006-01-02"), "2013-08-30"},

	{tBase, tBase.Add(90 * time.Minute), WithMax(Portuguese, 90*time.Minute, "01-02-2006"), "08-30-2013"},

	{tBase, tBase.Add(90 * time.Minute).Add(-1), NoMax(German), "vor einer Stunde"},
	{tBase, tBase.Add(90 * time.Minute).Add(-1), WithMax(German, 90*time.Minute, ""), "vor einer Stunde"},
	{tBase, tBase.Add(90 * time.Minute), WithMax(German, 90*time.Minute, German.DefaultLayout), "30.08.2013"},

	// Future
	{tBase.Add(24 * time.Hour), tBase, NoMax(English), "in one day"},

	{tBase.Add(2 * Month), tBase, NoMax(Portuguese), "daqui a 2 meses"},
	{tBase.Add(24 * time.Hour), tBase, NoMax(Portuguese), "daqui a um dia"},
	{tBase.Add(5 * time.Minute), tBase, NoMax(Portuguese), "daqui a 5 minutos"},
	{tBase.Add(1 * time.Minute), tBase, NoMax(Portuguese), "daqui a um minuto"},
	{tBase.Add(100 * time.Millisecond), tBase, NoMax(Portuguese), "daqui a menos de um segundo"},

	{tBase.Add(2 * Month), tBase, NoMax(German), "in 2 Monaten"},
	{tBase.Add(24 * time.Hour), tBase, NoMax(German), "in einem Tag"},
	{tBase.Add(5 * time.Minute), tBase, NoMax(German), "in 5 Minuten"},
	{tBase.Add(1 * time.Minute), tBase, NoMax(German), "in einer Minute"},
	{tBase.Add(100 * time.Millisecond), tBase, NoMax(German), "in einer Sekunde"},
}

// Test the FormatReference method
func TestFormatReference(t *testing.T) {
	for i, tt := range formatReferenceTests {
		actual := tt.cfg.FormatReference(tt.t, tt.ref)
		if actual != tt.expected {
			t.Errorf("%d) FormatReference(%s,%s): expected '%s', actual '%s'", i+1, tt.t, tt.ref, tt.expected, actual)
		}
	}
}
