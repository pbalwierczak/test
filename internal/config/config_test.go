package config

import (
	"testing"

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
