package db

import (
	url "golang/db/model"
	"golang/internal/util"
)

func (db *Database) GetUrls() ([]url.Url, error) {
	trace := util.CreateErrorContext("db.GetUrls")
	urls := []url.Url{}

	if _err := db.sqlx.Select(&urls, db.driver.GetUrls()); _err != nil {
		return urls, trace.Apply(_err)
	}

	return urls, nil
}

func (db *Database) DeleteUrl(id string) error {
	trace := util.CreateErrorContext("url.DeleteUrl")

	tx, txErr := db.sqlx.Beginx()
	if txErr != nil {
		return trace.Apply(txErr)
	}

	defer tx.Rollback()

	_, deleteErr := tx.Exec(
		db.driver.DeleteUrl(),
		id,
	)
	if deleteErr != nil {
		return trace.Apply(deleteErr)
	}

	if err := tx.Commit(); err != nil {
		return trace.Apply(err)
	}

	return nil

}

func (db *Database) SaveUrl(url url.Url) error {
	trace := util.CreateErrorContext("url.SaveUrl")

	tx, txErr := db.sqlx.Beginx()
	if txErr != nil {
		return trace.Apply(txErr)
	}

	defer tx.Rollback()

	_, saveErr := tx.Exec(
		db.driver.SaveUrl(),
		url.Id,
		url.Original,
	)
	if saveErr != nil {
		return trace.Apply(saveErr)
	}

	if err := tx.Commit(); err != nil {
		return trace.Apply(err)
	}

	return nil
}

func (db *Database) GetUrlByUrl(original string) (url.Url, error) {
	trace := util.CreateErrorContext("url.GetUrlByUrl")
	urlRecord := url.Url{}

	if _err := db.sqlx.Get(&urlRecord, db.driver.GetUrlByUrl(), original); _err != nil {
		return urlRecord, trace.Apply(_err)
	}

	return urlRecord, nil
}


func (db *Database) UpdateAccesses(id string, accesses int) error {
	trace := util.CreateErrorContext("db.UpdateAccesses")

	if _, err := db.sqlx.Exec(db.driver.UpdateAccesses(), accesses, id); err != nil {
		return trace.Apply(err)
	}

	return nil
}
