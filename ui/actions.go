package ui

import (
//	"strings"
//	"context"
	"github.com/jroimartin/gocui"
	"fmt"
//	"time"
//	pp "github.com/Sheenam3/x-tracer-techday/parse"
//"github.com/docker/docker/api/types"
//	"github.com/docker/docker/client"
)



func actionGlobalQuit(g *gocui.Gui, v *gocui.View) error {

	return gocui.ErrQuit
}


// View pods: Up
func actionViewConUp(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorUp(g, v, 2)
//	debug(g, "Select up in pods view")
	return nil
}

// View pods: Down
func actionViewConDown(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorDown(g, v, false)
//	debug(g, "Select down in pods view")
	return nil
}


//Display Probe Tools after Pod select
func actionViewConSelect(g *gocui.Gui, v *gocui.View) error {
        line,err  := getViewLine(g,v)
        if err != nil {
                return err
        }
        LOG_MOD = "con"
        errr := showSelectProbe(g)

        changeStatusContext(g, "SL")
        displayConfirmation(g, line+" Pod selected")
        return errr

}



// Probes


// View Probes: Up
func actionViewProbesUp(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorUp(g, v, 0)
	//debug(g, "Select up in namespaces view")
	return nil
}

// View Probes: Down
func actionViewProbesDown(g *gocui.Gui, v *gocui.View) error {
	moveViewCursorDown(g, v, false)
	//debug(g, "Select down in namespaces view")
	return nil
}


func actionViewProbesSelect(g *gocui.Gui, v *gocui.View) error {

	line, err := getViewLine(g, v)
  	LOG_MOD = "probe"

//	id := getContainerId(g)
//	pid := getPid(id)
//	var conid string
//	ctx := context.Background()
//        cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
//        if err != nil {
//                panic(err)
//        }
//
//        containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
//        for _, container := range containers {
//                out := strings.TrimLeft(container.Names[0],"/")
//                if line == out{
//                        conid = container.ID
//
//                }
//
//        }
//
//
//
//        topResult, err := cli.ContainerTop(context.Background(), conid,/*containers[con].ID*/ []string{"o","pid"})
//
//        if err != nil {
//                panic(err)
//        }

	if line == "All Probes"{
		G,lv := displayTcpLogs(g)
		G.SetViewOnTop("tcplogs")
        	G.SetCurrentView("tcplogs")
//		logtcptracer := make(chan pp.Log, 1)
//                go pp.RunTcptracer("tcptracer", logtcptracer, topResult.Processes[0][0])
//                go func() {
//                           for val := range logtcptracer {
//                                  parse := strings.Fields(string(val.Fulllog))
//                                  fmt.Fprintln(lv,"{Probe:%s |Sys_Time: %s |T: %s | PID:%s | PNAME:%s |IP->%s | SADDR:%s | DADDR:%s | SPORT:%s | DPORT:%s \n","tcptracer",parse[0],parse[1],parse[3],parse[4],parse[5],parse[6],parse[7],parse[8],parse[9])
//                                }
//                        }()

		e,ev := displayExecLogs(g)
                e.SetViewOnTop("execsnoop")
                e.SetCurrentView("execsnoop")

		c,cv := displayCacheLogs(g)
                c.SetViewOnTop("cachestat")
                c.SetCurrentView("cachestat")

		t,tv := displayTcplifeLogs(g)
                t.SetViewOnTop("tcplife")
                t.SetCurrentView("tcplife")


	        fmt.Fprintln(lv,  "TCP LOGS HERE")
	        fmt.Fprintln(ev,  "hi")
	        fmt.Fprintln(cv,  "p")
	        fmt.Fprintln(tv,  "s")

	}else{
		id := "0"
	        G,p,lv,id := showViewConLogs(g)
		displayConfirmation(g, line+" probe selected")
		startAgent(G,p,lv,line,id)
		G.SetViewOnTop("logs")
		G.SetCurrentView("logs")
	}



	return err
}

