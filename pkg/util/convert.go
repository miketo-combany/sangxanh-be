package util

import (
	"database/sql"
	"time"
)

func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "", // Valid is false if the string is empty
	}
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
