package vo

type TokenValidationRequest struct {
	ClientID string `json:"client_id"`
	Token    string `json:"token"`
}

type TokenValidationResponse struct {
	Data struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Result  struct {
			Sub               string `json:"sub"`
			Name              string `json:"name"`
			PreferredUsername string `json:"preferred_username"`
		} `json:"result"`
	} `json:"data"`
}
