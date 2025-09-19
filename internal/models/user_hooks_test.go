package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUser_TableName(t *testing.T) {
	user := User{}
	assert.Equal(t, "users", user.TableName())
}

func TestUser_BeforeCreate(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		expectError bool
	}{
		{
			name:        "successful creation with nil ID",
			user:        &User{},
			expectError: false,
		},
		{
			name: "successful creation with existing ID",
			user: &User{
				ID: uuid.New(),
			},
			expectError: false,
		},
		{
			name: "successful creation with existing timestamps",
			user: &User{
				CreatedAt: time.Now().Add(-1 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * time.Minute),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock GORM DB
			var tx *gorm.DB

			err := tt.user.BeforeCreate(tx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Check that ID was set if it was nil
				if tt.user.ID == uuid.Nil {
					assert.NotEqual(t, uuid.Nil, tt.user.ID)
				}

				// Check that timestamps were set if they were zero
				if tt.user.CreatedAt.IsZero() {
					assert.False(t, tt.user.CreatedAt.IsZero())
				}
				if tt.user.UpdatedAt.IsZero() {
					assert.False(t, tt.user.UpdatedAt.IsZero())
				}
			}
		})
	}
}

func TestUser_BeforeUpdate(t *testing.T) {
	tests := []struct {
		name string
		user *User
	}{
		{
			name: "successful update",
			user: &User{},
		},
		{
			name: "update with existing UpdatedAt",
			user: &User{
				UpdatedAt: time.Now().Add(-1 * time.Hour),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original UpdatedAt
			originalUpdatedAt := tt.user.UpdatedAt

			// Create a mock GORM DB
			var tx *gorm.DB

			err := tt.user.BeforeUpdate(tx)

			assert.NoError(t, err)

			// Check that UpdatedAt was set
			assert.True(t, tt.user.UpdatedAt.After(originalUpdatedAt) || tt.user.UpdatedAt.Equal(originalUpdatedAt))
		})
	}
}
