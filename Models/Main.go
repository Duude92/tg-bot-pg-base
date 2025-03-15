package Models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func CreateDb(host string, port int, user string, password string, dbName string) {
	var err error
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbName, port)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  connectionString, // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,             // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database")
	}
}
func AddUser(user *User) {
	ctx := db.FirstOrCreate(user)
	if ctx.RowsAffected == 0 {
		UpdateUser(user)
	}
}
func UpdateUser(user *User) {
	db.Save(user)
}
func DeleteUser(userId int64) {
	var existingUser User
	ctx := db.First(&existingUser, "id = ?", userId)
	if (ctx.Error) != nil {
		return
	}
	ctx = db.Delete(&existingUser)
}
