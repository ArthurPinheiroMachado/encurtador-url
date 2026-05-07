package postgres

import "fmt"

type Driver struct{}

func (c Driver) QueryString(user, pass, host, name, port string) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s port=%s sslmode=disable",
		user, pass, host, name, port,
	)
}

func (c Driver) InsertStage() string {
	return "INSERT INTO migrations(id, content) VALUES($1, $2)"
}

func (c Driver) LastPosition() string {
	return "SELECT MAX(id) FROM migrations"
}

func (c Driver) InitDatabase() []string {
	return []string{
		createMigrationTable(),
		createUrlTable(),
	}
}

func createMigrationTable() string {
	return `
	CREATE TABLE migrations(
			id INT NOT NULL,
			content TEXT NOT NULL,
			PRIMARY KEY(id)
	);`
}

func createUrlTable() string {
	return `
		CREATE TABLE url(
			id TEXT NOT NULL,
			original TEXT NOT NULL,
			accesses BIGINT DEFAULT 0,
			UNIQUE(original),
			PRIMARY KEY(id)
		)
	`
}

func (c Driver) GetUrls() string {
	return `
		SELECT * from url
	`
}

func (c Driver) DeleteUrl() string {
	return `
	DELETE FROM url WHERE id = $1

	`
}

func (c Driver) GetUrlByUrl() string {
	return `
	SELECT * from url WHERE original = $1
	`
}

func (c Driver) SaveUrl() string {
	return `
	INSERT INTO url(
		id, 
		original, 
		accesses
	)
	VALUES(
		$1, 
		$2, 
		0
	)
	`
}

func (c Driver) UpdateAccesses() string {
	return `
	UPDATE url SET accesses = $1 WHERE id = $2
	`
}
