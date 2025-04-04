package user

import (
	"errors"
	"fmt"
	"strconv"

	"time"
	"webapp/db"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User represents the user schema.

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement;column:user_id" json:"userId"`
	CreatedAt time.Time `gorm:"column:date_created" json:"dateCreated"`
	UpdatedAt time.Time `gorm:"column:date_modified" json:"dateUpdated"`
	FirstName string    `gorm:"column:first_name" json:"firstName"`
	LastName  string    `gorm:"column:last_name" json:"lastName"`
	Username  string    `gorm:"column:user_name;unique" json:"username"`
	Password  string    `gorm:"column:password_hash" json:"passwordHash"`
	Role      string    `gorm:"column:role;default:'USER'" json:"role"`
}

// CreateUser inserts a new user record into the users table using GORM.
func CreateUser(u *User) error {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return errors.New("error setting up connection to database")
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error hashing password")
	}
	u.Password = string(hashedPass)
	if err := gormDB.Create(u).Error; err != nil {
		fmt.Println(err)
		return errors.New("error inserting new user")
	}
	defer db.CloseDB()
	return nil
}

// GetUserByID retrieves a user from the database by its ID.
func GetUserByID(userID string) (*User, error) {
	sqlDB, err := db.GetMySQLConn()

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error connecting to MySQL")
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error initializing GORM")
	}

	if err := gormDB.AutoMigrate(&User{}); err != nil {
		fmt.Println(err)
		return nil, errors.New("error migrating user schema")
	}

	// Convert the string userID to a numeric value.
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	var retrievedUser User
	result := gormDB.First(&retrievedUser, uint(id))
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, errors.New("user not found")
	}

	return &retrievedUser, nil
}

// UpdateUser updates an existing user's fields in the database.
// It uses the user's ID (which should be already set) to perform the update.
func UpdateUser(u *User) error {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return errors.New("error setting up connection to database")
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error hashing password")
	}
	u.Password = string(hashedPass)
	// Update the user record identified by u.ID.
	if err := gormDB.Model(&User{}).Where("user_id = ?", u.ID).Updates(u).Error; err != nil {
		fmt.Println(err)
		return errors.New("error updating user")
	}

	fmt.Printf("User updated: %+v\n", u)

	return nil
}

// Authenticate checks if the provided username and password match a record in the database.
func Authenticate(username, password string) (*User, error) {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return nil, errors.New("error setting up connection to database")
	}

	var u User
	// Query the user by username.
	if err := gormDB.Where("user_name = ?", username).First(&u).Error; err != nil {
		fmt.Println(err)
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &u, nil
}

func IsUserExists(username string) bool {
	gormDB, err := db.GetOrmDatabase()
	if err != nil {
		return false
	}
	var u User
	// Query the user by username.
	if err := gormDB.Where("user_name = ?", username).First(&u).Error; err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
