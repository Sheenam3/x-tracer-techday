package ui

import (

	"context"
	"fmt"
	"strings"
//	"time"
	"github.com/docker/docker/api/types"
        "github.com/docker/docker/client"
	"github.com/jroimartin/gocui"
	"github.com/willf/pad"
)

// View: Overlay
func viewOverlay(g *gocui.Gui, lMaxX int, lMaxY int) error {
	if v, err := g.SetView("overlay", 0, 0, lMaxX, lMaxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = false
	}

	return nil
}

// View: Title bar
func viewTitle(g *gocui.Gui, lMaxX int, lMaxY int) error {
	if v, err := g.SetView("title", -1, -1, lMaxX, 1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = true
		v.BgColor = gocui.ColorDefault | gocui.AttrReverse
		v.FgColor = gocui.ColorDefault | gocui.AttrReverse

		// Content
		fmt.Fprintln(v, versionTitle(lMaxX))
		
	}

	return nil
}

func viewLogs(g *gocui.Gui, lMaxX int, lMaxY int) error {
	if v, err := g.SetView("logs", 2, 2, lMaxX-4, lMaxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Title = " Logs "
		v.Autoscroll = true
	}


	return nil
}


// View: Pods
func viewCon(g *gocui.Gui, lMaxX int, lMaxY int) error {
	if v, err := g.SetView("con", -1, 5, lMaxX, lMaxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = true
		v.Title = " Containers "
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 2)

		// Set as current view
		g.SetCurrentView(v.Name())

		// Content
		go viewPodsShowWithAutoRefresh(g)
	}

	return nil
}

// Auto refresh view pods
func viewPodsShowWithAutoRefresh(g *gocui.Gui) {
	go viewPodsRefreshList(g)
//	for {
		//	debug(g, fmt.Sprintf("View pods: Refreshing (%ds)", c.frequency))
//			go viewPodsRefreshList(g)
//	   }
}



func viewPodsRefreshList(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		lMaxX, _ := g.Size()
		//debug(g, "View pods: Actualize")
		v, err := g.View("con")
		if err != nil {
			return err
		}

		ctx := context.Background()
        	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
        	if err != nil {
                	panic(err)
        	}

        	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
        	if err != nil {
                	displayError(g,err)
			return nil
       		 }
		hideError(g)

		v.Clear()
		viewConAddLine(v, lMaxX, "NAME")
		fmt.Fprintln(v, strings.Repeat("─", lMaxX))


		if len(containers) > 0 {

	        	for _, container := range containers {
        	        	out := strings.TrimLeft(container.Names[0],"/")
                		viewConAddLine(v, lMaxX, out)
        		}

			if l, err := getViewLine(g, v); err != nil || l == "" {
					v.SetCursor(0, 2)
			}
		} else {
			v.SetCursor(0, 2)

			}

		return nil
	})

}



// View: Status bar
func viewStatusBar(g *gocui.Gui, lMaxX int, lMaxY int) error {
	if v, err := g.SetView("status", -1, lMaxY-2, lMaxX, lMaxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Settings
		v.Frame = false
		v.BgColor = gocui.ColorBlack
		v.FgColor = gocui.ColorWhite

		// Content
		changeStatusContext(g, "D")
	}

	return nil
}


// Change status context
func changeStatusContext(g *gocui.Gui, c string) error {
	lMaxX, _ := g.Size()
	v, err := g.View("status")
	if err != nil {
		return err
	}

	v.Clear()

	i := lMaxX + 4
	b := ""

	switch c {
	case "D":
		i = 150 + i
		b = b + frameText("↑") + " Up   "
		b = b + frameText("↓") + " Down   "
		b = b + frameText("D") + " Delete   "
		b = b + frameText("L") + " Show Logs   "

	case "SE":
		i = i + 100
		b = b + frameText("↑") + " Up   "
		b = b + frameText("↓") + " Down   "
		b = b + frameText("Enter") + " Select   "
	case "SL":
		i = i + 100
		b = b + frameText("↑") + " Up   "
		b = b + frameText("↓") + " Down   "
		b = b + frameText("L") + " Hide Logs   "
	}
	b = b + frameText("CTRL+C") + " Exit"

	fmt.Fprintln(v, pad.Left(b, i, " "))

	return nil
}


func viewConAddLine(v *gocui.View, maxX int, name string) {
	wN := maxX - 34 // 54 // TODO CPU + Memory #20
	if wN < 45 {
		wN = 45
	}
	line := pad.Right(name, wN, " ") 
		//pad.Right(cpu, 10, " ") + // TODO CPU + Memory #20
		//pad.Right(memory, 10, " ") + // TODO CPU + Memory #20
//		pad.Right(ready, 10, " ") +
//		pad.Right(status, 10, " ") +
//		pad.Right(restarts, 10, " ") +
//		pad.Right(age, 4, " ")
	fmt.Fprintln(v, line)
}
