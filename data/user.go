package data

import (
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int
	Name      string
	Email     string
	Password  []byte
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	UserId    int
	CreatedAt time.Time
}

// CreateSession creates a new session for existing user
func (u *User) CreateSession() (session Session, err error) {
	statement := "INSERT INTO sessions (uuid, user_id, created_at) VALUES ($1, $2, $3) RETURNING id, uuid, user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the Session struct
	uuidV4, err := uuid.NewV4()
	if err != nil {
		return
	}
	err = stmt.QueryRow(uuidV4, u.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.UserId, &session.CreatedAt)
	return
}

// Check if session is valid in the database
func (s *Session) Check() (err error) {
	err = Db.QueryRow("SELECT id, uuid, user_id, created_at FROM sessions WHERE uuid = $1", s.Uuid).
		Scan(&s.Id, &s.Uuid, &s.UserId, &s.CreatedAt)
	return
}

// DeleteByUUID deletes session for database
func (s *Session) DeleteByUUID() (err error) {
	statement := "DELETE FROM sessions WHERE uuid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.Uuid)
	return
}

// Create a new user, save user info into database
func (u *User) Create() (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := "INSERT INTO users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) RETURNING id, uuid, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// generate hash for user password
	bs, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return
	}
	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt.QueryRow(u.Name, u.Email, string(bs), time.Now()).
		Scan(&u.Id, &u.CreatedAt)
	return
}

// Get a single user by email
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, name, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}
