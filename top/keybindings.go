// Define key bindings.

package top

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/lesovsky/pgcenter/lib/stat"
	"os"
	"strconv"
	"strings"
	"time"
)

// Key represents particular key, a view where it should work and associated function.
type key struct {
	viewname string
	key      interface{}
	handler  func(g *gocui.Gui, v *gocui.View) error
}

// Setup key bindings and handlers.
func keybindings(g *gocui.Gui) error {
	var keys = []key{
		{"", gocui.KeyCtrlC, quit},
		{"", gocui.KeyCtrlQ, quit},
		{"sysstat", gocui.KeyArrowLeft, orderKeyLeft},
		{"sysstat", gocui.KeyArrowRight, orderKeyRight},
		{"sysstat", '<', switchSortOrder},
		{"sysstat", ',', toggleSysTables},
		{"sysstat", 'I', toggleIdleConns},
		{"sysstat", 'd', switchContextTo(stat.DatabaseView)},
		{"sysstat", 'r', switchContextTo(stat.ReplicationView)},
		{"sysstat", 't', switchContextTo(stat.TablesView)},
		{"sysstat", 'i', switchContextTo(stat.IndexesView)},
		{"sysstat", 's', switchContextTo(stat.SizesView)},
		{"sysstat", 'f', switchContextTo(stat.FunctionsView)},
		{"sysstat", 'v', switchContextTo(stat.VacuumView)},
		{"sysstat", 'a', switchContextTo(stat.ActivityView)},
		{"sysstat", 'x', switchContextTo(stat.StatementsView)},
		{"sysstat", 'Q', resetStat},
		{"sysstat", 'E', menuOpen(menuConfStyle)},
		{"sysstat", 'X', menuOpen(menuPgssStyle)},
		{"sysstat", 'l', showPgLog},
		{"sysstat", 'C', showPgConfig},
		{"sysstat", 'p', runPsql},
		{"sysstat", 'B', showAux(AUX_DISKSTAT)},
		{"sysstat", 'N', showAux(AUX_NICSTAT)},
		{"sysstat", 'L', showAux(AUX_LOGTAIL)},
		{"sysstat", 'R', dialogOpen(dialogPgReload)},
		{"sysstat", '/', dialogOpen(dialogFilter)},
		{"sysstat", '-', dialogOpen(dialogCancelQuery)},
		{"sysstat", '_', dialogOpen(dialogTerminateBackend)},
		{"sysstat", 'n', dialogOpen(dialogSetMask)},
		{"sysstat", 'm', showBackendMask},
		{"sysstat", 'k', dialogOpen(dialogCancelGroup)},
		{"sysstat", 'K', dialogOpen(dialogTerminateGroup)},
		{"sysstat", 'A', dialogOpen(dialogChangeAge)},
		{"sysstat", 'G', dialogOpen(dialogQueryReport)},
		{"sysstat", 'z', dialogOpen(dialogChangeRefresh)},
		{"dialog", gocui.KeyEsc, dialogCancel},
		{"dialog", gocui.KeyEnter, dialogFinish},
		{"menu", gocui.KeyEsc, menuClose},
		{"menu", gocui.KeyArrowUp, moveCursor(moveUp)},
		{"menu", gocui.KeyArrowDown, moveCursor(moveDown)},
		{"menu", gocui.KeyEnter, menuSelect},
		{"sysstat", 'h', showHelp},
		{"sysstat", gocui.KeyF1, showHelp},
		{"help", gocui.KeyEsc, closeHelp},
	}

	g.InputEsc = true

	for _, k := range keys {
		if err := g.SetKeybinding(k.viewname, k.key, gocui.ModNone, k.handler); err != nil {
			return fmt.Errorf("ERROR: failed to setup keybindings: %s", err)
		}
	}

	return nil
}

// Change interval of stats refreshing.
func changeRefresh(g *gocui.Gui, v *gocui.View, answer string) {
	answer = strings.TrimPrefix(string(v.Buffer()), dialogPrompts[dialogChangeRefresh])
	answer = strings.TrimSuffix(answer, "\n")

	if answer == "" {
		printCmdline(g, "Do nothing. Empty input.")
		return
	}

	interval, _ := strconv.Atoi(answer)

	switch {
	case interval < 1:
		printCmdline(g, "Should not be less than 1 second.")
		return
	case interval > 300:
		printCmdline(g, "Should not be more than 300 seconds.")
		return
	}

	refreshInterval = time.Duration(interval) * refreshMinGranularity
	do_update <- 1
}

// Quit program.
func quit(g *gocui.Gui, _ *gocui.View) error {
	close(do_update)
	close(do_exit)
	g.Close()
	os.Exit(0)
	return gocui.ErrQuit
}
