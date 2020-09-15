package ui

import (
	"github.com/jroimartin/gocui"
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


//Display Logs after Pod select
func actionViewConSelect(g *gocui.Gui, v *gocui.View) error {
        line,err  := getViewLine(g,v)
        if err != nil {
                return err
        }
//      maxX, maxY := g.Size()
        LOG_MOD = "con"
        errr := showViewConLogs(g)

        changeStatusContext(g, "SL")
//      viewLogs(g, maxX, maxY)
        displayConfirmation(g, line+" Pod selected")
        return errr

}
