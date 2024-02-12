package storage

import "ash/internal/dto"

type StorageIface interface {
	SaveData(data DataIface)
}

type DataIface interface {
	GetExecutionList() []dto.CommandIface
	GetCurrentDir() string
}
