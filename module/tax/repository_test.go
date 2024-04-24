package tax

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pphee/assessment-tax/store/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New() // Create instance of sqlmock
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:                 mockDB,
		PreferSimpleProtocol: true, // Disables implicit prepared statement usage
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestGetAllowanceConfig(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTaxRepository(db)

	expected := []model.AllowanceGorm{{AllowanceType: "personalDeduction", Amount: 50000}}
	rows := sqlmock.NewRows([]string{"allowance_type", "amount"}).
		AddRow("personalDeduction", 50000)

	mock.ExpectQuery(`SELECT \* FROM "allowance_gorms"`).WillReturnRows(rows)

	result, err := repo.GetAllowanceConfig()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet()) // Check if all expectations were met
}

func TestSetPersonalDeduction(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTaxRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "allowance_gorms"`).
		WithArgs("personalDeduction", float64(30000)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.SetPersonalDeduction(30000)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetKreceiptDeduction(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTaxRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "allowance_gorms"`).
		WithArgs("kReceipt", float64(15000)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.SetKreceiptDeduction(15000)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllowanceConfig_QueryError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTaxRepository(db)

	mock.ExpectQuery(`SELECT \* FROM "allowance_gorms"`).
		WillReturnError(gorm.ErrInvalidData)

	result, err := repo.GetAllowanceConfig()

	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrInvalidData)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetKreceiptDeduction_SaveError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTaxRepository(db)

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "allowance_gorms"`).
		WithArgs("kReceipt", float64(15000)).
		WillReturnError(gorm.ErrInvalidDB) // Simulate a database error

	mock.ExpectRollback()

	err := repo.SetKreceiptDeduction(15000)

	assert.ErrorIs(t, err, gorm.ErrInvalidDB)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetPersonalDeduction_SaveError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTaxRepository(db)

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "allowance_gorms"`).
		WithArgs("personalDeduction", float64(30000)).
		WillReturnError(gorm.ErrInvalidDB) // simulate an error

	mock.ExpectRollback()

	err := repo.SetPersonalDeduction(30000)

	assert.ErrorIs(t, err, gorm.ErrInvalidDB)
	assert.NoError(t, mock.ExpectationsWereMet())
}
