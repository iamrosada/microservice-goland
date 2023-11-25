package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

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

func (p *PromoHandlers) setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/check_alive", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// r.POST("/promo/create", createPromo)
	// r.GET("/promo/:id/codes", getPromo)
	r.POST("/promo/:id/apply_all", p.applyPromoAll)
	// r.POST("/promo/:id/apply_users", applyPromoUsers)
	// r.GET("/promo/:id/users", getPromoUsers)

	return r
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

func (p *PromoHandlers) applyPromoAll(c *gin.Context) {
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
