package service

import (
	"ecom/internal/database"

	"github.com/google/uuid"
)

type (
	ITestCreate interface {
		CreateTest(req *database.CreateTestParams) (database.Test, error)
	}

	ITestAdmin interface {
		GetTestById(id uuid.UUID) (database.Test, error)
		RemoveTest(id uuid.UUID) error
	}
)

var (
	localTestCreate ITestCreate
	localTestAdmin  ITestAdmin
)

func TestCreate() ITestCreate {
	if localTestCreate == nil {
		panic("TestCreate not initialized")
	}
	return localTestCreate
}

func InitTestCreate(testCreate ITestCreate) {
	localTestCreate = testCreate
}

func TestAdmin() ITestAdmin {
	if localTestAdmin == nil {
		panic("TestAdmin not initialized")
	}
	return localTestAdmin
}

func InitTestAdmin(testAdmin ITestAdmin) {
	localTestAdmin = testAdmin
}
