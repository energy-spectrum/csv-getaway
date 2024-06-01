package main

import (
	"context"
	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/repository/postgres"
	"csv-analyzer-api/internal/service/template"
	"csv-analyzer-api/internal/service/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // go postgres driver
	"github.com/jmoiron/sqlx"

	//_ "github.com/lib/pq"
	_ "github.com/santosh/gingo/docs"
	log "github.com/sirupsen/logrus"

	"csv-analyzer-api/internal/transport/http/route"
)

// @title           CSV API
// @version         1.0
// @description     getaway
// @host            localhost:8000
// @BasePath        /v1
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log.SetFormatter(new(log.JSONFormatter))

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	ctx := context.Background()

	pgClient := connectToDB(ctx, &cfg.Postgres)

	userService := user.NewService(cfg, postgres.NewUser(pgClient, cfg))
	templateService := template.NewService(cfg, postgres.NewTemplate(pgClient, cfg))

	runGinServer(cfg, userService,  templateService)
}

func connectToDB(ctx context.Context, cfg *config.Postgres) *sqlx.DB {
	log.Print(cfg.DSN)
	postgresClient, err := sqlx.ConnectContext(ctx, "pgx", cfg.DSN)
	if err != nil {
		log.Fatalf("failed to connect to Postgresql: %v", err)
	}

	log.Infof("connected to Postgresql")

	return postgresClient
}

func runGinServer(cfg *config.Configuration,
	userService user.UserService,
	templateService template.TemplateService) {
	ginEngine := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	ginEngine.Use(cors.New(config))

	serveImagesFolder(ginEngine)
	connectSwaggerToGin(ginEngine)

	router := ginEngine.Group("/v1/")
	route.Setup(cfg, userService, templateService, router)

	log.Infof("server running on address: %s", cfg.HTTPServer.Address)
	ginEngine.Run(cfg.HTTPServer.Address)
}

func serveImagesFolder(ginEngine *gin.Engine) {
	ginEngine.Static("/images", "./images")
}

func connectSwaggerToGin(ginEngine *gin.Engine) {
	// Serve the Swagger UI files
	ginEngine.Static("/swagger/", "./doc/swagger")
}
