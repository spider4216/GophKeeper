package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptData(t *testing.T) {
	service, err := prepareService()
	require.NoError(t, err)

	key := "12345678901234567890123456789012"
	data := "Hello Gopher"

	encrypted, err := service.EncryptData([]byte(data), []byte(key))
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
}

func TestDecryptData(t *testing.T) {
	service, err := prepareService()
	require.NoError(t, err)

	key := "12345678901234567890123456789012"
	data := "Hello Gopher"

	encrypted, err := service.EncryptData([]byte(data), []byte(key))
	require.NoError(t, err)

	decrypted, err := service.DecryptData(encrypted, []byte(key))
	require.NoError(t, err)
	assert.Equal(t, data, string(decrypted))
}
