package storage

import "ash/internal/dto"

type StorageIface interface {
	SaveData(data DataIface)
	GetTopHistoryByDirs(currentDir string, limit int) []StorageResult
	GetTopHistoryByPattern(prefix string, limit int) []StorageResult
}

type DataIface interface {
	GetExecutionList() []dto.CommandIface
	GetCurrentDir() string
}

type StorageResult interface {
	GetCommand() string
	GetDir() string
	GetUsedCount() int
}
