package migrator

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/todo-enjoers/backend_v1/internal/config"
)

var versionTableName = "versions"

type Migrator struct {
	cfg        *config.Config
	log        *zap.Logger
	migrator   *migrate.Migrator
	migrations fs.FS
	conn       *pgx.Conn
}

func New(ctx context.Context, cfg *config.Config, log *zap.Logger, migrations fs.FS) (*Migrator, error) {
	if cfg == nil {
		return nil, fmt.Errorf("must provide config")
	}

	if log == nil {
		log = zap.NewNop()
	}

	if migrations == nil {
		return nil, fmt.Errorf("must provide migrations")
	}

	m := &Migrator{
		migrator:   nil,
		conn:       nil,
		migrations: migrations,
		cfg:        cfg,
		log:        log.Named("migrator"),
	}

	err := multierr.Combine(
		m.initConn(ctx),
		m.initMigrator(ctx),
		m.loadMigrations(),
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Migrator) MigrateUp(ctx context.Context) error {
	return m.migrator.Migrate(ctx)
}

func (m *Migrator) MigrateDown(ctx context.Context) error {
	return m.migrator.MigrateTo(ctx, 0)
}

func (m *Migrator) MigrateTo(ctx context.Context, version int32) error {
	return m.migrator.MigrateTo(ctx, version)
}

func (m *Migrator) GetVersion(ctx context.Context) (int32, error) {
	return m.migrator.GetCurrentVersion(ctx)
}

func (m *Migrator) Close(ctx context.Context) {
	if m == nil || m.conn == nil {
		return
	}

	_ = m.conn.Close(ctx)
}

func (m *Migrator) loadMigrations() (err error) {
	err = m.migrator.LoadMigrations(m.migrations)
	if err != nil {
		m.log.Error("migrator failed to load migrations", zap.Error(err))
		return fmt.Errorf("migrator.LoadMigrations: %w", err)
	}
	return nil
}

func (m *Migrator) initConn(ctx context.Context) (err error) {
	m.conn, err = pgx.Connect(ctx, m.cfg.Postgres.GetURI())
	if err != nil {
		m.log.Error("failed to connect to database", zap.Error(err))
		return fmt.Errorf("pgx.Connect: %w", err)
	}
	m.log.Debug("connected to database")
	return nil
}

func (m *Migrator) initMigrator(ctx context.Context) (err error) {
	m.migrator, err = migrate.NewMigrator(ctx, m.conn, versionTableName)
	if err != nil {
		return fmt.Errorf("migrate.NewMigrator: %w", err)
	}
	return nil
}
