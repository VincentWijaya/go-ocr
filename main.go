package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/vincentwijaya/go-ocr/pkg/mailer"

	"github.com/vincentwijaya/go-ocr/internal/app/domain"
	"github.com/vincentwijaya/go-ocr/internal/app/handler"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/face"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/vehicle"
	"github.com/vincentwijaya/go-ocr/internal/app/usecase/validate"
	"github.com/vincentwijaya/go-pkg/log"

	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port             string `env:"PORT" envDefault:"9000"`
	Environment      string `env:"ENVIRONMENT" envDefault:"dev"`
	LogPath          string `env:"LOGPATH"`
	Level            string `env:"LEVEL"`
	Stdout           bool   `env:"STDOUT" envDefault:"true"`
	DBHost           string `env:"DBHOST"`
	DBPort           string `env:"DBPORT"`
	DBUser           string `env:"DBUSER"`
	DBPassword       string `env:"DBPASSWORD"`
	DBName           string `env:"DBNAME"`
	MailJetAPIKey    string `env:"MAILJETAPIKEY"`
	MailJetSecretKey string `env:"MAILJETSECRETKEY"`
	ALPRSecretKey    string `env:"ALPRSECRETKEY"`
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
	err := godotenv.Load()
	if err != nil {
		log.Error("unable to load .env file: ", err)
		return
	}

	config := Config{}
	env.Parse(&config)

	logConfig := log.LogConfig{
		StdoutFile: config.LogPath + infoFile,
		StderrFile: config.LogPath + errorFile,
		Level:      config.Level,
		Stdout:     config.Stdout,
	}
	log.InitLogger(config.Environment, logConfig, []string{})

	log.Info(banner)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to DB")
	}
	db.AutoMigrate(&domain.Member{}, &domain.Face{}, &domain.Vehicle{})

	mailjetClient := mailer.Init(config.MailJetAPIKey, config.MailJetSecretKey)

	// Repository
	vehicleRepo := vehicle.NewVehicleRepo(db)
	faceRepo := face.NewFaceRepo(db)

	// Usecase
	validateUC := validate.New(*vehicleRepo, *faceRepo, *mailjetClient, config.ALPRSecretKey)

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

	httpRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var filepath = path.Join("files/views", "index.html")
		var tmpl, err = template.ParseFiles(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, make(map[string]interface{}))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	httpRouter.Route("/v1", func(r chi.Router) {
		r.Post("/validate-vehicle", httpHandler.ValidateVehicleAndOwner)
	})

	log.Infof("Service Started on:%v", config.Port)
	err = http.ListenAndServe(config.Port, httpRouter)
	if err != nil {
		log.Info("Failed serving Chi Dispatcher:", err)
		return
	}
	log.Info("Serving Chi Dispatcher on port:", config.Port)
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
