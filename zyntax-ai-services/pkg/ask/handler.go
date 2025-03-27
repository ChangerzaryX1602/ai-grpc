package ask

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Zentrix-Software-Hive/zyntax-ai-services/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type askHandler struct {
	client AiServiceClient // Assuming your service is named AIService
	db     *gorm.DB
}

func NewAskHandler(askRoute fiber.Router, auth *handlers.RouterResources, grpcAddress string, db *gorm.DB) {
	handler, err := newAskHandler(grpcAddress, db)
	if err != nil {
		log.Fatalf("Failed to create ask handler: %v", err)
	}
	askRoute.Post("/", auth.ReqAuthHandler(0), handler.Ask)
	askRoute.Post("/history", auth.ReqAuthHandler(0), handler.CreateHistoryMe)
	askRoute.Get("/history", auth.ReqAuthHandler(0), handler.GetHistoriesMe)
	askRoute.Get("/history/messages/:id", auth.ReqAuthHandler(0), handler.GetHistoryMessageByHistoryID)
	askRoute.Put("/history/:id", auth.ReqAuthHandler(0), handler.UpdatePlaceholder)
	askRoute.Delete("/history/:id", auth.ReqAuthHandler(0), handler.DeleteHistory)
}

func newAskHandler(grpcAddress string, db *gorm.DB) (*askHandler, error) {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	// defer conn.Close() // Don't defer here, we need to keep the connection alive

	client := NewAiServiceClient(conn)
	return &askHandler{client: client, db: db}, nil
}

type Ask struct {
	Question  string `json:"question"`
	HistoryId int    `json:"history_id"`
}

type History struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement" swaggerignore:"true"`
	PlaceHolder string `json:"place_holder"`
}
type MapHistoryMessage struct {
	ID        int            `json:"id" gorm:"primaryKey;autoIncrement"`
	HistoryId int            `json:"history_id"`
	MessageId int            `json:"message_id"`
	History   History        `json:"history" gorm:"foreignKey:HistoryId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Message   HistoryMessage `json:"message" gorm:"foreignKey:MessageId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
type HistoryMessage struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
}
type MapUserHistory struct {
	ID         int      `json:"id" gorm:"primaryKey;autoIncrement"`
	MainUserID string   `json:"main_user_id"`
	HistoryID  int      `json:"history_id"`
	MainUser   MainUser `json:"main_user" gorm:"foreignKey:MainUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	History    History  `json:"history" gorm:"foreignKey:HistoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
type MainUser struct {
	ID     string `json:"id" gorm:"primaryKey"`
	NameTh string `json:"name_th"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
}

// @Summary Ask a question
// @Description Ask a question to the AI service
// @Tags Ask
// @Accept json
// @Produce json
// @Param question body Ask true "Question to ask"
// @Router /api/v1/ask/ [post]
// @Security ApiKeyAuth
func (h *askHandler) Ask(c *fiber.Ctx) error {
	ask := Ask{}
	err := c.BodyParser(&ask)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid request body")
	}
	if ask.Question == "" {
		return c.Status(http.StatusBadRequest).SendString("Query parameter 'question' is required")
	}
	isFeeList := []string{
		"ค่าธรรมเนียม", // fee
		"ค่าเทอม",      // tuition fee
		"ค่าเล่าเรียน", // tuition
		"ค่าลงทะเบียน", // registration fee
		"fee",
		"tuition fee",
		"tuition",
		"registration fee",
	}
	for _, fee := range isFeeList {
		if strings.Contains(strings.ToLower(ask.Question), fee) {
			//save message
			historyMessage := HistoryMessage{
				Question:  ask.Question,
				Answer:    "https://img2.pic.in.th/pic/.-650x900.png",
				CreatedAt: time.Now(),
			}
			if err := h.db.Create(&historyMessage).Error; err != nil {
				log.Printf("could not save history: %v", err)
				return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error saving history: %v", err))
			}
			//save map
			mapHistoryMessage := MapHistoryMessage{
				HistoryId: ask.HistoryId,
				MessageId: historyMessage.ID,
			}
			if err := h.db.Create(&mapHistoryMessage).Error; err != nil {
				log.Printf("could not save map: %v", err)
				return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error saving map: %v", err))
			}
			return c.SendString("https://img2.pic.in.th/pic/.-650x900.png")
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Increased timeout to 5 seconds
	defer cancel()

	// Create the request message.
	req := &AiRequest{ // Assuming your request message is named AskRequest
		Question: ask.Question,
	}

	// Call the gRPC method.
	r, err := h.client.Ask(ctx, req) // Assuming your method is named Ask
	if err != nil {
		log.Printf("could not ask: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error calling gRPC: %v", err))
	}
	//save message
	historyMessage := HistoryMessage{
		Question:  ask.Question,
		Answer:    r.GetAnswer(),
		CreatedAt: time.Now(),
	}
	if err := h.db.Create(&historyMessage).Error; err != nil {
		log.Printf("could not save history: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error saving history: %v", err))
	}
	//save map
	mapHistoryMessage := MapHistoryMessage{
		HistoryId: ask.HistoryId,
		MessageId: historyMessage.ID,
	}
	if err := h.db.Create(&mapHistoryMessage).Error; err != nil {
		log.Printf("could not save map: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error saving map: %v", err))
	}
	return c.SendString(r.GetAnswer()) // Assuming your response message has a field named Answer
}

// @Summary Create a history
// @Description Create a history
// @Tags History
// @Accept json
// @Produce json
// @Param history body History true "History to create"
// @Router /api/v1/ask/history [post]
// @Security ApiKeyAuth
func (h *askHandler) CreateHistoryMe(c *fiber.Ctx) error {
	history := History{}
	err := c.BodyParser(&history)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid request body")
	}
	if err := h.db.Create(&history).Error; err != nil {
		log.Printf("could not save history: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error saving history: %v", err))
	}
	mapUserHistory := MapUserHistory{
		MainUserID: c.Locals("user_id").(string),
		HistoryID:  history.ID,
	}
	if err := h.db.Create(&mapUserHistory).Error; err != nil {
		log.Printf("could not save map: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error saving map: %v", err))
	}
	return c.JSON(history)
}

// @Summary Get all histories
// @Description Get all histories
// @Tags History
// @Accept json
// @Produce json
// @Router /api/v1/ask/history [get]
// @Security ApiKeyAuth
func (h *askHandler) GetHistoriesMe(c *fiber.Ctx) error {
	var histories []MapUserHistory
	if err := h.db.Preload(clause.Associations).Where("main_user_id = ?", c.Locals("user_id").(string)).Find(&histories).Error; err != nil {
		log.Printf("could not get histories: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error getting histories: %v", err))
	}

	return c.JSON(histories)
}

// @Summary Update a history
// @Description Update a history
// @Tags History
// @Accept json
// @Produce json
// @Param id path string true "History ID"
// @Param history body History true "History to update"
// @Router /api/v1/ask/history/{id} [put]
// @Security ApiKeyAuth
func (h *askHandler) UpdatePlaceholder(c *fiber.Ctx) error {
	history := History{}
	historyId := c.Params("id")
	historyInt, err := strconv.Atoi(historyId)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid history ID")
	}
	err = c.BodyParser(&history)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid request body")
	}
	if err := h.db.Model(&history).Where("id = ?", historyInt).Update("place_holder", history.PlaceHolder).Error; err != nil {
		log.Printf("could not update history: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error updating history: %v", err))
	}
	return c.JSON(history)
}

// @Summary Delete a history
// @Description Delete a history
// @Tags History
// @Accept json
// @Produce json
// @Param id path string true "History ID"
// @Router /api/v1/ask/history/{id} [delete]
// @Security ApiKeyAuth
func (h *askHandler) DeleteHistory(c *fiber.Ctx) error {
	historyID := c.Params("id")
	if err := h.db.Delete(&MapUserHistory{}, "history_id = ?", historyID).Error; err != nil {
		log.Printf("could not delete)} history: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error deleting history: %v", err))
	}
	return c.SendStatus(http.StatusOK)

}

// @Summary Get history messages by history ID
// @Description Get history messages by history ID
// @Tags History
// @Accept json
// @Produce json
// @Param id path string true "History ID"
// @Router /api/v1/ask/history/messages/{id} [get]
// @Security ApiKeyAuth
func (h *askHandler) GetHistoryMessageByHistoryID(c *fiber.Ctx) error {
	historyID := c.Params("id")
	var historyMessages []MapHistoryMessage
	if err := h.db.Preload(clause.Associations).Where("history_id = ?", historyID).Find(&historyMessages).Error; err != nil {
		log.Printf("could not get history messages: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Error getting history messages: %v", err))
	}
	return c.JSON(historyMessages)
}
