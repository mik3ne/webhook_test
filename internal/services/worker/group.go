package worker

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type WorkerFunc func(cancelCtx context.Context, tasksChan chan int, wg *sync.WaitGroup) error

type Group struct {
	Log           *zap.Logger
	TasksChan     chan int
	WorkersNumber int
	WorkersCtx    context.Context
	WorkersCancel context.CancelFunc
	WorkerWG      sync.WaitGroup
	WorkerFn      WorkerFunc
}

func NewGroup(fn WorkerFunc, workersNumber int, taskChan chan int, log *zap.Logger) *Group {
	return &Group{
		Log:           log,
		WorkerFn:      fn,
		WorkersNumber: workersNumber,
		TasksChan:     taskChan,
	}
}

func (g *Group) Run() {

	g.Log.Debug("starting workers", zap.Int("workers", g.WorkersNumber))

	g.WorkersCtx, g.WorkersCancel = context.WithCancel(context.Background())

	g.WorkerWG = sync.WaitGroup{}
	g.WorkerWG.Add(g.WorkersNumber)

	for i := 0; i < g.WorkersNumber; i++ {
		go g.WorkerFn(g.WorkersCtx, g.TasksChan, &g.WorkerWG)
	}
}

func (g *Group) WaitAllDone() {
	g.WorkerWG.Wait()
}

func (g *Group) Stop() {

	g.Log.Debug("stopping workers")

	g.WorkersCancel()
	g.WorkerWG.Wait()
}
