package databases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	models "src/models"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "gpms"
)

type DBConnection struct {
	DB *sql.DB
}

func ConnectDB() DBConnection {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

	// open database
	db, err := sql.Open("postgres", psqlConn)
	CheckError(err)

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("DB Connected!")

	return DBConnection{db}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func (conn *DBConnection) CreatePasswordEntry(passwordInfo models.PasswordInfo) (int64, error) {
	var id int64 = 0

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := conn.DB.BeginTx(ctx, nil)
	if err != nil {
		return id, err
	}

	{
		stmt, err := tx.Prepare(INSERT_PASSWORD_ENTRY)

		if err != nil {
			return id, err
		}

		defer stmt.Close()

		err = stmt.QueryRow(
			passwordInfo.Client.Id,
			passwordInfo.Site.Name,
			passwordInfo.Site.Type_,
			passwordInfo.Site.Metadata,
			passwordInfo.Site.Username,
			passwordInfo.Entry.ReferenceId,
			passwordInfo.Entry.Length,
		).Scan(&id)

		if err != nil {
			return id, err
		}
	}
	{
		err := tx.Commit()

		if err != nil {
			return id, err
		}
	}

	return id, err
}

func (conn *DBConnection) RetrievePasswordInfo(entryTag models.PasswordEntryTag) (*models.PasswordGenerationInfo, error) {
	info := &models.PasswordGenerationInfo{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := conn.DB.BeginTx(ctx, nil)
	CheckError(err)

	err = tx.QueryRowContext(
		ctx,
		RETRIEVE_PASSWORD_ENTRY,
		entryTag.ClientId,
		entryTag.EntryId,
	).Scan(
		&info.Length,
		&info.Metadata,
	)

	if err != nil {
		fmt.Println(err)
	}

	return info, err
}

func (conn *DBConnection) ListTables() (string, error) {
	var s string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := conn.DB.BeginTx(ctx, nil)
	CheckError(err)

	x := `select table_name from information_schema.tables where table_schema = 'public'`

	err = tx.QueryRowContext(
		ctx,
		x,
	).Scan(
		&s,
	)

	return s, err
}

func (conn *DBConnection) RetrieveAPIKey(clientId int64) (*models.ClientAuthentication, error) {
	info := &models.ClientAuthentication{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := conn.DB.BeginTx(ctx, nil)
	CheckError(err)

	err = tx.QueryRowContext(
		ctx,
		RETRIEVE_API_KEY,
		clientId,
	).Scan(
		&info.ClientId,
		&info.APIKey,
	)

	return info, err
}
