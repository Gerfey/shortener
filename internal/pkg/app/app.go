package app

import (
	"github.com/Gerfey/shortener/internal/app/database"
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/google/uuid"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	middleware2 "github.com/Gerfey/shortener/internal/app/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/Gerfey/shortener/internal/app/handler"
	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ShortenerApp struct {
	settings *settings.Settings
	handler  *handler.URLHandler
	router   *chi.Mux
}

func NewShortenerApp(s *settings.Settings, db *database.Database) (*ShortenerApp, error) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	application := &ShortenerApp{}

	application.settings = s

	memoryRepository := repository.NewURLMemoryRepository()
	fileStorageService := service.NewFileStorage(s.Server.DefaultFilePath)

	err := loadFileData(fileStorageService, memoryRepository)
	if err != nil {
		return nil, err
	}

	go func() {
		<-c
		err := saveFileData(fileStorageService, memoryRepository)
		if err != nil {
			log.Errorf("failed to save file data: %v", err)
		}
		os.Exit(1)
	}()

	defer func() {
		err := saveFileData(fileStorageService, memoryRepository)
		if err != nil {
			log.Errorf("failed to save file data: %v", err)
		}
	}()

	shortenerService := service.NewShortenerService(memoryRepository)
	URLService := service.NewURLService(s)

	application.handler = handler.NewURLHandler(shortenerService, URLService, db)

	r := chi.NewRouter()

	r.Use(middleware2.LoggingMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware2.GzipMiddleware)

	application.router = r

	return application, nil
}

func saveFileData(fs *service.FileStorage, mr *repository.URLMemoryRepository) error {
	allMemoryURL := mr.All()

	for key, url := range allMemoryURL {
		urlInfo := service.URLInfo{
			UUID:        uuid.New().String(),
			ShortURL:    key,
			OriginalURL: url,
		}

		err := fs.Save(urlInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadFileData(fs *service.FileStorage, mr *repository.URLMemoryRepository) error {
	urlInfos, _ := fs.Load()

	for _, urlInfo := range urlInfos {
		err := mr.Save(urlInfo.ShortURL, urlInfo.OriginalURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ShortenerApp) Run() {
	log.Printf("Starting server: %v", a.settings.ServerAddress())

	a.router.Route("/", func(r chi.Router) {
		r.Get("/ping", a.handler.Ping)
		r.Post("/", a.handler.ShortenURLHandler)
		r.Get("/{id}", a.handler.RedirectURLHandler)
		r.Post("/api/shorten", a.handler.ShortenJSONHandler)
	})

	err := http.ListenAndServe(a.settings.ServerAddress(), a.router)
	if err != nil {
		log.Fatal(err)
	}
}
