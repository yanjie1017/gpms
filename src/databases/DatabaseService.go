package databases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	models "src/models"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type DBConnection struct {
	DB *sql.DB
}

func ConnectDB(config models.DatabaseConfiguration) DBConnection {
	contextLogger := log.WithFields(log.Fields{
		"host":     config.Host,
		"port":     config.Port,
		"user":     config.User,
		"database": config.Name,
	})

	contextLogger.Info("Connecting to database...")
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)

	// Open DB connection
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		contextLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to open database connection")
		panic(err)
	}

	// Check DB connection
	err = db.Ping()
	if err != nil {
		contextLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to connect to database")
		panic(err)
	}

	log.Info("Database connected")

	return DBConnection{db}
}

func (conn *DBConnection) CreatePasswordEntry(passwordInfo models.PasswordInfo) (int64, error) {
	contextLogger := log.WithFields(log.Fields{
		"client_id":          passwordInfo.Client.Id,
		"entry_reference_id": passwordInfo.Entry.ReferenceId,
		"query":              "INSERT_PASSWORD_ENTRY",
	})

	var id int64 = 0

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := conn.DB.BeginTx(ctx, nil)
	if err != nil {
		contextLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to execute query")
		return id, err
	}

	{
		stmt, err := tx.Prepare(INSERT_PASSWORD_ENTRY)

		if err != nil {
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to execute query")
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
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to execute query")
			return id, err
		}
	}
	{
		err := tx.Commit()

		if err != nil {
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Query execution failed")
			return id, err
		}
	}

	contextLogger.Info("Executed query")
	return id, err
}

func (conn *DBConnection) RetrievePasswordInfo(entryTag models.PasswordEntryTag) (*models.PasswordGenerationInfo, error) {
	contextLogger := log.WithFields(log.Fields{
		"client_id": entryTag.ClientId,
		"entry_id":  entryTag.EntryId,
		"query":     "RETRIEVE_PASSWORD_ENTRY",
	})

	info := &models.PasswordGenerationInfo{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := conn.DB.BeginTx(ctx, nil)
	if err != nil {
		contextLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to execute query")
		return nil, err
	}

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
		contextLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Query execution failed")
		return nil, err
	}

	return info, err
}

func (conn *DBConnection) RetrieveAPIKey(clientId int64) (*models.ClientAuthentication, error) {
	contextLogger := log.WithFields(log.Fields{
		"client_id": clientId,
		"query":     "RETRIEVE_API_KEY",
	})

	info := &models.ClientAuthentication{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := conn.DB.BeginTx(ctx, nil)
	if err != nil {
		contextLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to execute query")
		return nil, err
	}

	err = tx.QueryRowContext(
		ctx,
		RETRIEVE_API_KEY,
		clientId,
	).Scan(
		&info.ClientId,
		&info.APIKey,
	)

	if err != nil {
		contextLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Query execution failed")
		return nil, err
	}

	return info, err
}
