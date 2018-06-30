package models

import "rosella/dbhandler"

// SessionUID Returns the uid for a specified session.
func SessionUID(sessionID string) (string, error) {
	sqlQuery := "SELECT uid FROM public.session WHERE id = $1"
	var uid string
	err := dbhandler.GetDatabase().QueryRow(sqlQuery, sessionID).Scan(&uid)
	return uid, err
}
