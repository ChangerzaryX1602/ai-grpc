package domain

import "github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/models"

type FileRepository interface {
	CreateFile(file models.File) (models.File, error)
	GetFile(id string) (models.File, error)
	GetFiles() ([]models.File, error)
	DeleteFile(id string) error
	UpdateFile(file models.File) (models.File, error)
}
type FileService interface {
	CreateFile(file models.File) (models.File, error)
	GetFile(id string) (models.File, error)
	GetFiles() ([]models.File, error)
	DeleteFile(id string) error
	UpdateFile(file models.File) (models.File, error)
}
