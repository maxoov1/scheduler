package scheduler

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Scheduler struct {
	cancels   []context.CancelFunc
	waitGroup *sync.WaitGroup
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cancels:   []context.CancelFunc{},
		waitGroup: &sync.WaitGroup{},
	}
}

type ExecuteJobOptions struct {
	Job       any
	Arguments []any
	Timeout   time.Duration
}

func (s *Scheduler) ExecuteJob(ctx context.Context, opts ExecuteJobOptions) error {
	job, err := s.validateJob(opts.Job, opts.Arguments)
	if err != nil {
		return fmt.Errorf("failed to validate job: %w", err)
	}

	arguments := s.convertArgumentsToValues(opts.Arguments)

	ctx, cancel := context.WithCancel(ctx)
	s.cancels = append(s.cancels, cancel)

	s.waitGroup.Add(1)
	go s.startExecution(ctx, opts.Timeout, job, arguments)

	return nil
}

func (s *Scheduler) Shutdown() {
	for _, cancel := range s.cancels {
		cancel()
	}

	s.waitGroup.Wait()
}

func (s *Scheduler) validateJob(job any, arguments []any) (reflect.Value, error) {
	value := reflect.ValueOf(job)

	if value.Kind() != reflect.Func {
		return reflect.Value{}, fmt.Errorf("provided job is not a function")
	}

	numberOfArguments := value.Type().NumIn()

	if numberOfArguments != len(arguments) {
		return reflect.Value{}, fmt.Errorf("too few arguments provided")
	}

	return value, nil
}

func (s *Scheduler) convertArgumentsToValues(arguments []any) []reflect.Value {
	values := make([]reflect.Value, len(arguments))

	for _, argument := range arguments {
		values = append(values, reflect.ValueOf(argument))
	}

	return values
}

func (s *Scheduler) startExecution(ctx context.Context, t time.Duration, j reflect.Value, args []reflect.Value) {
	defer s.waitGroup.Done()

	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			j.Call(args)

		case <-ctx.Done():
			return
		}
	}
}
