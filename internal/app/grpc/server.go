package grpc

import (
	"context"
	"net"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/app/usecase"
	"github.com/Gerfey/shortener/internal/models"
	pb "github.com/Gerfey/shortener/proto/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type URLShortenerServer struct {
	pb.UnimplementedURLShortenerServer
	shortenUseCase   *usecase.ShortenUseCase
	userURLsUseCase  *usecase.UserURLsUseCase
	statsUseCase     *usecase.StatsUseCase
	redirectUseCase  *usecase.RedirectUseCase
	pingUseCase      *usecase.PingUseCase
	url              *service.URLService
}

func NewURLShortenerServer(shortener *service.ShortenerService, url *service.URLService, s *settings.Settings, r models.Repository) *URLShortenerServer {
	return &URLShortenerServer{
		shortenUseCase:   usecase.NewShortenUseCase(shortener, s),
		userURLsUseCase:  usecase.NewUserURLsUseCase(r, s),
		statsUseCase:     usecase.NewStatsUseCase(r, s),
		redirectUseCase:  usecase.NewRedirectUseCase(r),
		pingUseCase:      usecase.NewPingUseCase(r),
		url:              url,
	}
}

func (s *URLShortenerServer) ShortenURL(ctx context.Context, req *pb.ShortenURLRequest) (*pb.ShortenURLResponse, error) {
	if req.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "URL не может быть пустым")
	}

	result, err := s.shortenUseCase.ShortenURL(ctx, req.Url, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ShortenURLResponse{
		ShortUrl:     result.FullShortURL,
		AlreadyExists: result.AlreadyExists,
	}, nil
}

func (s *URLShortenerServer) ShortenJSON(ctx context.Context, req *pb.ShortenJSONRequest) (*pb.ShortenJSONResponse, error) {
	if req.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "URL не может быть пустым")
	}

	result, err := s.shortenUseCase.ShortenURL(ctx, req.Url, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ShortenJSONResponse{
		Result:       result.FullShortURL,
		AlreadyExists: result.AlreadyExists,
	}, nil
}

func (s *URLShortenerServer) ShortenBatch(ctx context.Context, req *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Пакет не может быть пустым")
	}

	var batchItems []models.BatchRequestItem
	for _, item := range req.Items {
		batchItems = append(batchItems, models.BatchRequestItem{
			CorrelationID: item.CorrelationId,
			OriginalURL:   item.OriginalUrl,
		})
	}

	results, err := s.shortenUseCase.ShortenBatch(ctx, batchItems, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var responseItems []*pb.BatchResponseItem
	for _, result := range results {
		responseItems = append(responseItems, &pb.BatchResponseItem{
			CorrelationId: result.CorrelationID,
			ShortUrl:      result.FullShortURL,
		})
	}

	return &pb.ShortenBatchResponse{
		Items: responseItems,
	}, nil
}

func (s *URLShortenerServer) GetUserURLs(ctx context.Context, req *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "ID пользователя не может быть пустым")
	}

	urls, err := s.userURLsUseCase.GetUserURLs(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(urls) == 0 {
		return &pb.GetUserURLsResponse{Items: []*pb.UserURLInfo{}}, nil
	}

	var items []*pb.UserURLInfo
	for _, url := range urls {
		items = append(items, &pb.UserURLInfo{
			ShortUrl:    url.ShortURL,
			OriginalUrl: url.OriginalURL,
		})
	}

	return &pb.GetUserURLsResponse{
		Items: items,
	}, nil
}

func (s *URLShortenerServer) DeleteUserURLs(ctx context.Context, req *pb.DeleteUserURLsRequest) (*pb.DeleteUserURLsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "ID пользователя не может быть пустым")
	}

	if len(req.Urls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Список URL не может быть пустым")
	}

	err := s.userURLsUseCase.DeleteUserURLs(ctx, req.UserId, req.Urls)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteUserURLsResponse{
		Success: true,
	}, nil
}

func (s *URLShortenerServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	result, err := s.statsUseCase.GetStats(ctx, req.ClientIp)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &pb.GetStatsResponse{
		Urls:  int32(result.URLs),
		Users: int32(result.Users),
	}, nil
}

func (s *URLShortenerServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	err := s.pingUseCase.Ping(ctx)
	if err != nil {
		return &pb.PingResponse{Available: false}, nil
	}

	return &pb.PingResponse{Available: true}, nil
}

func (s *URLShortenerServer) GetOriginalURL(ctx context.Context, req *pb.GetOriginalURLRequest) (*pb.GetOriginalURLResponse, error) {
	if req.ShortId == "" {
		return nil, status.Error(codes.InvalidArgument, "ID не может быть пустым")
	}

	result, err := s.redirectUseCase.GetOriginalURL(ctx, req.ShortId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "URL не найден")
	}

	return &pb.GetOriginalURLResponse{
		OriginalUrl: result.OriginalURL,
		Deleted:     result.IsDeleted,
	}, nil
}

func RunGRPCServer(settings *settings.Settings, shortener *service.ShortenerService, url *service.URLService, repository models.Repository) error {
	server := grpc.NewServer()
	urlShortenerServer := NewURLShortenerServer(shortener, url, settings, repository)
	pb.RegisterURLShortenerServer(server, urlShortenerServer)

	listener, err := net.Listen("tcp", settings.GRPCAddress())
	if err != nil {
		return err
	}

	logrus.Infof("Запуск gRPC сервера на %s", settings.GRPCAddress())
	return server.Serve(listener)
}
