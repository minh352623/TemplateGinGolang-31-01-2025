package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"ecom/pkg/security"
)

type DataRequest struct {
	TransactionCode string `json:"transactionCode"`
	UserID          string `json:"userID"`
}
type WebhookData struct {
	Status       string      `json:"status"`       // Capitalized to make it exported
	DataRequest  DataRequest `json:"dataRequest"`  // Capitalized to make it exported
	DataResponse interface{} `json:"dataResponse"` // Capitalized to make it exported
}

func CallWebhook(webhookUrl string, webhookData WebhookData) error {
	jsonData, err := json.Marshal(webhookData)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil

}

// retries 3 times
// func CallWebhookWithRetry(webhookUrl string, webhookData WebhookData, recipientAESKey string, securityService *security.SecurityService) error {
// 	if webhookUrl == "" {
// 		return errors.New("webhookUrl is empty")
// 	}
// 	webhookLogsService := service.NewWebhookLogsService(repo.NewWebhookLogsRepository())
// 	var err error
// 	for i := 0; i < 5; i++ {
// 		err = CallWebhookWithEncryption(webhookUrl, webhookData, recipientAESKey, securityService)
// 		if err == nil {
// 			webhookLogs := model.WebhookLog{
// 				URL:             webhookUrl,
// 				TransactionCode: webhookData.DataRequest.TransactionCode,
// 				Status:          consts.TransactionStatusSuccess,
// 				DateCreated:     time.Now(),
// 				DateUpdated:     time.Now(),
// 			}
// 			_, err = webhookLogsService.CreateWebhookLogs(webhookLogs)
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Println("Call webhook ", webhookUrl, " with retry", i, "success")
// 			return nil
// 		} else {
// 			fmt.Println("Call webhook ", webhookUrl, " with retry", i, "failed", err)
// 		}
// 		time.Sleep(time.Duration(i*5) * time.Second)
// 	}
// 	webhookLogs := model.WebhookLog{
// 		URL:             webhookUrl,
// 		TransactionCode: webhookData.DataRequest.TransactionCode,
// 		Status:          consts.TransactionStatusFailed,
// 		DateCreated:     time.Now(),
// 		DateUpdated:     time.Now(),
// 	}
// 	_, err = webhookLogsService.CreateWebhookLogs(webhookLogs)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Call webhook ", webhookUrl, " with retry failed")
// 	return errors.New("Call webhook " + webhookUrl + " with retry failed")
// }

func CallWebhookWithEncryption(webhookUrl string, webhookData WebhookData, recipientAESKey string, securityService *security.SecurityService) error {
	// Chuyển đổi payload thành JSON
	jsonData, err := json.Marshal(webhookData)
	if err != nil {
		return err
	}

	if recipientAESKey == "" || securityService == nil {
		return errors.New("recipientAESKey or securityService is empty")
	}
	// Mã hóa và ký dữ liệu
	encryptedData, err := securityService.EncryptAndSignEd25519(jsonData, securityService.PrivKey, recipientAESKey)
	if err != nil {
		return err
	}

	// Định dạng payload gửi đi
	payload := map[string]string{
		"data": encryptedData,
	}

	// Chuyển đổi payload thành JSON
	finalJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Tạo request
	request, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(finalJson))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	// Gửi request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	// defer resp.Body.Close()
	// Kiểm tra response
	if resp.StatusCode >= 400 {
		return errors.New("Call webhook " + webhookUrl + " with retry failed")
	}

	return nil
}
