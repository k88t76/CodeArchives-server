package models

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Session is
type Session struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	UserID    string `json:"userID"`
	UserName  string `json:"userName"`
	CreatedAt string `json:"createdAt"`
}

// User is
type User struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
}

func NewUser(id int64, uuid string, name string, password string, createdAt string) *User {
	return &User{
		id,
		uuid,
		name,
		password,
		createdAt,
	}
}

func NewSession(id int64, uuid string, userID string, userName string, createdAt string) *Session {
	return &Session{
		id,
		uuid,
		userID,
		userName,
		createdAt,
	}
}

func (s *Session) Check() (valid bool, err error) {
	cmd := fmt.Sprintf("SELECT id, uuid, user_id, user_name, created_at FROM %s WHERE uuid = ?", tableNameSessions)
	//err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
	//	Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	row := db.QueryRow(cmd, s.UUID)
	err = row.Scan(&s.ID, &s.UUID, &s.UserID, &s.UserName, &s.CreatedAt)
	if err != nil {
		valid = false
		return valid, err
	}
	valid = true
	return valid, err
}

func (s *Session) DeleteByUUID() error {
	//statement := "delete from sessions where uuid = $1"
	cmd := fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", tableNameSessions)
	//stmt, err := Db.Prepare(statement)
	_, err := db.Exec(cmd, s.UUID)
	if err != nil {
		return err
	}
	return err
}

func GetUser(sessionUUID string) *User {
	cmd := fmt.Sprintf("SELECT id, uuid, name, password, created_at FROM %s WHERE uuid = ?", tableNameUsers)
	//err = Db.QueryRow("SELECT uuid, name, email, created_at FROM users WHERE id = $1", session.UserId).
	row := db.QueryRow(cmd, sessionUUID)
	var user User
	err := row.Scan(&user.ID, &user.UUID, &user.Name, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil
	}
	return NewUser(user.ID, user.UUID, user.Name, user.Password, user.CreatedAt)
}

func SessionDeleteAll() error {
	cmd := fmt.Sprintf("DELETE FROM %s", tableNameSessions)
	_, err := db.Exec(cmd)
	return err
}

// Create is
func (u *User) Create() error {
	u.UUID = CreateUUID()
	cmd := fmt.Sprintf("INSERT INTO %s (uuid, name, password, created_at) VALUES (?, ?, ?, ?)", tableNameUsers)
	_, err := db.Exec(cmd, u.UUID, u.Name, Encrypt(u.Password), "2021-01-01")
	if err != nil {
		return err
	}
	return err
}

/*
func (u *User) CreateTmpUser() error {
	cmd := fmt.Sprintf("INSERT INTO %s (uuid, created_at) VALUES (?, ?)", tableNameTemporaryUsers)
	_, err := db.Exec(cmd, u.UUID, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return err
}
*/

// CreateSession is
func (u *User) CreateSession() error {
	cmd := fmt.Sprintf("INSERT INTO %s (uuid, user_id, user_name, created_at) VALUES (?, ?, ?, ?)", tableNameSessions)
	_, err := db.Exec(cmd, CreateUUID(), u.UUID, u.Name, time.Now().Format(time.RFC3339))
	return err
}

// Get the session for an existing user
func GetSession(uUUID string) *Session {
	cmd := fmt.Sprintf("SELECT id, uuid, user_id, user_name, created_at FROM %s WHERE user_id = ?", tableNameSessions)
	//err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = $1", user.Id).
	row := db.QueryRow(cmd, uUUID)
	var session Session
	err := row.Scan(&session.ID, &session.UUID, &session.UserID, &session.UserName, &session.CreatedAt)
	if err != nil {
		return nil
	}
	return NewSession(session.ID, session.UUID, session.UserID, session.UserName, session.CreatedAt)
}

// Delete user from database
func (u *User) Delete() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE UUID = ?", tableNameUsers)
	//statement := "delete from users where id = $1"
	_, err := db.Exec(cmd, u.UUID)
	if err != nil {
		return err
	}
	return err
}

// Update user information in the database
func (u *User) Update() error {
	cmd := fmt.Sprintf("UPDATE %s SET name = ? WHERE uuid = ?", tableNameUsers)
	//statement := "update users set name = $2, email = $3 where id = $1"
	_, err := db.Exec(cmd, u.Name, u.UUID)
	if err != nil {
		return err
	}
	return err
}

// Delete all users from database
func UserDeleteAll() error {
	cmd := fmt.Sprintf("DELETE FROM %s", tableNameUsers)
	_, err := db.Exec(cmd)
	return err
}

/*
// Get all users in the database and returns it
func Users() (users []User, err error) {
	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}
*/

// Get a single user given the UUID
func UserBySessionID(sessionID string) *User {
	session := GetSession(sessionID)
	userID := session.UserID
	cmd := fmt.Sprintf("SELECT id, uuid, name, password, created_at FROM %s WHERE uuid = ?", tableNameUsers)
	row := db.QueryRow(cmd, userID)
	var user User
	//err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
	_ = row.Scan(&user.ID, &user.UUID, &user.Name, &user.Password, &user.CreatedAt)
	//if err != nil {
	//	return nil
	//}
	return NewUser(user.ID, user.UUID, user.Name, user.Password, user.CreatedAt)
}

func SessionByUUID(uuid string) *Session {
	cmd := fmt.Sprintf("SELECT uuid, user_id, user_name, created_at FROM %s WHERE uuid = ?", tableNameSessions)
	row := db.QueryRow(cmd, uuid)
	var session Session
	//err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
	err := row.Scan(&session.ID, &session.UUID, &session.UserID, &session.UserName, &session.CreatedAt)
	if err != nil {
		return nil
	}
	return NewSession(session.ID, session.UUID, session.UserID, session.UserName, session.CreatedAt)
}

func CheckUser(user User) (User, bool) {
	cmd := fmt.Sprintf("SELECT uuid, name FROM %s WHERE name = ?", tableNameUsers)
	row := db.QueryRow(cmd, user.Name)
	var u User
	err := row.Scan(&u.UUID, &u.Name)
	fmt.Printf("user: %v\n", u)
	if err != nil {
		return u, false
	}
	return u, true
}

func GetUserNameBySessionID(sessionID string) string {
	cmd := fmt.Sprintf("SELECT user_name FROM %s WHERE uuid = ?", tableNameSessions)
	row := db.QueryRow(cmd, sessionID)
	var userName string
	err := row.Scan(&userName)
	if err != nil {
		return ""
	}
	return userName
}

func Encrypt(password string) string {
	cryptext := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	return cryptext
}
