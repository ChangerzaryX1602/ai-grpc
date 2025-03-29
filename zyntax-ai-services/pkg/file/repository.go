package file

import (
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/domain"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/models"
	"gorm.io/gorm"
)

type fileRepository struct {
	*gorm.DB
}

func NewFileRepository(db *gorm.DB) domain.FileRepository {
	return &fileRepository{db}
}
func (r *fileRepository) CreateFile(file models.File) (models.File, error) {
	if err := r.Create(&file).Error; err != nil {
		return models.File{}, err
	}
	return file, nil
}
func (r *fileRepository) GetFile(id string) (models.File, error) {
	var file models.File
	if err := r.First(&file, id).Error; err != nil {
		return models.File{}, err
	}
	return file, nil
}
func (r *fileRepository) GetFiles() ([]models.File, error) {
	var files []models.File
	if err := r.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
func (r *fileRepository) DeleteFile(id string) error {
	if err := r.Delete(&models.File{}, id).Error; err != nil {
		return err
	}
	return nil
}
func (r *fileRepository) UpdateFile(file models.File) (models.File, error) {
	if err := r.Save(&file).Error; err != nil {
		return models.File{}, err
	}
	return file, nil
}
