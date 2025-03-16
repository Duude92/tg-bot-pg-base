package Models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

func CreateDb(host string, port int, user string, password string, dbName string) error {
	var err error
	once.Do(func() {

		connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbName, port)
		if db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  connectionString, // data source name, refer https://github.com/jackc/pgx
			PreferSimpleProtocol: true,             // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
		}), &gorm.Config{}); err != nil {
			panic("failed to connect database")
		}
		if err = db.AutoMigrate(&User{}); err != nil {
			panic("failed to migrate database")
		}
	})
	return err
}
func AddUser(user *User) error {
	if db == nil {
		return gorm.ErrInvalidDB
	}
	ctx := db.FirstOrCreate(user)
	if ctx.Error != nil {
		log.Printf("Error adding user: %v", ctx.Error)
		return ctx.Error
	}
	if ctx.RowsAffected == 0 {
		UpdateUser(user)
	}
	log.Printf("New user added: %v", user)
	return nil
}
func UpdateUser(user *User) error {
	if db == nil {
		return gorm.ErrInvalidDB
	}
	ctx := db.Save(user)
	if ctx.Error != nil {
		log.Printf("Error updating user: %v", ctx.Error)
		return ctx.Error
	}
	log.Printf("Updated user: %v", user)
	return nil
}
func DeleteUser(userId int64) error {
	if db == nil {
		return gorm.ErrInvalidDB
	}
	var existingUser User
	ctx := db.First(&existingUser, "id = ?", userId)
	if (ctx.Error) != nil {
		log.Printf("No userId %d found\n%v", userId, ctx.Error)
		return ctx.Error
	}
	ctx = db.Delete(&existingUser)
	if (ctx.Error) != nil {
		log.Printf("Error deleting user: %v", ctx.Error)
		return ctx.Error
	}
	log.Printf("Deleted user: %v", userId)
	return nil
}
