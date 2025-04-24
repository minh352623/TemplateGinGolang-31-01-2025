package inittiallize

import (
	"database/sql"
	"ecom/global"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func checkErrorPanicC(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

func initPostgresC() {
	p := global.Config.Postgres
	// Check if the port is a string, and convert it to int
	port, err := strconv.Atoi(p.Port)
	if err != nil {
		checkErrorPanic(err, "Invalid port format")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", p.Host, p.User, p.Password, p.DBName, port, p.TimeZone)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		checkErrorPanic(err, "initPostgres error")
	}
	global.Pdbc = db
	SetPoolC()
	// migrateTable()
	// genTableDAO()
}

func SetPoolC() {
	p := global.Config.Postgres
	global.Pdbc.SetConnMaxIdleTime(time.Duration(p.MaxIdleConns))
	global.Pdbc.SetMaxOpenConns(p.MaxOpenConns)
	global.Pdbc.SetConnMaxLifetime(time.Duration(p.ConnMaxLifetime))
}

// func genTableDAOC() {
// 	g := gen.NewGenerator(gen.Config{
// 		OutPath: "./internal/model",
// 		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
// 	})

// 	g.UseDB(global.Pdb) // reuse your gorm db
// 	// g.GenerateModel("users")
// 	// g.GenerateModel("wallet")
// 	g.GenerateModel("transactions")

// 	//   // Generate the code
// 	g.Execute()
// }
