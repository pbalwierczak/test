package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		user := CreateUser()

		assert.NotEqual(t, uuid.Nil, user.ID)
		// Note: CreatedAt and UpdatedAt are set by SetTimestamps method
		assert.Zero(t, user.CreatedAt) // Will be set by SetTimestamps
		assert.Zero(t, user.UpdatedAt) // Will be set by SetTimestamps
	})

	t.Run("UserIDGeneration", func(t *testing.T) {
		user1 := CreateUser()
		user2 := CreateUser()

		// Ensure each user gets a unique ID
		assert.NotEqual(t, user1.ID, user2.ID)
		assert.NotEqual(t, uuid.Nil, user1.ID)
		assert.NotEqual(t, uuid.Nil, user2.ID)
	})

	t.Run("UserTableName", func(t *testing.T) {
		user := &User{}
		tableName := user.TableName()
		assert.Equal(t, "users", tableName)
	})

	t.Run("UserStructFields", func(t *testing.T) {
		user := CreateUser()

		// Test that all expected fields are present
		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.Zero(t, user.CreatedAt)
		assert.Zero(t, user.UpdatedAt)
		assert.Nil(t, user.DeletedAt) // Should be nil initially
		assert.Nil(t, user.Trips)     // Should be nil initially
	})

	t.Run("SetTimestamps", func(t *testing.T) {
		user := &User{}
		user.SetTimestamps()

		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())
	})

	t.Run("SetID", func(t *testing.T) {
		user := &User{}
		user.SetID()

		assert.NotEqual(t, uuid.Nil, user.ID)
	})
}
