package models

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"
)

type Session struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"createdAt"`
}

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

func NewSession(id int64, uuid string, token string, userID string, createdAt string) *Session {
	return &Session{
		id,
		uuid,
		token,
		userID,
		createdAt,
	}
}

func GetUser(sessionUUID string) (*User, error) {
	cmd := fmt.Sprintf("SELECT id, uuid, name, password, created_at FROM %s WHERE uuid = ?", tableNameUsers)

	row := db.QueryRow(cmd, sessionUUID)
	var user User
	err := row.Scan(&user.ID, &user.UUID, &user.Name, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return NewUser(user.ID, user.UUID, user.Name, user.Password, user.CreatedAt), nil
}

func (u *User) Create() error {
	u.UUID = CreateUUID()
	cmd := fmt.Sprintf("INSERT INTO %s (uuid, name, password, created_at) VALUES (?, ?, ?, ?)", tableNameUsers)
	_, err := db.Exec(cmd, u.UUID, u.Name, Encrypt(u.Password), time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02T15:04:05+09:00"))
	if err != nil {
		return err
	}
	return err
}

func (u *User) CreateSession() (string, error) {
	cmd := fmt.Sprintf("INSERT INTO %s (uuid, token, user_id, created_at) VALUES (?, ?, ?, ?)", tableNameSessions)
	token := CreateUUID()
	_, err := db.Exec(cmd, CreateUUID(), token, u.UUID, time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02T15:04:05+09:00"))
	return token, err
}

func (u *User) Validate() bool {
	cmd := fmt.Sprintf("SELECT name FROM %s WHERE name = ?", tableNameUsers)
	row := db.QueryRow(cmd, u.Name)
	var user User
	err := row.Scan(&user.ID, &user.UUID, &user.Name, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return true
	} else {
		return false
	}
}

func (u *User) Delete() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE UUID = ?", tableNameUsers)
	_, err := db.Exec(cmd, u.UUID)
	if err != nil {
		return err
	}
	return err
}

func (u *User) Update() error {
	cmd := fmt.Sprintf("UPDATE %s SET name = ? WHERE uuid = ?", tableNameUsers)
	_, err := db.Exec(cmd, u.Name, u.UUID)
	if err != nil {
		return err
	}
	return err
}

func CheckUser(user User) (string, bool, bool) {
	cmd := fmt.Sprintf("SELECT uuid FROM %s WHERE name = ?", tableNameUsers)
	row := db.QueryRow(cmd, user.Name)
	var u User
	err := row.Scan(&u.UUID)
	if err == sql.ErrNoRows {
		return "", false, false
	}
	cmd = fmt.Sprintf("SELECT uuid FROM %s WHERE name = ? AND password = ?", tableNameUsers)
	row = db.QueryRow(cmd, user.Name, Encrypt(user.Password))
	err = row.Scan(&u.UUID)
	if err == sql.ErrNoRows {
		return "", true, false
	}
	return u.UUID, true, true
}

func UpdateToken(userID string) (string, error) {
	token := CreateUUID()
	cmd := fmt.Sprintf("UPDATE %s SET token = ? WHERE user_id = ?", tableNameSessions)
	_, err := db.Exec(cmd, token, userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetUserNameByToken(token string) (string, error) {
	cmd := fmt.Sprintf("SELECT users.name FROM sessions LEFT JOIN users on sessions.user_id = users.uuid WHERE token = ?")
	row := db.QueryRow(cmd, token)
	var userName string
	err := row.Scan(&userName)
	if err != nil {
		return "", err
	}
	return userName, nil
}

func Encrypt(password string) string {
	cryptext := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	return cryptext
}
