package controller

import (
	"ecom/internal/service"
	"ecom/internal/vo"
	"ecom/pkg/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

type DepositController struct {
	depositService service.IDepositService
}

func NewDepositController(depositService service.IDepositService) *DepositController {
	return &DepositController{depositService: depositService}
}

// PingExample godoc
// @Summary Deposit
// @Schemes http
// @Description Deposit
// @Tags Deposit
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseData
// @Router /deposit [post]
// @Param data body vo.EncryptedRequest true "data"
// @SecurityScheme bearerAuth
// @Security bearerToken
// @BearerFormat JWT
// func (dc *DepositController) Deposit(c *gin.Context) {
// 	// get from body

// 	var encryptedRequest vo.EncryptedRequest
// 	err := c.ShouldBindJSON(&encryptedRequest)
// 	if err != nil {
// 		response.ErrorResponse(c, response.BadRequest, err.Error())
// 		return
// 	}

// 	// // Decrypt data
// 	var depositRequest vo.DepositRequest
// 	senderPubKey := global.Config.Security.CryptoKeys.Asymmetric.SenderPubKey
// 	recipientAESKey := global.Config.Security.CryptoKeys.Symmetric.AESKey

// 	decryptedData, err := global.SecurityService.DecryptAndVerifyEd25519(encryptedRequest.Data, recipientAESKey, senderPubKey)

// 	if err != nil {
// 		// call webhook
// 		go webhook.CallWebhookWithRetry(depositRequest.WebhookUrl, webhook.WebhookData{
// 			Status: consts.TransactionStatusFailed,
// 			DataRequest: webhook.DataRequest{
// 				TransactionCode: depositRequest.TransactionCode,
// 				UserID:          depositRequest.UserID,
// 			},
// 			DataResponse: nil,
// 		}, recipientAESKey, global.SecurityService)
// 		response.ErrorResponse(c, response.BadRequest, err.Error())
// 		return
// 	}
// 	fmt.Println("decryptedData", string(decryptedData))
// 	err = json.Unmarshal(decryptedData, &depositRequest)
// 	fmt.Println("depositRequest", depositRequest)
// 	if err != nil {
// 		// call webhook
// 		go webhook.CallWebhookWithRetry(depositRequest.WebhookUrl, webhook.WebhookData{
// 			Status: consts.TransactionStatusFailed,
// 			DataRequest: webhook.DataRequest{
// 				TransactionCode: depositRequest.TransactionCode,
// 				UserID:          depositRequest.UserID,
// 			},
// 			DataResponse: err.Error(),
// 		}, recipientAESKey, global.SecurityService)
// 		response.ErrorResponse(c, response.BadRequest, err.Error())
// 		return
// 	}

// 	result, err := dc.depositService.Deposit(depositRequest.UserID, depositRequest.Currency, depositRequest.Amount, depositRequest.RateUsd, depositRequest.ProviderKey, depositRequest.Platform, depositRequest.WebhookUrl, depositRequest.TransactionCode, depositRequest.RateCurrency)
// 	if err != nil {
// 		// call webhook
// 		go webhook.CallWebhookWithRetry(depositRequest.WebhookUrl, webhook.WebhookData{
// 			Status: consts.TransactionStatusFailed,
// 			DataRequest: webhook.DataRequest{
// 				TransactionCode: depositRequest.TransactionCode,
// 				UserID:          depositRequest.UserID,
// 			},
// 			DataResponse: err.Error(),
// 		}, recipientAESKey, global.SecurityService)
// 		response.ErrorResponse(c, response.BadRequest, err.Error())
// 		return
// 	}
// 	go webhook.CallWebhookWithRetry(depositRequest.WebhookUrl, webhook.WebhookData{
// 		Status: consts.TransactionStatusSuccess,
// 		DataRequest: webhook.DataRequest{
// 			TransactionCode: depositRequest.TransactionCode,
// 			UserID:          depositRequest.UserID,
// 		},
// 		DataResponse: result,
// 	}, recipientAESKey, global.SecurityService)
// 	response.SuccessResponse(c, response.Success, gin.H{"message": "Deposit successful", "result": result})
// }

func (dc *DepositController) Test(c *gin.Context) {
	// get from body
	// Add more logging here to debug
	fmt.Println("Test--------------------------------")
	var data vo.TestMQRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.ErrorResponse(c, response.BadRequest, err.Error())
		return
	}

	result, err := dc.depositService.Test(data.UserID, data.Email, data.MessageID, data.RoutingKey, data.HashKey)
	if err != nil {
		response.ErrorResponse(c, response.BadRequest, err.Error())
		return
	}

	response.SuccessResponse(c, response.Success, gin.H{"message": "Deposit successful", "result": result})
}
