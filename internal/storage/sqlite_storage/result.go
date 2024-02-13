package sqlite_storage

import "database/sql"

type storageResult struct {
	LastUsedTime sql.NullInt64  `db:"lastUsedTime"`
	usedCounter  sql.NullInt16  `db:"usedCounter"`
	dir          sql.NullString `db:"dir"`
	execWithArgs sql.NullString `db:"execWithArgs"`
}

func (s *storageResult) GetCommand() string {
	if s.execWithArgs.Valid {
		return s.execWithArgs.String
	}
	return ""
}

func (s *storageResult) GetDir() string {
	if s.dir.Valid {
		return s.dir.String
	}
	return ""
}

func (s *storageResult) GetUsedCount() int {
	if s.usedCounter.Valid {
		return int(s.usedCounter.Int16)
	}
	return 0
}
