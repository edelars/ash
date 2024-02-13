package storage

import "ash/internal/dto"

type StorageIface interface {
	SaveData(data DataIface)
	GetTopHistoryForCurrentDirAndAll(currentDir string, limit int) []StorageResult
	GetHistoryMathPrefix(prefix string, limit int) []StorageResult
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
