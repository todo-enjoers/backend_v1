package tern

import (
	"context"
	"io/fs"
)

type Tern interface {
	Migrate(ctx context.Context) error
	MigrateTo(ctx context.Context, version int32) error
	GetCurrentVersion(ctx context.Context) (int32, error)
	LoadMigrations(migrations fs.FS) error
}
