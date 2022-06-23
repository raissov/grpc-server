package delivery

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"server/config"
	"server/internal/models"
	"server/internal/service"
)

type Handler struct {
	service *service.Service
	logger  *zap.SugaredLogger
	cfg     *config.Configs
}

type TokenVerify struct {
	Token string `json:"access_token"`
}

func NewHandler(services *service.Service, logger *zap.SugaredLogger, cfg *config.Configs) *Handler {
	return &Handler{
		service: services,
		logger:  logger,
		cfg:     cfg,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/approve", h.Approve)
	router.GET("/get-messages", h.GetMessages)
	return router
}

func (h *Handler) Approve(c *gin.Context) {
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Errorf("Error occurred while reading request body: %s", err.Error())
		c.JSON(401, gin.H{
			"error": models.ErrInvalidInput.Error(),
		})
		return
	}
	var tokenVerify *TokenVerify
	err = json.Unmarshal(requestBody, &tokenVerify)
	if err != nil {
		h.logger.Errorf("Error occurred while unmarshalling request body: %s", err.Error())
		c.JSON(401, gin.H{
			"error": models.ErrInvalidInput.Error(),
		})
		return
	}

	err = h.service.VerifyToken.VerifyToken(tokenVerify.Token, "myapp:3000/verify-token")
	if err != nil {
		h.logger.Errorf("Error occurred while verifying token: %s", err.Error())
		c.JSON(401, gin.H{
			"error": models.ErrInvalidInput.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "approved",
	})
}

func (h *Handler) GetMessages(c *gin.Context) {
	messages, err := h.service.Message.GetMessages()
	if err != nil {
		h.logger.Errorf("Error occurred while getting messages: %s", err.Error())
		c.JSON(401, gin.H{
			"error": models.ErrInvalidInput.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"messages": messages,
	})
}
