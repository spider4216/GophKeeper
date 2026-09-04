package services

import (
	"encoding/json"
	"testing"

	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptData(t *testing.T) {
	store := map[string][]string{}
	service, err := prepareService(store)
	require.NoError(t, err)

	key := "12345678901234567890123456789012"
	data := "Hello Gopher"

	encrypted, err := service.EncryptData([]byte(data), []byte(key))
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted)
}

func TestDecryptData(t *testing.T) {
	store := map[string][]string{}
	service, err := prepareService(store)
	require.NoError(t, err)

	key := "12345678901234567890123456789012"
	data := "Hello Gopher"

	encrypted, err := service.EncryptData([]byte(data), []byte(key))
	require.NoError(t, err)

	decrypted, err := service.DecryptData(encrypted, []byte(key))
	require.NoError(t, err)
	assert.Equal(t, data, string(decrypted))
}

func TestCreateUserPassItem(t *testing.T) {
	store := map[string][]string{}
	service, err := prepareService(store)
	require.NoError(t, err)

	login := "test_login"
	pass := "qwerty"
	title := "Website"

	req := models.LoginPassReq{
		Login: login,
		Pass:  pass,
		Title: title,
	}

	key := "12345678901234567890123456789012"

	err = service.CreateUserPassItem(t.Context(), req, key, 1)
	require.NoError(t, err)

	items, ok := store["items"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(items))

	meta, ok := store["metadata"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(meta))

	pends, ok := store["pending_changes"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(pends))

	itemsMap := map[string]any{}

	err = json.Unmarshal([]byte(items[0]), &itemsMap)
	require.NoError(t, err)

	ciphertext, ok := itemsMap["ciphertext"]
	assert.Equal(t, true, ok)

	txt, ok := ciphertext.(string)

	assert.Equal(t, true, ok)

	origin, err := service.DecryptData(txt, []byte(key))

	require.NoError(t, err)

	var loginpass models.LoginPassFmt

	err = json.Unmarshal(origin, &loginpass)
	require.NoError(t, err)

	assert.Equal(t, login, loginpass.Login)
	assert.Equal(t, pass, loginpass.Pass)

	metaMap := map[string]any{}

	err = json.Unmarshal([]byte(meta[0]), &metaMap)
	require.NoError(t, err)

	val, ok := metaMap["value"]
	assert.Equal(t, true, ok)

	assert.Equal(t, title, val)

}
