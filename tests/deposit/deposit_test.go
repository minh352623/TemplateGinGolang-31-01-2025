package deposit

import (
    "context"
    "testing"

    "ecom/internal/service/deposit"
    "ecom/internal/service/deposit/types"

    "github.com/stretchr/testify/assert"
)

type DepositTestCase struct {
    name    string
    setup   func()
    input   *types.Deposit
    want    *types.Deposit
    wantErr error
}

type mockDepositHandler struct {
    resp *types.Deposit
    err  error
}

func (m *mockDepositHandler) ProcessDeposit(ctx context.Context, d *types.Deposit) (*types.Deposit, error) {
    return m.resp, m.err
}

type mockDepositInfo struct {
    resp *types.Deposit
    err  error
}

func (m *mockDepositInfo) GetDepositById(ctx context.Context, id string) (*types.Deposit, error) {
    return m.resp, m.err
}

func TestDepositHandler_ProcessDeposit(t *testing.T) {
    cases := []DepositTestCase{
        {
            name: "success",
            setup: func() {
                deposit.InitDepositHandler(&mockDepositHandler{resp: &types.Deposit{ID: "abc"}})
            },
            input: &types.Deposit{ID: "abc"},
            want:  &types.Deposit{ID: "abc"},
            wantErr: nil,
        },
        {
            name: "error",
            setup: func() {
                deposit.InitDepositHandler(&mockDepositHandler{resp: nil, err: assert.AnError})
            },
            input:   &types.Deposit{ID: "err"},
            want:    nil,
            wantErr: assert.AnError,
        },
    }

    for _, c := range cases {
        c.setup()
        got, err := deposit.DepositHandler().ProcessDeposit(context.Background(), c.input)
        if c.wantErr != nil {
            assert.Error(t, err)
            assert.Equal(t, c.want, got)
            continue
        }
        assert.NoError(t, err)
        assert.Equal(t, c.want, got)
    }
}

func TestDepositInfo_GetDepositById(t *testing.T) {
    cases := []DepositTestCase{
        {
            name: "success",
            setup: func() {
                deposit.InitDepositInfo(&mockDepositInfo{resp: &types.Deposit{ID: "xyz"}})
            },
            input: &types.Deposit{ID: "xyz"},
            want:  &types.Deposit{ID: "xyz"},
            wantErr: nil,
        },
        {
            name: "error",
            setup: func() {
                deposit.InitDepositInfo(&mockDepositInfo{resp: nil, err: assert.AnError})
            },
            input:   &types.Deposit{ID: "err"},
            want:    nil,
            wantErr: assert.AnError,
        },
    }

    for _, c := range cases {
        c.setup()
        got, err := deposit.DepositInfo().GetDepositById(context.Background(), c.input.ID)
        if c.wantErr != nil {
            assert.Error(t, err)
            assert.Equal(t, c.want, got)
            continue
        }
        assert.NoError(t, err)
        assert.Equal(t, c.want, got)
    }
}