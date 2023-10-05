package repository

import (
	"SergeyProject/pkg/logger"
	"github.com/hashicorp/go-hclog"
)

const (
	Increment = iota
	Decrement
)

type (
	CmdType int

	Command struct {
		cmd       CmdType
		name      string
		val       int
		replyChan chan int
	}
)

type CounterRepository struct {
	logger  hclog.Logger
	storage map[string]int
	cmdChan chan Command
}

func NewCounterRepository() *CounterRepository {
	cr := &CounterRepository{logger: logger.GetLogger(), storage: make(map[string]int), cmdChan: make(chan Command)}
	go cr.ProcessConcurrency()
	return cr
}

var replyChan chan int

func (cr *CounterRepository) ProcessConcurrency() {
	replyChan = make(chan int)
	for data := range cr.cmdChan {
		switch data.cmd {
		//То есть если канал принимает инкремент, то при наличии элемента в списке (storage), значение увеличивается на число
		//а в значение replyChan канал, записывается новое значение для параметра списка
		case Increment:
			if _, ok := cr.storage[data.name]; ok {
				cr.storage[data.name] += data.val
				data.replyChan <- cr.storage[data.name]
			}
		case Decrement:
			if _, ok := cr.storage[data.name]; ok {
				cr.storage[data.name] -= data.val
				data.replyChan <- cr.storage[data.name]
			}
		}
	}
}

func (cr *CounterRepository) SetValue(name string, val int) int {
	cr.storage[name] = val
	return cr.storage[name]
}

func (cr *CounterRepository) GetStorage() map[string]int {
	return cr.storage
}

func (cr *CounterRepository) IncreaseValue(name string, val int) int {
	cr.cmdChan <- Command{cmd: Increment, name: name, val: val, replyChan: replyChan}
	return <-replyChan
}

func (cr *CounterRepository) DecreaseValue(name string, val int) int {
	cr.cmdChan <- Command{cmd: Decrement, name: name, val: val, replyChan: replyChan}
	return <-replyChan
}
