package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vincentwijaya/go-ocr/internal/app/domain"
	"github.com/vincentwijaya/go-ocr/internal/app/handler"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/face"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/vehicle"
	"github.com/vincentwijaya/go-ocr/internal/app/usecase/validate"
	"github.com/vincentwijaya/go-pkg/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/gcfg.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Server   ServerConfig
	Log      LogConfig
	Database DBConfig
	IMGBB    IMGBBConfig
}

type ServerConfig struct {
	Port        string
	Environment string
}

type LogConfig struct {
	LogPath string
	Level   string
	Stdout  bool
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

type IMGBBConfig struct {
	APIKey string
}

const fileLocation = "/etc/ocr/"
const devFileLocation = "files/etc/ocr/"
const fileName = "ocr.%s.yaml"
const devFileName = "ocr.yaml"

const infoFile = "ocr.info.log"
const errorFile = "ocr.error.log"

const banner = `
________  ________  ________     
|\   __  \|\   ____\|\   __  \    
\ \  \|\  \ \  \___|\ \  \|\  \   
 \ \  \\\  \ \  \    \ \   _  _\  
  \ \  \\\  \ \  \____\ \  \\  \| 
   \ \_______\ \_______\ \__\\ _\ 
    \|_______|\|_______|\|__|\|__|
`

func main() {
	//Read config
	var config Config
	location, fileName := getConfigLocation()
	err := gcfg.ReadFileInto(&config, location+fileName)
	if err != nil {
		log.Error("Failed to start service:", err)
		return
	}

	logConfig := log.LogConfig{
		StdoutFile: config.Log.LogPath + infoFile,
		StderrFile: config.Log.LogPath + errorFile,
		Level:      config.Log.Level,
		Stdout:     config.Log.Stdout,
	}
	log.InitLogger(config.Server.Environment, logConfig, []string{})

	log.Info(banner)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", config.Database.Host, config.Database.User, config.Database.Password, config.Database.DB, config.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to DB")
	}
	db.AutoMigrate(&domain.Member{}, &domain.Face{}, &domain.Vehicle{})

	// Repository
	vehicleRepo := vehicle.NewVehicleRepo(db)
	faceRepo := face.NewFaceRepo(db)

	// Usecase
	validateUC := validate.New(*vehicleRepo, *faceRepo)

	// Handler
	httpHandler := handler.New(validateUC)

	checker := systemCheck{
		pinger: map[string]Tester{},
	}

	httpRouter := chi.NewRouter()

	//CORS
	// httpRouter.Use(cors.Handler(cors.Options{
	// 	// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{""},
	// 	AllowCredentials: false,
	// 	MaxAge:           300,
	// 	Debug:            true,
	// }))

	httpRouter.Use(middleware.RequestID)
	httpRouter.Use(middleware.RealIP)
	httpRouter.Use(middleware.Timeout(60 * time.Second))

	httpRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("nothing here"))
	})
	httpRouter.Get("/ping", checker.ping)

	httpRouter.Route("/v1", func(r chi.Router) {
		r.Post("/validate-vehicle", httpHandler.ValidateVehicleAndOwner)
	})

	log.Infof("Service Started on:%v", config.Server.Port)
	err = http.ListenAndServe(config.Server.Port, httpRouter)
	if err != nil {
		log.Info("Failed serving Chi Dispatcher:", err)
		return
	}
	log.Info("Serving Chi Dispatcher on port:", config.Server.Port)
}

//-----------[ Pinger ]-----------------

type Tester interface {
	Ping() error
}

type systemCheck struct {
	pinger map[string]Tester
}

func (sys *systemCheck) ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func getConfigLocation() (string, string) {
	env := os.Getenv("ENV")
	location := devFileLocation
	name := devFileName
	if env == "staging" || env == "production" || env == "development" {
		location = fileLocation
		name = fmt.Sprintf(fileName, env)
	}
	return location, name
}
