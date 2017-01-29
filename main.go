package main

import (
	"fmt"
	"os"

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
	window.SetTitle("GTK DrawingArea")
	window.Connect("destroy", gtk.MainQuit)

	vbox := gtk.NewVBox(true, 0)
	vbox.SetBorderWidth(5)
	drawingarea := gtk.NewDrawingArea()

	var gdkwin *gdk.Window
	var pixmap *gdk.Pixmap
	var gc *gdk.GC
	// colors := []string{
	// 	"black",
	// 	"gray",
	// 	"blue",
	// 	"purple",
	// 	"red",
	// 	"orange",
	// 	"yellow",
	// 	"green",
	// 	"darkgreen",
	// }

	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		allocation := drawingarea.GetAllocation()
		pixmap = gdk.NewPixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)
		gc = gdk.NewGC(pixmap.GetDrawable())
		gc.SetRgbFgColor(gdk.NewColor("white"))
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
	// window.Fullscreen()

	gdkwin = drawingarea.GetWindow()

	angle1 := 0
	angle2 := 360 * 64

	// duration := 4000
	timeout := 50
	perc := 6
	direction := 1

	glib.TimeoutAdd(uint(timeout), func() bool {
		var windowW, windowH int
		window.GetSize(&windowW, &windowH)

		maxSize := windowH / 100 * 60
		halfSize := maxSize / 2
		centerX := windowW / 2
		centerY := windowH / 2
		startX := centerX - maxSize/2
		startY := centerY - maxSize/2

		fmt.Println(windowW, windowH, maxSize, halfSize, centerX, centerY)

		if perc >= 95 || perc <= 5 {
			direction = -direction
		}

		perc += direction

		fmt.Printf("perc = %d, direction = %d\n", perc, direction)

		gc.SetRgbFgColor(gdk.NewColor("black"))
		pixmap.GetDrawable().DrawArc(gc, true, startX-10, startY-10, maxSize+20, maxSize+20, angle1, angle2)
		gc.SetRgbFgColor(gdk.NewColor("white"))
		pixmap.GetDrawable().DrawArc(gc, true, startX, startY, maxSize, maxSize, angle1, angle2)

		size := maxSize / 100 * perc
		x := startX + (maxSize-size)/2
		y := startY + (maxSize-size)/2
		fmt.Println(size, x)

		color := gdk.NewColor("blue")
		gc.SetRgbFgColor(color)
		pixmap.GetDrawable().DrawArc(gc, true, x, y, size, size, angle1, angle2)

		drawingarea.QueueDraw()
		return true
	})

	gtk.Main()
}
