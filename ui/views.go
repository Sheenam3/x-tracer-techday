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

func viewTcpLogs(g *gocui.Gui, lMaxX int, lMaxY int) error {
	if v, err := g.SetView("tcplogs", 1, 1, lMaxX/2, lMaxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = " Tcp Logs "
		v.Autoscroll = true
		v.Wrap = true

		v.SetCursor(1,3)


	}

	return nil
}


func viewTcpLifeLogs(g *gocui.Gui, lMaxX int, lMaxY int) error {
	if v, err := g.SetView("tcplife", lMaxX/2 , 1, lMaxX, lMaxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = " TcpLife "
		v.Autoscroll = true
		v.Wrap = true

		v.SetCursor(1,3)


	}

	return nil
}

func viewExecSnoopLogs(g *gocui.Gui, lMaxX int, lMaxY int) error {
	if v, err := g.SetView("execsnoop", 1, lMaxY/2 , lMaxX/2, lMaxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = " ExecSnoop "
		v.Autoscroll = true
		v.Wrap = true

		v.SetCursor(1,3)


	}

	return nil
}

func viewCacheStatLogs(g *gocui.Gui, lMaxX int, lMaxY int) error {
	if v, err := g.SetView("cachestat", lMaxX/2 , lMaxY/2, lMaxX, lMaxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "CacheStat"
		v.Autoscroll = true
		v.Wrap = true

		v.SetCursor(1,3)


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




func viewProbes(g *gocui.Gui, lMaxX int, lMaxY int) error {
	w := lMaxX / 2
	h := lMaxY / 4
	minX := (lMaxX / 2) - (w / 2)
	minY := (lMaxY / 2) - (h / 2)
	maxX := minX + w
	maxY := minY + h
	// Main view
	if v, err := g.SetView("probes", minX, minY, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// Configure view
		v.Title = " Select Probes "
		v.Frame = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		viewProbeNames(g)

	}
	return nil
}


func viewProbeNames(g *gocui.Gui){
	g.Update(func(g *gocui.Gui) error {
	
		v, err := g.View("probes")
		if err != nil {
			return err
		}


		probes := getProbeNames()

	//var pn []string

	v.Clear()
	
	if len(probes) >= 0 {
		for i, _ := range probes{
			fmt.Fprintln(v, probes[i])
		}
	}else {
	
	}
	
	setViewCursorToLine(g, v, probes, "tcptracer")
	
	return nil
	
	})

}


func getProbeNames()[]string{

	pn := []string {"tcptracer", "tcpconnect", "tcpaccept", "tcplife", "execsnoop", "biosnoop", "cachestat", "All Probes"}
	return pn

}

