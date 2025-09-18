package config

import (
	"testing"

	"scootin-aboot/pkg/database"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigLoad(t *testing.T) {
	config, err := Load()
	require.NoError(t, err)
	assert.NotEmpty(t, config.APIKey)
	assert.NotEmpty(t, config.ServerPort)
}

func TestDatabaseConnection(t *testing.T) {
	// Load test configuration
	config, err := Load()
	require.NoError(t, err)

	// Test database connection using the new database package
	dsn := config.GetDatabaseDSN()
	db, err := database.ConnectDatabase(dsn)
	require.NoError(t, err)
	require.NotNil(t, db)

	// Test that we can get the underlying sql.DB
	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NotNil(t, sqlDB)

	// Test ping
	err = sqlDB.Ping()
	require.NoError(t, err)

	// Close connection
	err = sqlDB.Close()
	require.NoError(t, err)
}

func TestDatabaseDSN(t *testing.T) {
	config := &Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "test_db",
		DBUser:     "test_user",
		DBPassword: "test_password",
		DBSSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=test_user password=test_password dbname=test_db sslmode=disable"
	actual := config.GetDatabaseDSN()
	assert.Equal(t, expected, actual)
}
