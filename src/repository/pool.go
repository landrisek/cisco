package repository

import (
	"sync"
)

type Pool struct {
	taskChan chan func()
	wg       sync.WaitGroup
}

func NewPool(capacity int) *Pool {
	pool := &Pool{
		taskChan: make(chan func(), capacity),
	}
	for i := 0; i < capacity; i++ {
		go pool.worker()
	}
	return pool
}

func (p *Pool) Schedule(task func()) {
	p.wg.Add(1)
	p.taskChan <- task
}

func (p *Pool) worker() {
	for task := range p.taskChan {
		task()
		p.wg.Done()
	}
}

func (p *Pool) Wait() {
	close(p.taskChan)
	p.wg.Wait()
}