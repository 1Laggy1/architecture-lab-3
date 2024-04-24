package ui

import (
    "image"
    "image/color"
    "log"

    "golang.org/x/exp/shiny/driver"
    "golang.org/x/exp/shiny/screen"
    "golang.org/x/image/draw"
    "golang.org/x/mobile/event/key"
    "golang.org/x/mobile/event/lifecycle"
    "golang.org/x/mobile/event/mouse"
    "golang.org/x/mobile/event/paint"
    "golang.org/x/mobile/event/size"
)

type Visualizer struct {
    Title         string
    Debug         bool
    OnScreenReady func(s screen.Screen)

    w    screen.Window
    tx   chan screen.Texture
    done chan struct{}

    sz  size.Event
    pos image.Rectangle

    // Додаткові поля для відображення хрестика
    crossX, crossY int
}

func (pv *Visualizer) Main() {
    pv.tx = make(chan screen.Texture)
    pv.done = make(chan struct{})
    pv.pos.Max.X = 800
    pv.pos.Max.Y = 800
    driver.Main(pv.run)
}

func (pv *Visualizer) Update(t screen.Texture) {
    pv.tx <- t
}

func (pv *Visualizer) run(s screen.Screen) {
    w, err := s.NewWindow(&screen.NewWindowOptions{
        Title: pv.Title,
				Width:  800,
        Height: 800,
    })
    if err != nil {
        log.Fatal("Failed to initialize the app window:", err)
    }
    defer func() {
        w.Release()
        close(pv.done)
    }()

    if pv.OnScreenReady != nil {
        pv.OnScreenReady(s)
    }

    pv.w = w
		pv.crossX = 400
    pv.crossY = 400

    events := make(chan interface{})
    go func() {
        for {
            e := w.NextEvent()
            if pv.Debug {
                log.Printf("new event: %v", e)
            }
            if detectTerminate(e) {
                close(events)
                break
            }
            events <- e
        }
    }()

    var t screen.Texture

    for {
        select {
        case e, ok := <-events:
            if !ok {
                return
            }
            pv.handleEvent(e, t)

        case t = <-pv.tx:
            w.Send(paint.Event{})
        }
    }
}

func detectTerminate(e interface{}) bool {
    switch e := e.(type) {
    case lifecycle.Event:
        if e.To == lifecycle.StageDead {
            return true // Window destroy initiated.
        }
    case key.Event:
        if e.Code == key.CodeEscape {
            return true // Esc pressed.
        }
    }
    return false
}

func (pv *Visualizer) handleEvent(e interface{}, t screen.Texture) {
    switch e := e.(type) {
    case size.Event: // Оновлення даних про розмір вікна.
        pv.sz = e

    case error:
        log.Printf("ERROR: %s", e)

    case mouse.Event:
        if t == nil {
            if e.Button == mouse.ButtonRight {
                // Конвертуємо координати миші до цілих чисел перед присвоєнням
                pv.crossX, pv.crossY = int(e.X), int(e.Y)
                pv.drawDefaultUI()
            }
        }

    case paint.Event:
        if t == nil {
            pv.drawDefaultUI()
        } else {
            pv.w.Scale(pv.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
        }
        pv.w.Publish()
    }
}

func (pv *Visualizer) drawDefaultUI() {
	pv.w.Fill(pv.sz.Bounds(), color.RGBA{0, 255, 0, 255}, draw.Src) // Зелений фон.

	// Малювання жовтого хрестика
	crossX := pv.crossX
	crossY := pv.crossY
	crossLen := 200  // Нова довжина хрестика
	crossThickness := 20  // Нова товщина хрестика

	// Вертикальна лінія
	verticalLine := image.Rect(crossX-crossThickness/2, crossY-crossLen/2, crossX+crossThickness/2, crossY+crossLen/2)
	pv.w.Fill(verticalLine, color.RGBA{255, 255, 0, 255}, draw.Src)

	// Горизонтальна лінія
	horizontalLine := image.Rect(crossX-crossLen/2, crossY-crossThickness/2, crossX+crossLen/2, crossY+crossThickness/2)
	pv.w.Fill(horizontalLine, color.RGBA{255, 255, 0, 255}, draw.Src)
}
