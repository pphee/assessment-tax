package store

import (
	"github.com/pphee/assessment-tax/store/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgresStore struct {
	DB *gorm.DB
}

func NewPostgresStore(dsn string) *PostgresStore {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error getting database connection: ", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS assessment_tax;").Error; err != nil {
		log.Fatal("Failed to create schema: ", err)
	}

	if err := db.Exec("SET search_path TO assessment_tax;").Error; err != nil {
		log.Fatal("Failed to set schema: ", err)
	}

	if err := db.Exec("CREATE TYPE assessment_tax.allowance_type AS ENUM ('PersonalDefault', 'PersonalMax', 'DonationMax', 'KReceiptDefault', 'KReceiptMax');").Error; err != nil {
		log.Fatal("Failed to create ENUM type: ", err)
	}

	if err := db.AutoMigrate(&modelgorm.AllowanceGorm{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	if err := modelgorm.InitializeData(db); err != nil {
		log.Fatal("Failed to initialize data: ", err)
	}

	return &PostgresStore{DB: db}
}
