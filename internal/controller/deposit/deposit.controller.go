package controller

import (
	"ecom/internal/service/deposit"
	"ecom/internal/service/deposit/types"
	"ecom/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Deposit = new(DepositController)

type DepositController struct {
}

// @Summary Process Deposit
// @Description Process Deposit
// @Tags Deposit
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} types.Deposit
// @Router /deposit/{id} [post]
// @Security bearerAuth
// @Security bearerToken
// @BearerFormat JWT
func (c *DepositController) ProcessDeposit(ctx *gin.Context) {
	deposit, err := deposit.DepositHandler().ProcessDeposit(ctx, &types.Deposit{
		ID: ctx.Param("id"),
	})
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessResponse(ctx, http.StatusOK, deposit)
}
