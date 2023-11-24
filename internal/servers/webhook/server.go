package webhook

import (
	"context"
	"sync"
	"time"
	"webhook/internal/services"
	"webhook/internal/services/worker"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type WorkerGroupInterface interface {
	Run()
	WaitAllDone()
	Stop()
}

type RateLimiterInterface interface {
	Take() time.Time
}

type RequestSenderInterface interface {
	Send(url string, reqNum int) (int, error)
}

type TaskGeneratorInterface interface {
	GenerateTasks(requestAmount int) chan int
}

type Server struct {
	Settings      Settings
	RequstSender  RequestSenderInterface
	RateLimiter   RateLimiterInterface
	TaskGenerator TaskGeneratorInterface
	WorkerGroup   WorkerGroupInterface
	Log           *zap.Logger
}

func NewServer(lc fx.Lifecycle, logger *zap.Logger, requestSender *services.RequestSender, tasksGenerator *services.TaskGenerator, rateLimiter *services.RateLimiter) *Server {
	server := Server{
		Log:           logger,
		RequstSender:  requestSender,
		TaskGenerator: tasksGenerator,
		RateLimiter:   rateLimiter,
	}
	return &server
}

func (s *Server) Start() error {

	s.Log.Debug("starting webhook server")

	// Loading tasks
	tasksChan := s.TaskGenerator.GenerateTasks(s.Settings.RequestAmount)
	defer close(tasksChan)

	s.Log.Debug("tasks generated", zap.Int("tasks_amount", s.Settings.RequestAmount))

	// Starting workers group

	s.WorkerGroup = worker.NewGroup(s.Worker, s.Settings.WorkersNumber, tasksChan, s.Log)
	s.WorkerGroup.Run()

	s.Log.Debug("workgroups started")

	s.WorkerGroup.WaitAllDone()

	s.Log.Debug("all done. Press CTRL+C to exit")

	return nil
}

func (s *Server) Stop() error {

	s.Log.Debug("stopping webhook server")

	s.WorkerGroup.Stop()

	s.Log.Debug("stopped webhook server")

	return nil
}

func (s *Server) Worker(ctx context.Context, taskChan chan int, wg *sync.WaitGroup) error {

	s.Log.Debug("starting worker")
	defer wg.Done()

	for {

		if len(taskChan) <= 0 {
			s.Log.Debug("no more tasks. exiting worker")
			return nil
		}

		s.RateLimiter.Take()

		select {
		case requestNum := <-taskChan:

			s.Log.Debug("processing request", zap.Int("request_num", requestNum))

			statusCode, err := s.RequstSender.Send(s.Settings.TargetURL, requestNum)
			if err != nil {
				s.Log.Error("sending request", zap.Error(err))
			}

			s.Log.Debug("processed request", zap.Int("request_num", requestNum), zap.Int("status_code", statusCode))

		case <-ctx.Done():

			s.Log.Debug("stopping worker")
			return nil

		default:
		}

	}
}
