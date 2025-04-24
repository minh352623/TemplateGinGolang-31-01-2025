package controller

import (
	"database/sql"
	"ecom/global"
	"ecom/internal/database"
	"ecom/internal/messaging"
	"ecom/internal/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TestController struct {
	testService service.ITestService
}

func NewTestController(testService service.ITestService) *TestController {
	return &TestController{
		testService: testService,
	}
}

func (c *TestController) GetTestById(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println("id", id)
	idUUID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	test, err := c.testService.GetTestById(idUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get test"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": test})
}

// test update
func (c *TestController) UpdateTest(ctx *gin.Context) {
	messagesUpdate := []messaging.BodyMessage{
		{
			Action: "update",
			Data: database.UpdateTestParams{
				ID:      uuid.MustParse("82f048f3-e760-44ca-b5f8-067238a52ef6"),
				Name:    "test",
				Balance: sql.NullString{String: "100", Valid: true},
			},
		},
		{
			Action: "update",
			Data: database.UpdateTestParams{
				ID:      uuid.MustParse("82f048f3-e760-44ca-b5f8-067238a52ef6"),
				Name:    "test",
				Balance: sql.NullString{String: "100", Valid: true},
			},
		},
	}
	for _, message := range messagesUpdate {
		go func() {
			// params := message.Data.(database.UpdateTestParams)
			// c.testService.UpdateTest(&params)
			messageJson, err := json.Marshal(message)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal message"})
				return
			}
			global.RabbitMQManager.PublishToExchange(global.Config.Exchange.Test, "82f048f3-e760-44ca-b5f8-067238a52ef6", string(messageJson))
		}()
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Update test"})
}
