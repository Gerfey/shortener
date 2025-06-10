package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Gerfey/shortener/internal/app/service"
	"github.com/Gerfey/shortener/internal/app/settings"
	"github.com/Gerfey/shortener/internal/mock"
	"github.com/Gerfey/shortener/internal/models"
	pb "github.com/Gerfey/shortener/proto/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestURLShortenerServer_ShortenURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	mockRepo.EXPECT().
		FindShortURL(gomock.Any(), "http://example.com").
		Return("", errors.New("not found"))

	mockRepo.EXPECT().
		Save(gomock.Any(), gomock.Any(), "http://example.com", "user123").
		Return("abc123", nil)

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.ShortenURLRequest{
		Url:    "http://example.com",
		UserId: "user123",
	}

	resp, err := server.ShortenURL(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8080/abc123", resp.ShortUrl)
	assert.False(t, resp.AlreadyExists)
}

func TestURLShortenerServer_ShortenURL_EmptyURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})
	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.ShortenURLRequest{
		Url:    "",
		UserId: "user123",
	}

	_, err := server.ShortenURL(context.Background(), req)
	require.Error(t, err)
	
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestURLShortenerServer_ShortenJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	mockRepo.EXPECT().
		FindShortURL(gomock.Any(), "http://example.com").
		Return("", errors.New("not found"))

	mockRepo.EXPECT().
		Save(gomock.Any(), gomock.Any(), "http://example.com", "user123").
		Return("abc123", nil)

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.ShortenJSONRequest{
		Url:    "http://example.com",
		UserId: "user123",
	}

	resp, err := server.ShortenJSON(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8080/abc123", resp.Result)
	assert.False(t, resp.AlreadyExists)
}

func TestURLShortenerServer_GetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})
	
	urls := []models.URLPair{
		{ShortURL: "abc123", OriginalURL: "http://example.com"},
		{ShortURL: "def456", OriginalURL: "http://example.org"},
	}

	mockRepo.EXPECT().
		GetUserURLs(gomock.Any(), "user123").
		Return(urls, nil)

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.GetUserURLsRequest{
		UserId: "user123",
	}

	resp, err := server.GetUserURLs(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, resp.Items, 2)
	assert.Equal(t, "http://localhost:8080/abc123", resp.Items[0].ShortUrl)
	assert.Equal(t, "http://example.com", resp.Items[0].OriginalUrl)
	assert.Equal(t, "http://localhost:8080/def456", resp.Items[1].ShortUrl)
	assert.Equal(t, "http://example.org", resp.Items[1].OriginalUrl)
}

func TestURLShortenerServer_GetUserURLs_EmptyUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})
	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.GetUserURLsRequest{
		UserId: "",
	}

	_, err := server.GetUserURLs(context.Background(), req)
	require.Error(t, err)
	
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestURLShortenerServer_ShortenBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	// Мок для SaveBatch
	mockRepo.EXPECT().
		SaveBatch(gomock.Any(), gomock.Any(), "user123").
		Return(nil)

	// Моки для GetShortURL, которые вызываются после SaveBatch
	mockRepo.EXPECT().
		FindShortURL(gomock.Any(), "http://example.com").
		Return("abc123", nil)

	mockRepo.EXPECT().
		FindShortURL(gomock.Any(), "http://example.org").
		Return("def456", nil)

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.ShortenBatchRequest{
		UserId: "user123",
		Items: []*pb.BatchItem{
			{CorrelationId: "1", OriginalUrl: "http://example.com"},
			{CorrelationId: "2", OriginalUrl: "http://example.org"},
		},
	}

	resp, err := server.ShortenBatch(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, resp.Items, 2)
}

func TestURLShortenerServer_ShortenBatch_EmptyBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.ShortenBatchRequest{
		UserId: "user123",
		Items:  []*pb.BatchItem{},
	}

	_, err := server.ShortenBatch(context.Background(), req)
	require.Error(t, err)
	
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestURLShortenerServer_DeleteUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	mockRepo.EXPECT().
		DeleteUserURLsBatch(gomock.Any(), []string{"abc123", "def456"}, "user123").
		Return(nil)

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.DeleteUserURLsRequest{
		UserId: "user123",
		Urls:   []string{"abc123", "def456"},
	}

	resp, err := server.DeleteUserURLs(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestURLShortenerServer_DeleteUserURLs_EmptyUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.DeleteUserURLsRequest{
		UserId: "",
		Urls:   []string{"abc123", "def456"},
	}

	_, err := server.DeleteUserURLs(context.Background(), req)
	require.Error(t, err)
	
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestURLShortenerServer_DeleteUserURLs_EmptyURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.DeleteUserURLsRequest{
		UserId: "user123",
		Urls:   []string{},
	}

	_, err := server.DeleteUserURLs(context.Background(), req)
	require.Error(t, err)
	
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestURLShortenerServer_GetStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
		TrustedSubnet:          "192.168.0.0/24",
	})

	mockRepo.EXPECT().
		All(gomock.Any()).
		Return(map[string]string{
			"abc123": "http://example.com",
			"def456": "http://example.org",
		})

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.GetStatsRequest{
		ClientIp: "192.168.0.1",
	}

	resp, err := server.GetStats(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, int32(2), resp.Urls)
}

func TestURLShortenerServer_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	mockRepo.EXPECT().
		Ping(gomock.Any()).
		Return(nil)

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.PingRequest{}

	resp, err := server.Ping(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, resp.Available)
}

func TestURLShortenerServer_Ping_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	mockRepo.EXPECT().
		Ping(gomock.Any()).
		Return(errors.New("database unavailable"))

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.PingRequest{}

	resp, err := server.Ping(context.Background(), req)
	require.NoError(t, err)
	assert.False(t, resp.Available)
}

func TestURLShortenerServer_GetOriginalURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	mockRepo.EXPECT().
		Find(gomock.Any(), "abc123").
		Return("http://example.com", true, false)

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.GetOriginalURLRequest{
		ShortId: "abc123",
	}

	resp, err := server.GetOriginalURL(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "http://example.com", resp.OriginalUrl)
	assert.False(t, resp.Deleted)
}

func TestURLShortenerServer_GetOriginalURL_EmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})
	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.GetOriginalURLRequest{
		ShortId: "",
	}

	_, err := server.GetOriginalURL(context.Background(), req)
	require.Error(t, err)
	
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestURLShortenerServer_GetOriginalURL_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	s := settings.NewSettings(settings.ServerSettings{
		ServerShortenerAddress: "http://localhost:8080",
		ShutdownTimeout:        5 * time.Second,
	})

	mockRepo.EXPECT().
		Find(gomock.Any(), "nonexistent").
		Return("", false, false)

	shortenerService := service.NewShortenerService(mockRepo)
	urlService := service.NewURLService(s)
	
	server := NewURLShortenerServer(shortenerService, urlService, s, mockRepo)

	req := &pb.GetOriginalURLRequest{
		ShortId: "nonexistent",
	}

	_, err := server.GetOriginalURL(context.Background(), req)
	require.Error(t, err)
	
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}
