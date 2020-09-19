package ui

import (
	"strings"
	"context"
	"github.com/jroimartin/gocui"
	"fmt"
//	"time"
	pp "github.com/Sheenam3/x-tracer-techday/parse"
"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

	conName,err := getSelectedCon(g)
//	pid := getPid(id)
	var conid string
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
//
//
//
        topResult, err := cli.ContainerTop(ctx, conid,/*containers[con].ID*/ []string{"o","pid"})

        if err != nil {
                panic(err)
        }

	if line == "All Probes"{
		G,lv := displayTcpLogs(g)
		G.SetViewOnTop("tcplogs")
        	G.SetCurrentView("tcplogs")

	//	lvv.Autoscroll = true
		logtcptracer := make(chan pp.Log, 1)
                go pp.RunTcptracer("tcptracer", logtcptracer, topResult.Processes[0][0])
                go func() {
                           for val := range logtcptracer {
                                  parse := strings.Fields(string(val.Fulllog))
                                  fmt.Fprintln(lv,"{Probe:" + "TCPTRACER" + "|" + "Sys_Time:" + parse[0]  + "|" + "T:" + parse[1]  + "|" + "PID:"  + parse[3]  + "|" + " PNAME:"  + parse[4]  + "|" + "IP->"  + parse[5]  + "|" + "SADDR:"  + parse[6]  + "|" + "DADDR:" + parse[7]  + "|" + "SPORT:" + parse[8]  + "|" + "DPORT:"  + parse[9])
                                }
                        }()

		logtcpconnect := make(chan pp.Log, 1)
                        go pp.RunTcpconnect("tcpconnect", logtcpconnect, topResult.Processes[0][0])
                        go func() {
                                for val := range logtcpconnect {
                                        parse := strings.Fields(string(val.Fulllog))
			fmt.Fprintln(lv,"{Probe:" + "TCPCONNECT" + "|" + "Sys_Time:" + parse[0]  + "|" + "T:" + parse[1]  + "|" + "PID:"  + parse[3]  + "|" + " PNAME:"  + parse[4]  + "|" + "IP->"  + parse[5]  + "|" + "SADDR:" + parse[6]  + "|" + "DADDR:" + parse[7]  + "|" + "DPORT:"  + parse[8])
                                }
                        }()

		logtcpaccept := make(chan pp.Log, 1)
                        go pp.RunTcpaccept("tcpaccept", logtcpaccept, topResult.Processes[0][0])
                        go func() {
                                for val := range logtcpaccept {
                                        parse := strings.Fields(string(val.Fulllog))
                                  fmt.Fprintln(lv,"{Probe:" + "TCPACCEPT" + "|" + "Sys_Time:" + parse[0]  + "|" + "T:" + parse[1]  + "|" + "PID:"  + parse[3]  + "|" + " PNAME:"  + parse[4]  + "|" + "IP->"  + parse[5]  + "|" + "RADDR:"  + parse[6]  + "|" + "RPORT:" + parse[7]  + "|" + "LADDR:" + parse[8]  + "|" + "LPORT:"  + parse[9])
                                }
                        }()


		e,ev := displayExecLogs(g)
                e.SetViewOnTop("execsnoop")
                e.SetCurrentView("execsnoop")
//		ev.Autoscroll = true
		logexecsnoop := make(chan pp.Log, 1)
                        go pp.RunExecsnoop("execsnoop", logexecsnoop, topResult.Processes[0][0])
                        go func() {
                                for val := range logexecsnoop {
                                        parse := strings.Fields(string(val.Fulllog))
			fmt.Fprintln(ev,"Sys_Time:" + parse[0]  + "|" + "T:" + parse[1]  + "|" + "PID:"  + parse[4]  + "|" + " PNAME:"  + parse[3]  + "|" + "PPID->"  + parse[5]  + "|" + "RET:" + parse[6]  + "|" + "ARGS:" + parse[7])
                                }
                        }()

		c,cv := displayCacheLogs(g)
                c.SetViewOnTop("cachestat")
                c.SetCurrentView("cachestat")
//		cv.Autoscroll = true
		logcachetop := make(chan pp.Log, 1)
                        go pp.RunCachetop("cachestat", logcachetop, topResult.Processes[0][0])
                        go func() {
                                for val := range logcachetop {
                                        parse := strings.Fields(string(val.Fulllog))
					fmt.Fprintln(cv,"{Probe:" + "CACHESTAT" + "|" + "Sys_Time:" + parse[0]  + "|" + "PID:"  + parse[1]  + "|" + "UID:"  + parse[2]  + "|" + "CMD->"  + parse[3]  + "|" + "HITS:"  + parse[5]  + "|" + "MISS:" + parse[6]  + "|" + "DIRTIES:" + parse[7]  + "|" + "READ_HIT%:"  + parse[8]  +  "W_HIT%:"  + parse[9])
                                }
                        }()

		logbiosnoop := make(chan pp.Log, 1)
                        go pp.RunBiosnoop("biosnoop", logbiosnoop, topResult.Processes[0][0])
                        go func() {
                                for val := range logbiosnoop {
                                        parse := strings.Fields(string(val.Fulllog))
                                       	fmt.Fprintln(cv,"{Probe:" + "BIOSNOOP" + "|" + "Sys_Time:" + parse[0]  + "|" + "T:"  + parse[1]  + "|"  +  "PNAME:"  + parse[2] + "|" + "PID:"  + parse[3]  + "|" + "DISK->"  + parse[4]  + "|" + "R/W:"  + parse[5]  + "|" + "SECTOR:" + parse[6]  + "|" + "BYTES:" + parse[7]  + "|" + "LAT(ms):"  + parse[9])
                                }
                        }()

		t,tv := displayTcplifeLogs(g)
                t.SetViewOnTop("tcplife")
                t.SetCurrentView("tcplife")
//		tv.Autoscroll = true
		logtcplife := make(chan pp.Log, 1)
                        go pp.RunTcplife("tcplife", logtcplife, topResult.Processes[0][0])
                        go func() {
                                for val := range logtcplife {
                                        parse := strings.Fields(string(val.Fulllog))
					fmt.Fprintln(tv,"{Probe:" + "TCPLIFE" + "|" + "Sys_Time:" + parse[0]  + "|" + "PID:"  + parse[2]  + "|" + "PNAME:"  + parse[3]  + "|" + "LADDR"  + parse[4]  + "|" + "LPORT:"  + parse[5]  + "|" + "RADDR:" + parse[6]  + "|" + "RPORT:" + parse[7]  + "|" + "TX_KB"  + parse[8]  +  "RX_KB:"  + parse[9] + "|" + "MS:" + parse[10])
                                }
                        }()


//	        fmt.Fprintln(lv,  "TCP LOGS HERE")
//	        fmt.Fprintln(ev,  topResult.Processes[0][0])


	}else{
		
	        G,p,lv := showViewConLogs(g)
		displayConfirmation(g, line+" probe selected")
		startAgent(G,p,lv,line)
		G.SetViewOnTop("logs")
		G.SetCurrentView("logs")
	}

/*for {
                time.Sleep(time.Duration(1) * time.Second)
}*/

	return err
}

