package db

import (
	"fmt"

	"github.com/orders-service/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDBConnection(config *app.Configuration) (*gorm.DB, error) {
	databaseURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName,
	)

	// gorm.Open() : responsible for connecting with db
	database, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return database, err
	}
	database.AutoMigrate(&Customer{})
	database.AutoMigrate(&Product{})
	database.AutoMigrate(&Order{})
	database.AutoMigrate(&OrderProduct{})
	return database, nil

}

// func LoadSampleData(dbConn *gorm.DB) {
// 	// Read the contents of the .sql file
// 	scriptFile := "./loadSampleData.sql"
// 	scriptBytes, err := ioutil.ReadFile(scriptFile)
// 	if err != nil {
// 		log.Fatal("Failed to read script file:", err)
// 	}

// 	script := string(scriptBytes)

// 	// Execute the SQL script
// 	err = dbConn.Exec(script).Error
// 	if err != nil {
// 		log.Fatal("Failed to execute script:", err)
// 	}
// 	fmt.Println("Script executed successfully.")
// }
