package main

import (
	"os"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type point struct {
	x int
	y int
}

func main() {
	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Go Breathing")
	window.Connect("destroy", gtk.MainQuit)

	vbox := gtk.NewVBox(true, 0)
	vbox.SetBorderWidth(5)
	drawingarea := gtk.NewDrawingArea()

	var gdkwin *gdk.Window
	var pixmap *gdk.Pixmap
	var gc *gdk.GC

	window.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		key := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		// fmt.Println(key.Keyval)

		switch key.Keyval {
		case 113, 65307:
			gtk.MainQuit()
		}
	})

	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		allocation := drawingarea.GetAllocation()
		pixmap = gdk.NewPixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)
		gc = gdk.NewGC(pixmap.GetDrawable())
		color := gdk.NewColorRGB(0xe6e6, 0xe9e9, 0xecec)
		gc.SetRgbFgColor(color)
		gc.SetRgbBgColor(gdk.NewColor("white"))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
	})

	drawingarea.Connect("expose-event", func() {
		if pixmap == nil {
			return
		}
		gdkwin.GetDrawable().DrawDrawable(gc, pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
	})

	drawingarea.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.POINTER_MOTION_HINT_MASK | gdk.BUTTON_PRESS_MASK))
	vbox.Add(drawingarea)

	window.Add(vbox)
	window.SetSizeRequest(400, 400)
	window.ShowAll()
	window.Fullscreen()

	gdkwin = drawingarea.GetWindow()

	angle1 := 0
	angle2 := 360 * 64

	percLowerLimit := 2500
	percUpperLimit := 9500
	duration := time.Second * 4
	tickDuration := time.Second * 4 / 7000
	timeout := 10
	direction := 1

	percAtom := int32(percLowerLimit + 1)

	go func() {
		for {
			perc := int(atomic.LoadInt32(&percAtom))

			if perc >= percUpperLimit || perc <= percLowerLimit {
				direction = -direction
				time.Sleep(duration)
			} else {
				time.Sleep(tickDuration)
			}

			atomic.AddInt32(&percAtom, int32(direction))
		}
	}()

	previousPerc := 0

	glib.TimeoutAdd(uint(timeout), func() bool {
		var windowW, windowH int
		window.GetSize(&windowW, &windowH)

		maxSize := windowH / 100 * 60
		centerX := windowW / 2
		centerY := windowH / 2
		startX := centerX - maxSize/2
		startY := centerY - maxSize/2

		perc := int(atomic.LoadInt32(&percAtom))

		// fmt.Printf("perc = %d, direction = %d\n", perc, direction)

		gc.SetRgbFgColor(gdk.NewColor("black"))
		pixmap.GetDrawable().DrawArc(gc, true, startX-10, startY-10, maxSize+20, maxSize+20, angle1, angle2)
		gc.SetRgbFgColor(gdk.NewColor("white"))
		pixmap.GetDrawable().DrawArc(gc, true, startX, startY, maxSize, maxSize, angle1, angle2)

		size := maxSize / 100 * perc / 100
		x := startX + (maxSize-size)/2
		y := startY + (maxSize-size)/2

		color := gdk.NewColorRGB(0x2222, 0x6f6f, 0xa0a0)
		gc.SetRgbFgColor(color)
		pixmap.GetDrawable().DrawArc(gc, true, x, y, size, size, angle1, angle2)

		if previousPerc != perc {
			drawingarea.QueueDraw()
			previousPerc = perc
			// fmt.Println("Skipping redraw queue")
		}

		return true
	})

	gtk.Main()
}
