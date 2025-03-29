package file

import (
	"os"
	"strconv"

	"github.com/Zentrix-Software-Hive/zyntax-ai-services/internal/handlers"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/domain"
	"github.com/Zentrix-Software-Hive/zyntax-ai-services/pkg/models"
	"github.com/gofiber/fiber/v2"
	helpers "github.com/zercle/gofiber-helpers"
)

type fileHandler struct {
	handlers.RouterResources
	service domain.FileService
}

func NewFileHandler(router fiber.Router, service domain.FileService) {
	handler := &fileHandler{
		service: service,
	}
	router.Get("/", handler.GetFiles())
	router.Get("/:id", handler.GetFile())
	router.Post("", handler.CreateFile())
	router.Put("/:id", handler.UpdateFile())
	router.Delete("/:id", handler.DeleteFile())
}

// @Summary Get all files
// @Description Get all files
// @Tags files
// @Accept json
// @Produce json
// @Router /api/v1/files [get]
// @Security ApiKeyAuth
func (h *fileHandler) GetFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		files, err := h.service.GetFiles()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(files)
	}
}

// @Summary Get a file by ID
// @Description Get a file by ID
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Router /api/v1/files/{id} [get]
// @Security ApiKeyAuth
func (h *fileHandler) GetFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		file, err := h.service.GetFile(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(file)
	}
}

// @Summary Create a new file
// @Description Create a new file
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Router /api/v1/files [post]
// @Security ApiKeyAuth
func (h *fileHandler) CreateFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseForm{
				Errors: []helpers.ResponseError{
					{
						Code:    fiber.StatusBadRequest,
						Message: err.Error(),
					},
				},
			})
		}
		err = os.WriteFile("../ai/assets/"+file.Filename, []byte(""), 0644)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var fileType models.File
		fileType.Name = file.Filename
		fileType.Path = "../ai/assets/" + file.Filename
		createdFile, err := h.service.CreateFile(fileType)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(createdFile)
	}
}

// @Summary Update a file
// @Description Update a file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Param file body models.File true "File object"
// @Router /api/v1/files/{id} [put]
// @Security ApiKeyAuth
func (h *fileHandler) UpdateFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var file models.File
		if err := c.BodyParser(&file); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}
		file.ID = idInt
		updatedFile, err := h.service.UpdateFile(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(updatedFile)
	}
}

// @Summary Delete a file
// @Description Delete a file
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Router /api/v1/files/{id} [delete]
// @Security ApiKeyAuth
func (h *fileHandler) DeleteFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := h.service.DeleteFile(id); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
