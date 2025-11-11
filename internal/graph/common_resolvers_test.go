package graph

import (
	"context"
	"testing"

	"gin-crud-api/internal/database"
	"gin-crud-api/internal/middleware"
	"gin-crud-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	// Create test database client
	client := testutil.NewTestEntClient(t)
	defer client.Close()

	// Create repositories
	deptRepo := database.NewEntDepartmentRepo(client)
	empRepo := database.NewEntEmployeeRepo(client)
	projRepo := database.NewEntProjectRepo(client)

	// Create resolver with dependencies
	resolver := NewResolver(deptRepo, empRepo, projRepo)

	// Create context with request ID
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	// Call health check
	result, err := resolver.Query().Health(ctx)

	// Assert success
	require.NoError(t, err)
	assert.Equal(t, "ok", result)
}
