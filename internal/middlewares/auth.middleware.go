package middlewares

import (
	"github.com/gin-gonic/gin"
)

type TokenValidationRequest struct {
	ClientID string `json:"client_id"`
	Token    string `json:"token"`
}

type TokenValidationResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Message string `json:"message"`
		Result  struct {
			Sub               string `json:"sub"`
			Name              string `json:"name"`
			PreferredUsername string `json:"preferred_username"`
		} `json:"result"`
	} `json:"data"`
	Success bool `json:"success"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 	token := c.GetHeader("Authorization")
		// 	if token == "" {
		// 		response.ErrorResponse(c, response.Unauthorized, "Missing token")
		// 		c.Abort()
		// 		return
		// 	}

		// 	if len(token) > 7 && token[:7] == "Bearer " {
		// 		token = token[7:]
		// 	}

		// 	validationReq := TokenValidationRequest{
		// 		ClientID: global.Config.TokenValidation.ClientID,
		// 		Token:    token,
		// 	}

		// 	reqBody, err := json.Marshal(validationReq)
		// 	if err != nil {
		// 		response.ErrorResponse(c, response.InternalServerError, "Failed to encode request")
		// 		c.Abort()
		// 		return
		// 	}

		// 	apiURL := global.Config.TokenValidation.URL
		// 	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
		// 	if err != nil {
		// 		response.ErrorResponse(c, response.InternalServerError, "Failed to create request")
		// 		c.Abort()
		// 		return
		// 	}
		// 	req.Header.Set("Content-Type", "application/json")
		// 	req.Header.Set("x-api-key", global.Config.TokenValidation.XApiKey) // change to basic auth

		// 	req.Header.Set("Content-Type", "application/json")
		// 	client := &http.Client{}
		// 	resp, err := client.Do(req)
		// 	if err != nil {
		// 		response.ErrorResponse(c, response.InternalServerError, "Failed to validate token")
		// 		c.Abort()
		// 		return
		// 	}
		// 	defer resp.Body.Close()

		// 	body, err := io.ReadAll(resp.Body)
		// 	if err != nil {
		// 		response.ErrorResponse(c, response.InternalServerError, "Failed to read response")
		// 		c.Abort()
		// 		return
		// 	}

		// 	var validationResp TokenValidationResponse
		// 	err = json.Unmarshal(body, &validationResp)
		// 	if err != nil {
		// 		response.ErrorResponse(c, response.InternalServerError, "Invalid response format")
		// 		c.Abort()
		// 		return
		// 	}

		// 	if !validationResp.Success || validationResp.Code != 200 {
		// 		response.ErrorResponse(c, response.Unauthorized, validationResp.Message)
		// 		c.Abort()
		// 		return
		// 	}

		// 	var userData = validationResp.Data.Result

		// 	userRepo := repo.NewUserRepository()
		// 	user, _ := userRepo.GetUserByID(userData.Sub)

		// 	if user.ID == "" {
		// 		user, err = userRepo.CreateUser(model.User{
		// 			ID:          userData.Sub,
		// 			ProviderKey: "google",
		// 			NickName:    userData.Name,
		// 		})
		// 		if err != nil {
		// 			response.ErrorResponse(c, response.InternalServerError, "Failed to create user")
		// 			c.Abort()
		// 		}
		// 	}

		// 	c.Set("userId", user.ID)
		// 	fmt.Println("🚀 ~ returnfunc ~ userId:", user.ID)

		// 	c.Set("userData", userData)

		c.Next()
	}
}
