package database

func DeleteSessions() error {
	if err := db.Exec("delete from sessions").Error; err != nil {
		return err
	}
	return nil
}
