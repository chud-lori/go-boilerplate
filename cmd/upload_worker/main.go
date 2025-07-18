package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/chud-lori/go-boilerplate/config"
	"github.com/chud-lori/go-boilerplate/domain/entities"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/chud-lori/go-boilerplate/infrastructure/api_clients"
	"github.com/chud-lori/go-boilerplate/infrastructure/cache"
	"github.com/chud-lori/go-boilerplate/infrastructure/queue"
	"github.com/chud-lori/go-boilerplate/internal/utils"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type UploadHandlerDeps struct {
	RedisCache ports.Cache
	Uploader   *api_clients.MockUploader
}

type UploadJobHandlerFunc func(ctx context.Context, payload []byte) error

func NewUploadJobHandler(deps UploadHandlerDeps, logger *logrus.Entry) UploadJobHandlerFunc {
	return func(ctx context.Context, payload []byte) error {
		var job entities.UploadJobMessage
		if err := json.Unmarshal(payload, &job); err != nil {
			log.Printf("Failed to unmarshal upload job: %v", err)
			return nil
		}

		if job.RequestID != "" {
			logger = logger.WithField("request_id", job.RequestID)
		}
		statusKey := "upload_status:" + job.UploadID

		deps.RedisCache.Set(ctx, statusKey, []byte(string(entities.UploadStatusUploading)), time.Hour)
		logger.Printf("[Worker] Processing upload for post %s, file %s", job.PostID, job.FileName)

		url, err := deps.Uploader.Upload(ctx, job.FileName, job.FileData)

		if err != nil {
			deps.RedisCache.Set(ctx, statusKey, []byte(string(entities.UploadStatusFailed)), time.Hour)
			logger.Errorf("[Worker] Upload failed: %v", err)
			return nil
		}

		deps.RedisCache.Set(ctx, statusKey, []byte(string(entities.UploadStatusSuccess)), time.Hour)
		logger.Infof("[Worker] Upload successful: %s", url)
		return nil
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load environment variables")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	baseLogger := logger.NewLogger(cfg.LogLevel)
	workerLogger := baseLogger.WithFields(logrus.Fields{
		"layer":  "consumer",
		"driver": cfg.RabbitMQURL,
	})

	jobQueue, err := queue.NewRabbitMQJobQueue(cfg.RabbitMQURL, baseLogger)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	redisCache, err := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, baseLogger)
	if err != nil {
		baseLogger.Fatalf("Failed to connect to Redis: %v", err)
	}
	uploader := api_clients.NewMockUploader()

	handler := NewUploadJobHandler(UploadHandlerDeps{
		RedisCache: redisCache,
		Uploader:   uploader,
	}, workerLogger)

	ctx := context.Background()
	shutdown := utils.GracefullShutdown(ctx, 10*time.Second, map[string]utils.Operation{
		"rabbitmq": func(ctx context.Context) error {
			return jobQueue.Close()
		},
		"redis": func(ctx context.Context) error {
			return redisCache.Close()
		},
	})

	queueName := "post_upload_queue"
	baseLogger.Printf("[UploadWorker] Waiting for jobs on queue: %s...", queueName)
	err = jobQueue.ConsumeJobs(ctx, queueName, func(payload []byte) error {
		return handler(ctx, payload)
	})
	if err != nil {
		baseLogger.Fatalf("[UploadWorker] Error consuming jobs: %v", err)
	}

	<-shutdown
	baseLogger.Println("Upload worker shutdown complete.")
}
