package migrator

import (
	"testing"
)

func TestStatusStruct(t *testing.T) {
	s := &Status{
		Version: 3,
		Dirty:   false,
		Pending: 2,
		Applied: 3,
		Total:   5,
	}

	if s.Version != 3 {
		t.Errorf("Version = %d, want 3", s.Version)
	}
	if s.Dirty {
		t.Error("Dirty = true, want false")
	}
	if s.Pending != 2 {
		t.Errorf("Pending = %d, want 2", s.Pending)
	}
	if s.Applied != 3 {
		t.Errorf("Applied = %d, want 3", s.Applied)
	}
	if s.Total != 5 {
		t.Errorf("Total = %d, want 5", s.Total)
	}
}

func TestMigrationInfoStruct(t *testing.T) {
	m := MigrationInfo{
		Version: 1,
		Name:    "create_users",
		Applied: true,
	}

	if m.Version != 1 {
		t.Errorf("Version = %d, want 1", m.Version)
	}
	if m.Name != "create_users" {
		t.Errorf("Name = %s, want create_users", m.Name)
	}
	if !m.Applied {
		t.Error("Applied = false, want true")
	}
}
