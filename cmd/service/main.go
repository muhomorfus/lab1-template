package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/muhomorfus/lab1-template/internal/generated"
	"github.com/muhomorfus/lab1-template/internal/openapi"
	"github.com/muhomorfus/lab1-template/internal/person"
	"github.com/muhomorfus/lab1-template/internal/repository"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	db, err := sqlx.Connect("postgres", cfg.dsn())
	if err != nil {
		return fmt.Errorf("connect to db: %w", err)
	}

	repo := repository.New(db)
	mgr := person.New(repo)
	server := openapi.New(mgr)
	router := fiber.New()
	generated.RegisterHandlers(router, server)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-exit

		_ = db.Close()
		_ = router.Shutdown()
	}()

	if err := router.Listen(cfg.listerAddress()); err != nil {
		return fmt.Errorf("listen http server: %w", err)
	}

	return nil
}

type config struct {
	PostgresHost     string `envconfig:"PGHOST" required:"true"`
	PostgresPort     int    `envconfig:"PGPORT" required:"true"`
	PostgresUser     string `envconfig:"PGUSER" required:"true"`
	PostgresPassword string `envconfig:"PGPASSWORD" required:"true"`
	PostgresDB       string `envconfig:"PGDB" required:"true"`
	PostgresSSL      bool   `envconfig:"PGSSL" default:"false"`
	Port             string `envconfig:"PORT" required:"true"`
}

func (c config) dsn() string {
	sslMode := ""
	if !c.PostgresSSL {
		sslMode = "sslmode=disable"
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s", c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, sslMode)
}

func (c config) listerAddress() string {
	return fmt.Sprintf("0.0.0.0:%s", c.Port)
}
