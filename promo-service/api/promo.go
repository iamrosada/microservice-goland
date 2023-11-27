package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/entity"
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/usecase"

	"github.com/gin-gonic/gin"
)

var userMicroserviceURL = "http://localhost:8000"

type PromoHandlers struct {
	PromoUseCase *usecase.PromotionUsecaseImpl
}

func NewPromoHandlers(useCase *usecase.PromotionUsecaseImpl) *PromoHandlers {
	return &PromoHandlers{
		PromoUseCase: useCase,
	}
}

func (p *PromoHandlers) SetupRouter(router *gin.Engine) {
	// Use the provided router instead of creating a new one
	// r := gin.Default() // Don't create a new router; use the provided one

	// Health check endpoint
	router.GET("/check_alive", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Uncommented routes
	router.POST("/promo/create", p.CreatePromoHandler)
	router.POST("/promo/:id/codes", p.PromoCodeHandler)
	// router.POST("/promo/:id/apply_users", p.ApplyPromoUsersHandler)
	// router.GET("/promo/:id/users", p.GetPromoUsersHandler)

	// Active route with your ApplyPromoAllHandler
	router.POST("/promo/:id/apply_all", p.ApplyPromoAllHandler)
}

func fetchUsersForPromo(url string) (map[string][]uint, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch users for promo, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string][]uint
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *PromoHandlers) ApplyPromoAllHandler(c *gin.Context) {
	id := c.Param("id")
	uInt32Val, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var request struct {
		UserIDs []int `json:"user_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the user microservice to apply the promo to all users
	url := fmt.Sprintf("%s/users/promo_type/%d/available", userMicroserviceURL, uInt32Val)
	usersIDs, err := fetchUsersForPromo(url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, userID := range usersIDs["users_id"] {
		randomNumber := rand.Intn(9) + 1
		appliedPromotion := entity.UserPromotion{
			PromotionID: uint(uInt32Val),
			UserID:      uint(userID),
			Type:        int(randomNumber),
		}

		// Save the applied promotion to the user microservice
		saveURL := fmt.Sprintf("%s/api/users/promo/%d/apply", userMicroserviceURL, uInt32Val)
		reqBody, _ := json.Marshal(appliedPromotion)
		resp, err := http.Post(saveURL, "application/json", bytes.NewBuffer(reqBody))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save applied promotion"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promo applied to all users"})
}

func (p *PromoHandlers) PromoCodeHandler(c *gin.Context) {
	// Parse and validate promoID
	promoID := c.Param("id")
	uInt32Val, err := strconv.ParseUint(promoID, 10, 32)

	// Validate other input parameters as needed

	// Log the received request
	log.Printf("Received request to add codes to promotion with ID %d", uInt32Val)

	// Bind JSON request
	var request struct {
		CodesPromo []string `json:"codes"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	// Check for zero-length codes
	if len(request.CodesPromo) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No codes provided"})
		return
	}

	// Convert codes to lowercase for case-insensitive duplicate check
	lowercaseCodes := make(map[string]bool)
	for _, code := range request.CodesPromo {
		lowercaseCode := strings.ToLower(code)
		if _, exists := lowercaseCodes[lowercaseCode]; exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate code found", "code": code})
			return
		}
		lowercaseCodes[lowercaseCode] = true
	}

	// Handle the use case
	err = p.PromoUseCase.AddCodesToPromotion(uint(uInt32Val), request.CodesPromo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply promotion", "details": err.Error()})
		return
	}

	// Respond to the client
	c.JSON(http.StatusOK, gin.H{"message": "Promotion applied successfully"})
}

func (p *PromoHandlers) CreatePromoHandler(c *gin.Context) {
	var request struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Type        int    `json:"type"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPromotion := entity.Promotion{
		Name:        request.Name,
		Slug:        request.Slug,
		URL:         request.URL,
		Description: request.Description,
		Type:        request.Type,
	}

	err := p.PromoUseCase.CreatePromo(&newPromotion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create promotion", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion created successfully"})
}
