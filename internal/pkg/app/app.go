package app

import (
	"github.com/Gerfey/shortener/internal/app/repository"
	"github.com/google/uuid"
	"net/http"

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

func NewShortenerApp(s *settings.Settings) (*ShortenerApp, error) {

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

	defer saveFileData(fileStorageService, memoryRepository)

	shortenerService := service.NewShortenerService(memoryRepository)
	URLService := service.NewURLService(s)

	application.handler = handler.NewURLHandler(shortenerService, URLService)

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

		fs.Save(urlInfo)
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
	log.Println("Starting server...")

	a.router.Route("/", func(r chi.Router) {
		r.Post("/", a.handler.ShortenURLHandler)
		r.Get("/{id}", a.handler.RedirectURLHandler)
		r.Post("/api/shorten", a.handler.ShortenJSONHandler)
	})

	err := http.ListenAndServe(a.settings.ServerAddress(), a.router)
	if err != nil {
		log.Fatal(err)
	}
}
