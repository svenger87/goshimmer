package workerpool

import (
	"runtime"
)

var DEFAULT_OPTIONS = &Options{
	WorkerCount: runtime.NumCPU(),
	QueueSize:   runtime.NumCPU(),
}

func WorkerCount(workerCount int) Option {
	return func(args *Options) {
		args.WorkerCount = workerCount
	}
}

func QueueSize(queueSize int) Option {
	return func(args *Options) {
		args.QueueSize = queueSize
	}
}

type Options struct {
	WorkerCount int
	QueueSize   int
}

func (options Options) Override(optionalOptions ...Option) *Options {
	result := &options
	for _, option := range optionalOptions {
		option(result)
	}

	return result
}

type Option func(*Options)
