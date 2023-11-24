package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iamrosada/microservice-goland/user-service/internal/user/usecase"
)

type UserPromoHandlers struct {
	UserPromoUseCase *usecase.PromotionUseCase
}

func NewUserPromoHandlers(useCase *usecase.PromotionUseCase) *UserPromoHandlers {
	return &UserPromoHandlers{
		UserPromoUseCase: useCase,
	}
}

func (p *UserPromoHandlers) SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/promo_type/:type/available", p.GetAvailableUsersHandler)
			users.GET("/promo/:id/applied", p.GetAppliedUsersHandler)
			users.POST("/promo/:id/apply", p.ApplyPromotionHandler)
		}
	}
}

func (p *UserPromoHandlers) GetAvailableUsersHandler(c *gin.Context) {
	promoType := c.Param("type")
	println(promoType)
	Int32Val, _ := strconv.ParseInt(promoType, 10, 32)

	// Validate promoType and other input as needed

	users, err := p.UserPromoUseCase.GetAvailableUsers(int(Int32Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_ids": users})
}

func (p *UserPromoHandlers) GetAppliedUsersHandler(c *gin.Context) {
	promoID := c.Param("id")
	uInt32Val, err := strconv.ParseUint(promoID, 10, 32)

	// Validate promoID and other input as needed

	users, err := p.UserPromoUseCase.GetAppliedUsers(uint(uInt32Val))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_ids": users})
}

func (p *UserPromoHandlers) ApplyPromotionHandler(c *gin.Context) {
	promoID := c.Param("id")
	uInt32Val, err := strconv.ParseUint(promoID, 10, 32)

	// Validate promoID and other input as needed

	var request struct {
		UserIDs []int `json:"user_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = p.UserPromoUseCase.ApplyPromotion(uint(uInt32Val), request.UserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion applied successfully"})
}
