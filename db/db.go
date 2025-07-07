package db

import (
	"fmt"
	"os"
	"time"

	"github.com/Gabriel-Schiestl/qq-framework-log-golang/logger"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"

	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type GormDatabase struct {
	DB *gorm.DB
}

type DatabaseOptions struct {
    ServiceName string
    Config      IDbProvider
}

//dd:ignore
func NewDatabase(opts DatabaseOptions, entities ...any) (*GormDatabase, error) {
	sLog := logger.Get()

	// Definir valor padrão para ConnectionRetries se não estiver definido
	dbConnectionRetries := opts.Config.GetConnectionRetries()
	if dbConnectionRetries == 0 {
		dbConnectionRetries = 3 // Valor padrão, por exemplo, 3 tentativas
	}

	newLogger := gormLogger.New(
		sLog,
		gormLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	
	var err error

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", opts.Config.GetUser(), opts.Config.GetPass(), opts.Config.GetHost(), opts.Config.GetPort(), opts.Config.GetName())

	serviceName := opts.ServiceName
    if serviceName == "" {
        serviceName = os.Getenv("APP_NAME") + "-postgres"
    }
	
	sqltrace.Register("pgx", &stdlib.Driver{}, sqltrace.WithChildSpansOnly())
	sqlDb, err := sqltrace.Open("pgx", dbURI, sqltrace.WithServiceName(serviceName))
	if err != nil {
		sLog.Errorf("failed to open sqltrace, error: %s", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDb}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		sLog.Errorf("failed to open gormtrace, error: %s", err)
		return nil, err
	}

	time.Sleep(5 * time.Second)

	sqlDb, err = db.DB()
	if err != nil {
		sLog.Errorf("failed to get sql.DB from gorm.DB, error: %s", err)
		return nil, err
	}

	sqlDb.SetMaxIdleConns(opts.Config.GetMaxIdleConnections())
	sqlDb.SetMaxOpenConns(opts.Config.GetMaxOpenConnections())
	sqlDb.SetConnMaxLifetime(time.Duration(opts.Config.GetConnectionMaxLifetime()) * time.Minute)

	if os.Getenv("DB_AUTO_MIGRATE") == "true" {
        if len(entities) > 0 {
            if err := db.AutoMigrate(entities...); err != nil {
                sLog.Errorf("failed to auto migrate entities, error: %s", err)
                return nil, err
            }
            sLog.Infof("Auto migrated %d entities successfully", len(entities))
        } else {
            sLog.Info("No entities provided for auto migration")
        }
    }

	return &GormDatabase{
		DB: db,
	}, nil
}
