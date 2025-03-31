package inittiallize

import (
	"ecom/global"
	"fmt"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func initPostgresSetting() {
	p := global.Config.PostgresSetting
	// Check if the port is a string, and convert it to int
	port, err := strconv.Atoi(p.Port)
	if err != nil {
		checkErrorPanic(err, "Invalid port format")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", p.Host, p.User, p.Password, p.DBName, port, p.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		checkErrorPanic(err, "initPostgres error")
	}
	global.PdbSetting = db
	SetPoolDbSetting()
	genTableDAODbSetting()

}

func SetPoolDbSetting() {
	p := global.Config.PostgresSetting
	pql, err := global.PdbSetting.DB()
	if err != nil {
		fmt.Println("get pql error", err)
	}
	pql.SetConnMaxIdleTime(time.Duration(p.MaxIdleConns))
	pql.SetMaxOpenConns(p.MaxOpenConns)
	pql.SetConnMaxLifetime(time.Duration(p.ConnMaxLifetime))
}

func genTableDAODbSetting() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/model",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	global.Logger.Info("genTableDAODbSetting")
	g.UseDB(global.PdbSetting)
	// g.GenerateModel("project")
	g.GenerateModel("webhook_logs")
	// g.GenerateModel("wallet_integrations")
	// g.GenerateModel("wallet_integration_currencies")
	// g.GenerateModel("cycle")
	// g.GenerateModel("platform_interest_rates")
	// g.GenerateModel("transaction_type")

	//   // Generate the code
	g.Execute()
}
