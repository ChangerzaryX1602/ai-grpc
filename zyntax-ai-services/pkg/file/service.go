package file

import (
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/domain"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/models"
)

type fileService struct {
	repository domain.FileRepository
}

func NewFileService(repository domain.FileRepository) domain.FileService {
	return &fileService{
		repository: repository,
	}
}
func (fs *fileService) CreateFile(file models.File) (models.File, error) {
	return fs.repository.CreateFile(file)
}
func (fs *fileService) GetFile(id string) (models.File, error) {
	return fs.repository.GetFile(id)
}
func (fs *fileService) GetFiles() ([]models.File, error) {
	files, err := fs.repository.GetFiles()
	if err != nil {
		return nil, err
	}
	return files, nil
}
func (fs *fileService) DeleteFile(id string) error {
	return fs.repository.DeleteFile(id)
}
func (fs *fileService) UpdateFile(file models.File) (models.File, error) {
	updated, err := fs.repository.UpdateFile(file)
	if err != nil {
		return models.File{}, err
	}
	return updated, nil
}
