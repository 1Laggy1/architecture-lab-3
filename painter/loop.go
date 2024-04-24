package painter

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправлення останнього разу у Receiver

	mq messageQueue

	stop    chan struct{}
	stopReq bool
	stopped chan struct{}
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.mq = messageQueue{}
	// TODO: стартувати цикл подій.
	l.stopped = make(chan struct{})
	go func() {
		for !l.stopReq || !l.mq.empty() {
			op := l.mq.pull()
			update := op.Do(l.next)
			if update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
		close(l.stopped)
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
l.Post(OperationFunc(func(screen.Texture) {
	l.stopReq = true
}))
<-l.stopped
}

// TODO: Реалізувати чергу подій.
type messageQueue struct{
	queue []Operation
}

func (mq *messageQueue) push(op Operation) {
mq.queue = append(mq.queue, op) 
}

func (mq *messageQueue) pull() Operation {
	if len(mq.queue) == 0 {
		return nil
	}
	op := mq.queue[0]
	mq.queue = mq.queue[1:]
	return op
}

func (mq *messageQueue) empty() bool {
	return len(mq.queue) == 0
}
