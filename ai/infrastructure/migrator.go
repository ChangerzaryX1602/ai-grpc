package infrastructure

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	checkAndAutoMigrate(db)
}
func checkAndAutoMigrate(db *gorm.DB, model ...interface{}) {
	for _, m := range model {
		// Check if the table does not exist
		if !db.Migrator().HasTable(m) {
			// Auto migrate the table if it does not exist
			if err := db.AutoMigrate(m); err != nil {
				log.Fatalf("Failed to auto migrate table: %v", err)
			}
			fmt.Printf("Table %T created successfully.\n", m)
		} else {
			fmt.Printf("Table %T already exists.\n", m)
		}
	}
}
