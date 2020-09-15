package ui


import (
	"fmt"
//	"context"
	"log"
//	"os"
//	"io"
	"strings"
	"time"
	"github.com/jroimartin/gocui"
)


var (
	g	*gocui.Gui
)


var version = "master"
var LOG_MOD string = "pod"

// Configure globale keys
var keys []Key = []Key{
	Key{"", gocui.KeyCtrlC, actionGlobalQuit},
//	Key{"", gocui.KeyCtrlD, actionGlobalToggleViewDebug},
//	Key{"container", gocui.KeyCtrlN, actionGlobalToggleViewNamespaces},
	Key{"con", gocui.KeyArrowUp, actionViewConUp},
	Key{"con", gocui.KeyArrowDown, actionViewConDown},
	//Key{"con", 'd', actionViewPodsDelete},
	Key{"con", gocui.KeyEnter, actionViewConSelect},
	//Key{"logs", l, actionStreamLogs},
//	Key{"logs", 'l', actionViewPodsLogsHide},
//	Key{"logs", gocui.KeyArrowUp, actionViewPodsLogsUp},
//	Key{"logs", gocui.KeyArrowDown, actionViewPodsLogsDown},
//	Key{"namespaces", gocui.KeyArrowUp, actionViewNamespacesUp},
//	Key{"namespaces", gocui.KeyArrowDown, actionViewNamespacesDown},
//	Key{"namespaces", gocui.KeyEnter, actionViewNamespacesSelect},
}




// Entry Point of the x-tracer
func InitGui() {
//	c := getConfig()

	// Ask version
/*	if c.askVersion {
		fmt.Println(versionFull())
		os.Exit(0)
	}

	// Ask Help
	if c.askHelp {
		fmt.Println(versionFull())
		fmt.Println(HELP)
		os.Exit(0)
	}*/

	// Only used to check errors
	//getClientSet()

	G, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer G.Close()

	G.Highlight = true
	G.SelFgColor = gocui.ColorGreen

	G.SetManagerFunc(uiLayout)

	if err := uiKey(G); err != nil {
		log.Panicln(err)
	}

	if err := G.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}



// Define the UI layout
func uiLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	//viewDebug(g, maxX, maxY)
	viewLogs(g, maxX, maxY)
	//viewNamespaces(g, maxX, maxY)
	viewOverlay(g, maxX, maxY)
	viewTitle(g, maxX, maxY)
	viewCon(g, maxX, maxY)
	viewStatusBar(g, maxX, maxY)

	return nil
}



// Move view cursor to the bottom
func moveViewCursorDown(g *gocui.Gui, v *gocui.View, allowEmpty bool) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	nextLine, err := getNextViewLine(g, v)
	if err != nil {
		return err
	}
	if !allowEmpty && nextLine == "" {
		return nil
	}
	if err := v.SetCursor(cx, cy+1); err != nil {
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

// Move view cursor to the top
func moveViewCursorUp(g *gocui.Gui, v *gocui.View, dY int) error {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if cy > dY {
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// Get view line (relative to the cursor)
func getViewLine(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	return l, err
}

// Get the next view line (relative to the cursor)
func getNextViewLine(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy + 1); err != nil {
		l = ""
	}

	return l, err
}

// Set view cursor to line
func setViewCursorToLine(g *gocui.Gui, v *gocui.View, lines []string, selLine string) error {
	ox, _ := v.Origin()
	cx, _ := v.Cursor()
	for y, line := range lines {
		if line == selLine {
			if err := v.SetCursor(ox, y); err != nil {
				if err := v.SetOrigin(cx, y); err != nil {
					return err
				}
			}
		}
	}
	return nil
}




// Get docker name form line
func getConNameFromLine(line string) string {
	if line == "" {
		return ""
	}

	i := strings.Index(line, " ")
	if i == -1 {
		return line
	}

	return line[0:i]
}



// Get selected pod
func getSelectedCon(g *gocui.Gui) (string, error) {
	v, err := g.View("con")
	if err != nil {
		return "", err
	}
	l, err := getViewLine(g, v)
	if err != nil {
		return "", err
	}
	p := getConNameFromLine(l)

	return p, nil
}


// Show views logs
func showViewConLogs(g *gocui.Gui) error {
	vn := "logs"

	switch LOG_MOD {
	case "con":
		// Get current selected pod
		p, err := getSelectedCon(g)
		if err != nil {
			return err
		}

		lv, err := g.View(vn)
                if err != nil {
                        return err
                }
                lv.Clear()

		fmt.Fprintln(lv, "Container you choose is: " + p)
	}


//	debug(g, "Action: Show view logs")
	g.SetViewOnTop(vn)
	g.SetCurrentView(vn)

	return nil
}





























































































// Display error
func displayError(g *gocui.Gui, e error) error {
	lMaxX, lMaxY := g.Size()
	minX := lMaxX / 6
	minY := lMaxY / 6
	maxX := 5 * (lMaxX / 6)
	maxY := 5 * (lMaxY / 6)

	if v, err := g.SetView("errors", minX, minY, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Title = " ERROR "
		v.Frame = true
		v.Wrap = true
		v.Autoscroll = true
		v.BgColor = gocui.ColorRed
		v.FgColor = gocui.ColorWhite

		// Content
		v.Clear()
		fmt.Fprintln(v, e.Error())

		// Send to forground
		g.SetCurrentView(v.Name())
	}

	return nil
}

// Hide error box
func hideError(g *gocui.Gui) {
	g.DeleteView("errors")
}

// Display confirmation message
func displayConfirmation(g *gocui.Gui, m string) error {
	lMaxX, lMaxY := g.Size()

	if v, err := g.SetView("confirmation", -1, lMaxY-3, lMaxX, lMaxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = false

		// Content
		fmt.Fprintln(v, textPadCenter(m, lMaxX))

		// Auto-hide message
		hide := func() {
			hideConfirmation(g)
		}
		time.AfterFunc(time.Duration(2)*time.Second, hide)
	}

	return nil
}

// Hide confirmation message
func hideConfirmation(g *gocui.Gui) {
	g.DeleteView("confirmation")
}

