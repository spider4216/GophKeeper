package services

import (
	"net/http"
	"testing"

	"github.com/spider4216/GophKeeper/internal/client/repositories/reptest"
	"github.com/spider4216/GophKeeper/internal/logger"
	commonRep "github.com/spider4216/GophKeeper/internal/repository/reptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateEmailFormat(t *testing.T) {
	cli := &http.Client{}
	logger := logger.Init("debug")
	commonRep := commonRep.NewRepository(logger)
	rep := reptest.NewRepository(logger, commonRep)

	service, err := New(
		WithHTTPClient(cli),
		WithHost(""),
		WithRepo(rep),
		WithLogger(logger),
	)

	require.NoError(t, err)

	cases := []struct {
		name      string
		email     string
		exprected bool
	}{
		{
			name:      "case negative #1",
			email:     "test",
			exprected: false,
		},
		{
			name:      "case negative #2",
			email:     "example@",
			exprected: false,
		},
		{
			name:      "case negative #3",
			email:     "example@test",
			exprected: false,
		},
		{
			name:      "case positive #4",
			email:     "example@test.com",
			exprected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exprected, service.ValidateEmailFormat(tc.email))
		})
	}
}
