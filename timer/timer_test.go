package timer

import (
	"github.com/sirupsen/logrus"
	"log"
	"testing"
)

type PrintHello struct {
	status Status
}

func (p *PrintHello) Description() string {
	return "print hello"
}

func (p *PrintHello) Pause() {
	p.status = 1
}

func (p *PrintHello) Start() {
	p.status = 0
}

func (p *PrintHello) Status() Status {
	return p.status
}

func (p *PrintHello) Runnable() bool {
	return true
}

func (p *PrintHello) Name() JobName {
	return "print_hello"
}

func (p PrintHello) Spec() string {
	return "@every 5s"
}

func (p *PrintHello) Run() {
	log.Println("hello")
}

func TestNewTimer(t *testing.T) {
	var timer = NewTimer(logrus.New().WithField("filter", "pkg.timer.test"))
	_ = timer.AddJob(&PrintHello{})
	err := timer.AddJob(&PrintHello{})

	if err == ErrJobExist {
		t.Log("check ok!")
	}

	timer.Run()

	t.Log(timer.Tasks().GetJob("print_hello").EntryID)
	t.Log(timer.Tasks().GetJob("print_hello").Name)

	timer.Tasks().Pause("print_hello")

	timer.Tasks().Start("print_hello")
}
