package sqlite_storage

import "database/sql"

type storageResult struct {
	LastUsedTime sql.NullInt64  `db:"lastUsedTime"`
	UsedCounter  sql.NullInt16  `db:"usedCounter"`
	Dir          sql.NullString `db:"dir"`
	ExecWithArgs sql.NullString `db:"execWithArgs"`
}

func (s *storageResult) GetCommand() string {
	if s.ExecWithArgs.Valid {
		return s.ExecWithArgs.String
	}
	return ""
}

func (s *storageResult) GetDir() string {
	if s.Dir.Valid {
		return s.Dir.String
	}
	return ""
}

func (s *storageResult) GetUsedCount() int {
	if s.UsedCounter.Valid {
		return int(s.UsedCounter.Int16)
	}
	return 0
}
