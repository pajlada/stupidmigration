package stupidmigration

import (
	"database/sql"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func createMigrationsTableIfItDoesNotAlreadyExist(db *sql.DB) error {
	const createTableQuery = `CREATE TABLE IF NOT EXISTS migrations(version TEXT);`
	_, err := db.Exec(createTableQuery)
	return err
}

func insertVersionRow(db *sql.DB) error {
	const query = `INSERT INTO migrations (version) VALUES ('0')`
	_, err := db.Exec(query)
	return err
}

// getCurrentVersion returns the current version from the database. if no version row is there, it will insert a new row with the default value 0
func getCurrentVersion(db *sql.DB) (uint64, error) {
	const query = `SELECT version FROM migrations;`
	row := db.QueryRow(query)
	var versionString string

	if err := row.Scan(&versionString); err != nil {
		if err == sql.ErrNoRows {
			err = insertVersionRow(db)
		}
		return 0, err
	}

	return strconv.ParseUint(versionString, 10, 64)
}

func getMigrations(migrationsPath string, version uint64) (migrations []migration, err error) {
	files, err := filepath.Glob(migrationsPath + "/*.sql")
	if err != nil {
		return
	}

	for _, f := range files {
		b := path.Base(f)
		parts := strings.Split(b, "-")
		migrationVersion, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		if version >= migrationVersion {
			// This migration is older than our database version (which means we've already applied it). Ignore it
			continue
		}
		migrations = append(migrations, migration{
			path:    f,
			version: migrationVersion,
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	return
}

// Migrate attempts to load all migrations from migrationsPath and apply them to the given db
func Migrate(migrationsPath string, db *sql.DB) error {
	if err := createMigrationsTableIfItDoesNotAlreadyExist(db); err != nil {
		return err
	}

	version, err := getCurrentVersion(db)
	if err != nil {
		return err
	}

	migrations, err := getMigrations(migrationsPath, version)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if err := migration.Migrate(db); err != nil {
			return err
		}
		if migration.err != nil {
			return err
		}
	}

	return nil
}
