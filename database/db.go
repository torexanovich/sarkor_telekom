package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func InitDB() *sql.DB {
	var err error
	db, err = sql.Open("sqlite3", "sarkor1.db")
	if err != nil {
		panic(err)
	}

	// Creating tables
	query := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT,
		password TEXT,
		name TEXT,
		age INTEGER
	);
	CREATE TABLE IF NOT EXISTS phones (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		phone TEXT,
		description TEXT,
		is_mobile INTEGER
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func CreateUser(login, password, name string, age int) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO users(login, password, name, age) VALUES(?, ?, ?, ?)`, login, string(hashedPassword), name, age)
	if err != nil {
		return err
	}
	fmt.Printf("New user: login=%s, password=%s, name=%s, age=%d\n", login, password, name, age)

	return nil
}

func AuthenticateUser(login, password string) (string, int, error) {
	var (
		savedPassword string
		userID        int
	)
	err := db.QueryRow(`SELECT password, id FROM users WHERE login = ?`, login).Scan(&savedPassword, &userID)
	if err != nil {
		return "", 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(password))
	if err != nil {
		return "", 0, err
	}

	jwtHandler := JWTHandler{
		SigninKey: "123",
	}

	token, err := jwtHandler.GenerateToken(login, userID) // Include the userID in the token
	if err != nil {
		return "", 0, err
	}

	return token, userID, nil
}

type User struct {
	ID   int    `json:"user_id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func GetUser(name string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, name, age FROM users WHERE name = ?", name).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreatePhoneNumber(phone, description string, isMobile bool, userID int) error {
	_, err := db.Exec("INSERT INTO phones(user_id, phone, description, is_mobile) VALUES (?, ?, ?, ?)", userID, phone, description, isMobile)
	if err != nil {
		return err
	}

	return nil
}

func GetUsersByPhoneNumber(phone string) ([]*User, error) {
	rows, err := db.Query("SELECT u.id, u.name, u.age FROM users u JOIN phones p ON u.id = p.user_id WHERE p.phone = ?", phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func UpdatePhoneNumber(phoneID int, phone, description string, isMobile bool, userID int) error {
	_, err := db.Exec("UPDATE phones SET phone = ?, description = ?, is_mobile = ? WHERE id = ? AND user_id = ?", phone, description, isMobile, phoneID, userID)
	if err != nil {
		return err
	}

	return nil
}

func DeletePhoneNumber(phoneID, userID int) error {
	_, err := db.Exec("DELETE FROM phones WHERE id = ? AND user_id = ?", phoneID, userID)
	if err != nil {
		return err
	}

	return nil
}

type JWTHandler struct {
	SigninKey string
}

func (jwtHandler *JWTHandler) GenerateToken(login string, userID int) (string, error) {
	claims := jwt.MapClaims{
		"login":   login,
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // token expiration time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the signing key
	signedToken, err := token.SignedString([]byte(jwtHandler.SigninKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (jwtHandler *JWTHandler) ExtractClaims(tokenStr string, signingKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, err
	}

	return claims, nil
}

func GetCurrentUserID(c *gin.Context) (int, error) {
	tokenStr, err := c.Cookie("SESSTOKEN") // Getting the token from the cookie
	if err != nil {
		return 0, err
	}

	jwtHandler := JWTHandler{
		SigninKey: "123",
	}

	claims, err := jwtHandler.ExtractClaims(tokenStr, jwtHandler.SigninKey)
	if err != nil {
		return 0, err
	}

	userIDF, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid userID")
	}

	userID := int(userIDF)

	return userID, nil
}
