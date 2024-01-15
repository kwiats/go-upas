package storage

import (
	"database/sql"
	"log"

	"github.com/kwiats/go-upas/models"
	"github.com/kwiats/go-upas/utils"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*models.UserCreate) error
	DeleteUser(string) error
	UpdateUser(*models.UserDTO) error
	GetUsers() ([]*models.UserDTO, error)
	GetUserByID(string) (*models.UserDTO, error)
	GetUserCredentials(*models.UserLogin) (*models.UserLogin, error)
	SignInUser(*models.UserLogin) (*utils.TokenResponse, error)
}

func (dbs *DataBaseStore) CreateUser(user *models.UserCreate) error {
	transaction, err := dbs.db.Begin()
	defer func() {
		if err != nil {
			if rErr := transaction.Rollback(); rErr != nil {
				log.Printf("Error rolling back transaction: %v", rErr)
			}
			log.Printf("Transaction rolled back: %v", err)
		}
	}()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return err
	}

	profilID, err := InsertProfile(transaction, user)

	if err != nil {
		log.Printf("Failed to insert profile: %v", err)
		return err
	}
	userID, err := InsertUser(transaction, user, profilID)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}

	if err := transaction.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	log.Printf("Created user(UUID=%+v) and profile(UUID=%+v) object.", userID, profilID)
	return nil
}

func InsertUser(transaction *sql.Tx, user *models.UserCreate, profilID string) (string, error) {
	userQuery := `INSERT INTO users ( username, password, email,  profil_id)
				VALUES ($1, $2, $3, $4) RETURNING id`

	var userID string
	err := transaction.QueryRow(
		userQuery,
		user.Username,
		user.Password,
		user.Email,
		profilID,
	).Scan(&userID)
	return userID, err
}

func InsertProfile(transaction *sql.Tx, user *models.UserCreate) (string, error) {
	profileQuery := `INSERT INTO profils (image_profile, first_name, last_name)
                 VALUES ($1, $2, $3) RETURNING id`

	var profilID string
	err := transaction.QueryRow(
		profileQuery,
		user.ImageProfile,
		user.FirstName,
		user.LastName,
	).Scan(&profilID)
	return profilID, err
}

func (dbs *DataBaseStore) GetUsers() ([]*models.UserDTO, error) {
	query := "SELECT id, username, email FROM users"
	rows, err := dbs.db.Query(query)
	if err != nil {
		return nil, err
	}
	users := []*models.UserDTO{}
	for rows.Next() {
		user := new(models.UserDTO)
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (dbs *DataBaseStore) UpdateUser(*models.UserDTO) error {
	return nil
}

func (dbs *DataBaseStore) DeleteUser(id string) error {
	// Start a transaction
	transaction, err := dbs.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rErr := transaction.Rollback(); rErr != nil {
				log.Printf("Error rolling back transaction: %v", rErr)
			}
			log.Printf("Transaction rolled back: %v", err)
		}
	}()

	deleteProfileQuery := "DELETE FROM profils WHERE id = (SELECT profil_id FROM users WHERE id = $1)"
	if _, err := transaction.Exec(deleteProfileQuery, id); err != nil {
		return err
	}

	deleteUserQuery := "DELETE FROM users WHERE id = $1"
	if _, err := transaction.Exec(deleteUserQuery, id); err != nil {
		return err
	}

	if err := transaction.Commit(); err != nil {
		return err
	}

	return nil
}
func (dbs *DataBaseStore) GetUserByID(id string) (*models.UserDTO, error) {
	user := models.UserDTO{}
	query := "SELECT id, username, email FROM users WHERE users.id = $1"
	row := dbs.db.QueryRow(query, id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found, return nil without error
		}
		return nil, err
	}
	return &user, nil
}

func (dbs *DataBaseStore) GetUserCredentials(user *models.UserLogin) (*models.UserLogin, error) {
	userLogin := new(models.UserLogin)

	query := "SELECT id, username, password FROM users WHERE users.username = $1"
	row := dbs.db.QueryRow(query, user.Username)

	err := row.Scan(
		&userLogin.ID,
		&userLogin.Username,
		&userLogin.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found, return nil without error
		}
		return nil, err
	}
	return userLogin, nil
}
