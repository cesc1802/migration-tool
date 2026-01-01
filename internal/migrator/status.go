package migrator

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
)

// Status represents migration status info
type Status struct {
	Version uint
	Dirty   bool
	Pending int
	Applied int
	Total   int
}

// Status returns current migration status
func (mg *Migrator) Status() (*Status, error) {
	version, dirty, err := mg.m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return nil, err
	}

	pending, applied, total := mg.countMigrations(version)

	return &Status{
		Version: version,
		Dirty:   dirty,
		Pending: pending,
		Applied: applied,
		Total:   total,
	}, nil
}

// countMigrations counts pending/applied migrations relative to current version
func (mg *Migrator) countMigrations(currentVersion uint) (pending, applied, total int) {
	src := mg.sourceDriver

	v, err := src.First()
	for err == nil {
		total++
		if v <= currentVersion && currentVersion != 0 {
			applied++
		} else {
			pending++
		}
		v, err = src.Next(v)
	}

	// Handle ErrNotExist (no more migrations)
	if err != nil && err != os.ErrNotExist {
		return 0, 0, 0
	}

	return pending, applied, total
}

// MigrationInfo represents a single migration entry
type MigrationInfo struct {
	Version uint
	Name    string
	Applied bool
}

// GetMigrationList returns list of migrations with applied status
func (mg *Migrator) GetMigrationList(currentVersion uint) []MigrationInfo {
	var list []MigrationInfo
	src := mg.sourceDriver

	v, err := src.First()
	for err == nil {
		_, name, _ := src.ReadUp(v)
		list = append(list, MigrationInfo{
			Version: v,
			Name:    name,
			Applied: v <= currentVersion && currentVersion != 0,
		})
		v, err = src.Next(v)
	}

	return list
}
