package ui


import (
	"fmt"
//	"exec"
	"context"
	"log"
	"os/exec"
	"io"
	"strings"
	"time"
	"github.com/jroimartin/gocui"
"github.com/docker/docker/api/types"
        "github.com/docker/docker/client"
pp "github.com/Sheenam3/x-tracer-techday/parse"
)


var (
	g	*gocui.Gui
)


var version = "master"
var LOG_MOD string = "con"

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

	Key{"probes", gocui.KeyArrowUp, actionViewProbesUp},
	Key{"probes", gocui.KeyArrowDown, actionViewProbesDown},
	Key{"probes", gocui.KeyEnter, actionViewProbesSelect},
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

	viewTcpLogs(g, maxX, maxY)
	viewTcpLifeLogs(g, maxX, maxY)
	viewExecSnoopLogs(g, maxX, maxY)
	viewCacheStatLogs(g, maxX, maxY)
	viewLogs(g, maxX, maxY)
	viewProbes(g, maxX, maxY)
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


func showSelectProbe(g *gocui.Gui) error {

	switch LOG_MOD {
	case "con":
		//Choose probe tool
		g.SetViewOnTop("probes")
		g.SetCurrentView("probes")
		changeStatusContext(g,"SE")
	}
	return nil
}

func getContainerId(g *gocui.Gui) (string){
	p, err := getSelectedCon(g)
                if err != nil {
                        fmt.Println(err)
                }
	id, err := exec.Command("sudo", "docker", "ps", "--no-trunc", "-aqf", fmt.Sprintf("name=%s",p)).Output() 
		if err != nil {
		    log.Fatal(err)
		}

	conId := string(id)

//
//	pid, err := exec.Command("sudo", "docker", "inspect", "-f", "'{{.State.Pid}}'", fmt.Sprintf("%s",conId)).Output()
//        if err != nil {
//           fmt.Println("kya h:",err)
//        }
//        ppid := string(pid)
//        out := strings.TrimLeft(strings.TrimRight(ppid,"'"),"'")
//	var s string
//	if len(out) > 0 {
//        	s = out[:len(out)-2]
//	}
        


	return conId

}

// Show views logs
func showViewConLogs(g *gocui.Gui) (*gocui.Gui,string,io.Writer){//,string) {
	vn := "logs"

	switch LOG_MOD {
	case "probe":
		// Get current selected pod
		p, err := getSelectedCon(g)
		if err != nil {
			fmt.Println(err)
		}

		lv, err := g.View(vn)
                if err != nil {
                        fmt.Println(err)
                }
                lv.Clear()

		id, err := exec.Command("sudo", "docker", "ps", "--no-trunc", "-aqf", fmt.Sprintf("name=%s",p)).Output() 
		if err != nil {
		    log.Fatal(err)
		}

		conId := string(id)
		fmt.Fprintln(lv, "Container you choose is: " + p)
		fmt.Fprintln(lv, "Container ID:", conId)
		return g,p,lv
	}



//	g.SetViewOnTop(vn)
//	g.SetCurrentView(vn)

	return nil,"ok",nil
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


func getPid(conId string)(string){

	id, err := exec.Command("sudo", "docker", "inspect", "-f", "'{{.State.Pid}}'", fmt.Sprintf("%s",conId)).Output()
        if err != nil {
           fmt.Println("kya h:",err)
        }
	cid := string(id)
        out := strings.TrimLeft(strings.TrimRight(cid,"'"),"'")
        s := out[:len(out)-2]

	return s

}


func startAgent(g *gocui.Gui, conName string, o io.Writer, probeName string) error {


        var conid string

        conid = getContainerId(g)
	ctx := context.Background()
        cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
        if err != nil {
                panic(err)
        }

        containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
        for _, container := range containers {
                out := strings.TrimLeft(container.Names[0],"/")
              if conName == out{
                      conid = container.ID

              }

        }

	

        topResult, err := cli.ContainerTop(context.Background(), conid,/*containers[con].ID*/ []string{"o","pid"})




//	fmt.Fprintln(o, topResult.Processes[0][0])
	fmt.Fprintln(o, "Probe Selected is --->",probeName + "\n" + probeName + " logs will be streamed shortly..." )

	switch probeName {

        case "tcptracer":
		logtcptracer := make(chan pp.Log, 1)
                go pp.RunTcptracer("tcptracer", logtcptracer, topResult.Processes[0][0])
                go func() {
                           for val := range logtcptracer {
                                  parse := strings.Fields(string(val.Fulllog))
                                  fmt.Fprintln(o,"{Probe:" + "TCPTRACER" + "|" + "Sys_Time:" + parse[0]  + "|" + "T:" + parse[1]  + "|" + "PID:"  + parse[3]  + "|" + " PNAME:"  + parse[4]  + "|" + "IP->"  + parse[5]  + "|" + "SADDR:"  + parse[6]  + "|" + "DADDR:" + parse[7]  + "|" + "SPORT:" + parse[8]  + "|" + "DPORT:"  + parse[9])
                                }
                        }()
        case "tcpconnect":
                logtcpconnect := make(chan pp.Log, 1)
                        go pp.RunTcpconnect("tcpconnect", logtcpconnect, topResult.Processes[0][0])
                        go func() {
                                for val := range logtcpconnect {
                                        parse := strings.Fields(string(val.Fulllog))
			fmt.Fprintln(o,"{Probe:" + "TCPCONNECT" + "|" + "Sys_Time:" + parse[0]  + "|" + "T:" + parse[1]  + "|" + "PID:"  + parse[3]  + "|" + " PNAME:"  + parse[4]  + "|" + "IP->"  + parse[5]  + "|" + "SADDR:" + parse[6]  + "|" + "DADDR:" + parse[7]  + "|" + "DPORT:"  + parse[8])
                                }
                        }()

	case "tcpaccept":
		logtcpaccept := make(chan pp.Log, 1)
                        go pp.RunTcpaccept("tcpaccept", logtcpaccept, topResult.Processes[0][0])
                        go func() {
                                for val := range logtcpaccept {
                                        parse := strings.Fields(string(val.Fulllog))
                                  fmt.Fprintln(o,"{Probe:" + "TCPACCEPT" + "|" + "Sys_Time:" + parse[0]  + "|" + "T:" + parse[1]  + "|" + "PID:"  + parse[3]  + "|" + " PNAME:"  + parse[4]  + "|" + "IP->"  + parse[5]  + "|" + "RADDR:"  + parse[6]  + "|" + "RPORT:" + parse[7]  + "|" + "LADDR:" + parse[8]  + "|" + "LPORT:"  + parse[9])
                                }
                        }()

	case "tcplife":
		logtcplife := make(chan pp.Log, 1)
                        go pp.RunTcplife("tcplife", logtcplife, topResult.Processes[0][0])
                        go func() {
                                for val := range logtcplife {
                                        parse := strings.Fields(string(val.Fulllog))
					fmt.Fprintln(o,"{Probe:" + "TCPLIFE" + "|" + "Sys_Time:" + parse[0]  + "|" + "PID:"  + parse[2]  + "|" + "PNAME:"  + parse[3]  + "|" + "LADDR"  + parse[4]  + "|" + "LPORT:"  + parse[5]  + "|" + "RADDR:" + parse[6]  + "|" + "RPORT:" + parse[7]  + "|" + "TX_KB"  + parse[8]  +  "RX_KB:"  + parse[9] + "|" + "MS:" + parse[10])
                                }
                        }()


	case "execsnoop":
		logexecsnoop := make(chan pp.Log, 1)
                        go pp.RunExecsnoop("execsnoop", logexecsnoop, topResult.Processes[0][0])
                        go func() {
                                for val := range logexecsnoop {
                                        parse := strings.Fields(string(val.Fulllog))
			fmt.Fprintln(o,"Sys_Time:" + parse[0]  + "|" + "T:" + parse[1]  + "|" + "PID:"  + parse[4]  + "|" + " PNAME:"  + parse[3]  + "|" + "PPID->"  + parse[5]  + "|" + "RET:" + parse[6]  + "|" + "ARGS:" + parse[7])
                                }
                        }()

	case "cachestat":
		logcachetop := make(chan pp.Log, 1)
                        go pp.RunCachetop("cachestat", logcachetop, topResult.Processes[0][0])
                        go func() {
                                for val := range logcachetop {
                                        parse := strings.Fields(string(val.Fulllog))
					fmt.Fprintln(o,"{Probe:" + "CACHESTAT" + "|" + "Sys_Time:" + parse[0]  + "|" + "PID:"  + parse[1]  + "|" + "UID:"  + parse[2]  + "|" + "CMD->"  + parse[3]  + "|" + "HITS:"  + parse[5]  + "|" + "MISS:" + parse[6]  + "|" + "DIRTIES:" + parse[7]  + "|" + "READ_HIT%:"  + parse[8]  +  "W_HIT%:"  + parse[9])
                                }
                        }()

	case "biosnoop":
		logbiosnoop := make(chan pp.Log, 1)
                        go pp.RunBiosnoop("biosnoop", logbiosnoop, topResult.Processes[0][0])
                        go func() {
                                for val := range logbiosnoop {
                                        parse := strings.Fields(string(val.Fulllog))
                                       	fmt.Fprintln(o,"{Probe:" + "BIOSNOOP" + "|" + "Sys_Time:" + parse[0]  + "|" + "T:"  + parse[1]  + "|"  +  "PNAME:"  + parse[2] + "|" + "PID:"  + parse[3]  + "|" + "DISK->"  + parse[4]  + "|" + "R/W:"  + parse[5]  + "|" + "SECTOR:" + parse[6]  + "|" + "BYTES:" + parse[7]  + "|" + "LAT(ms):"  + parse[9])
                                }
                        }()
	}


		//displayLogs(g)

	

	return nil
}


func displayTcpLogs(g *gocui.Gui)(*gocui.Gui,io.Writer){

		lv, err := g.View("tcplogs")
                if err != nil {
                        fmt.Println(err)
                }
                lv.Clear()

		return g,lv

//		fmt.Fprintln(lv, "TCP Logs Here")


}


func displayTcplifeLogs(g *gocui.Gui)(*gocui.Gui,io.Writer){

		lv, err := g.View("tcplife")
                if err != nil {
                        fmt.Println(err)
                }
                lv.Clear()

		return g,lv

//		fmt.Fprintln(lv, "TCP Logs Here")


}

func displayExecLogs(g *gocui.Gui)(*gocui.Gui,io.Writer){

		lv, err := g.View("execsnoop")
                if err != nil {
                        fmt.Println(err)
                }
                lv.Clear()

		return g,lv

//		fmt.Fprintln(lv, "TCP Logs Here")


}


func displayCacheLogs(g *gocui.Gui)(*gocui.Gui,io.Writer){

		lv, err := g.View("cachestat")
                if err != nil {
                        fmt.Println(err)
                }
                lv.Clear()

		return g,lv

//		fmt.Fprintln(lv, "TCP Logs Here")


}
