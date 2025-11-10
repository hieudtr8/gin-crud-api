//go:build tools
// +build tools

// This file declares dependencies for tools used in development.
// These imports ensure that `go mod tidy` keeps the tool dependencies in go.mod.
// See: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
)

// Note: For CLI tools (like ent), use `go run <package>@latest` instead of importing here
// Example: go run entgo.io/ent/cmd/ent@latest new --target internal/ent/schema Project
