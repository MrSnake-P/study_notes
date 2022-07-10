// 使用通道来监视程序的执行时间，
// 在开发需要调度后台处理任务的程序时，
// 这种模式非常有用
//
// runner 在给定的超时时间内执行一组任务
// 并且在操作系统发送中断信号时结束这些任务
//
// 程序可以在分配的时间内完成工作，正常终止
// 程序没有及时完成工作，“自杀”
// 接受到操作系统发送的中断事件，
// 程序立刻试图清理状态并终止工作。

package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

type Runner struct {
	// 从系统收到的信号
	interrupt chan os.Signal

	// 报告处理任务已经完成
	complete chan error

	// 报告处理任务已经超时
	timeout <-chan time.Time

	// 持有一组以索引顺序依次执行的函数
	tasks []func(int)
}

// 任务超时错误
var ErrTimeout = errors.New("received timeout")

// 接收到操作系统的事件时返回
var ErrInterrupt = errors.New("received interrupt")

func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

// Add添加任务
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// Start执行所有任务，并监视通道时间
func (r *Runner) Start() error {
	// 接受所有终端信号
	signal.Notify(r.interrupt, os.Interrupt)

	// 并发执行不同的任务
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

func (r *Runner) run() error {
	for id, task := range r.tasks {
		if r.gotInterrupt() {
			return ErrInterrupt
		}

		// 执行已经注册的任务
		task(id)
	}
	return nil
}

func (r *Runner) gotInterrupt() bool {
	select {
	case <-r.interrupt:
		// 停止接受后续的任何信号
		signal.Stop(r.interrupt)
		return true

	default:
		// 默认继续运行
		return false
	}
}
