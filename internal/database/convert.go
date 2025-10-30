package database

import (
	"database/sql"
	"log"
	"time"
)

func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func ToNullTime(s string) sql.NullTime {
	if s == "" {
		return sql.NullTime{Valid: false}
	}

	var t time.Time
	var err error
	switch len(s) {
	case len(time.RFC1123Z):
		t, err = time.Parse(time.RFC1123Z, s)
	case len(time.RFC1123):
		t, err = time.Parse(time.RFC1123, s)
	case len(time.RFC822Z):
		t, err = time.Parse(time.RFC822Z, s)
	case len(time.RFC822):
		t, err = time.Parse(time.RFC822, s)
	case len(time.RFC3339):
		t, err = time.Parse(time.RFC3339, s)
	case len(time.UnixDate):
		t, err = time.Parse(time.UnixDate, s)
	default:
		t, err = time.Parse(time.DateTime, s)
	}

	if err != nil {
		log.Printf("couldn't parse time %q: %v", s, err)
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: t.UTC(), Valid: true}
}
