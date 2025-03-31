package inittiallize

import (
	"ecom/global"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func checkErrorPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

func initPostgres() {
	p := global.Config.Postgres
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
	global.Pdb = db
	SetPool()
	// migrateTable()
	// genTableDAO()
}

func SetPool() {
	p := global.Config.Postgres
	pql, err := global.Pdb.DB()
	if err != nil {
		fmt.Println("get pql error", err)
	}
	pql.SetConnMaxIdleTime(time.Duration(p.MaxIdleConns))
	pql.SetMaxOpenConns(p.MaxOpenConns)
	pql.SetConnMaxLifetime(time.Duration(p.ConnMaxLifetime))
}

func genTableDAO() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/model",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(global.Pdb) // reuse your gorm db
	// g.GenerateModel("users")
	// g.GenerateModel("wallet")
	g.GenerateModel("transactions")

	//   // Generate the code
	g.Execute()
}
