package util

import "database/sql"

func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "", // Valid is false if the string is empty
	}
}
