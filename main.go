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

	gdkwin = drawingarea.GetWindow()

	angle1 := 0
	angle2 := 360 * 64
	innerSize := 600
	startPos := 50
	tick := 0
	timeout := 10
	duration := 4000
	direction := innerSize * timeout / duration

	glib.TimeoutAdd(uint(timeout), func() bool {
		tick += direction
		fmt.Printf("tick = %d, direction = %d\n", tick, direction)

		gc.SetRgbFgColor(gdk.NewColor("black"))
		pixmap.GetDrawable().DrawArc(gc, true, startPos-10, startPos-10, innerSize+20, innerSize+20, angle1, angle2)
		gc.SetRgbFgColor(gdk.NewColor("white"))
		pixmap.GetDrawable().DrawArc(gc, true, startPos, startPos, innerSize, innerSize, angle1, angle2)

		x := startPos + tick/2
		y := startPos + tick/2
		size := innerSize - tick
		if size <= 0 || size >= innerSize {
			direction = -direction
		}
		color := gdk.NewColorRGB(0, 0, 0)
		gc.SetRgbFgColor(color)
		pixmap.GetDrawable().DrawArc(gc, true, x, y, size, size, angle1, angle2)
		drawingarea.QueueDraw()

		return true
	})

	gtk.Main()
}
