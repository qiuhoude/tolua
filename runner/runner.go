package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

type Runner struct {
	tasks     []func(int)      // 要执行的任务
	complete  chan error       //用于通知任务全部完成
	timeout   <-chan time.Time // 超时时间
	interrupt chan os.Signal   // 中断通知
}

// 超时错误
var ErrTimeout = errors.New("received timeout")

// 中断错误
var ErrInterrupt = errors.New("received interrupt")

func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

// 添加需要执行的任务
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// 检测是否接受到中断信号
func (r *Runner) isInterrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

// 执行
func (r *Runner) run() error {
	for id, task := range r.tasks {
		//执行任务，执行的过程中接收到中断信号时，返回中断错误
		if r.isInterrupt() {
			return ErrInterrupt
		}
		//执行任务
		task(id)
	}
	//如果任务全部执行完，还没有接收到中断信号，则返回nil
	return nil
}

func (r *Runner) Start() error {
	// 希望接收哪些系统信号
	signal.Notify(r.interrupt, os.Interrupt)

	go func() {
		r.complete <- r.run()
	}()

	select {
	case err := <-r.complete:
		return err
	case <-r.timeout:
		return ErrTimeout

	}
}
