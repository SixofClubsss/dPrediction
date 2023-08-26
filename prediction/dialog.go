package prediction

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ownerObjects struct {
	synced  bool
	service struct {
		run     *widget.Button
		payouts struct {
			check   *widget.Check
			enabled bool
		}
		transactions struct {
			check   *widget.Check
			enabled bool
		}
	}
	sports struct {
		end     *dwidget.DeroAmts
		amt     *dwidget.DeroAmts
		game    *widget.Select
		league  *widget.SelectEntry
		feed    *widget.SelectEntry
		deposit *dwidget.DeroAmts
		set     *widget.Button
		cancel  *widget.Button
		payout  *widget.SelectEntry
	}
	predict struct {
		end     *dwidget.DeroAmts
		mark    *widget.Entry
		amt     *dwidget.DeroAmts
		name    *widget.SelectEntry
		feed    *widget.SelectEntry
		deposit *dwidget.DeroAmts
		set     *widget.Button
		post    *widget.Button
		pay     *widget.Button
		cancel  *widget.Button
	}
}

var owner ownerObjects

// Check if prediction is for on chain values
func isOnChainPrediction(s string) bool {
	if s == "DERO-Difficulty" || s == "DERO-Block Time" || s == "DERO-Block Number" {
		return true
	}

	return false
}

// Check which on chain values are required
func onChainPrediction(s string) int {
	switch s {
	case "DERO-Difficulty":
		return 1
	case "DERO-Block Time":
		return 2
	case "DERO-Block Number":
		return 3
	default:
		return 0
	}
}

// dPrediction owner control objects for side menu
//   - Pass side menu window to reset to
func predictionOpts(window fyne.Window) fyne.CanvasObject {
	pred := []string{"DERO-BTC", "XMR-BTC", "BTC-USDT", "DERO-USDT", "XMR-USDT", "DERO-Difficulty", "DERO-Block Time", "DERO-Block Number"}
	owner.predict.name = widget.NewSelectEntry(pred)
	owner.predict.name.SetPlaceHolder("Name:")
	owner.predict.name.OnChanged = func(s string) {
		if isOnChainPrediction(s) {
			opts := []string{rpc.DAEMON_RPC_REMOTE1, rpc.DAEMON_RPC_REMOTE2, rpc.DAEMON_RPC_REMOTE5, rpc.DAEMON_RPC_REMOTE6}
			owner.predict.feed.SetOptions(opts)
			if owner.predict.feed.Text != opts[1] {
				owner.predict.feed.SetText(opts[0])
			}
			owner.predict.feed.SetPlaceHolder("Node:")
			owner.predict.feed.Refresh()
		} else {
			opts := []string{"dReams Client"}
			owner.predict.feed.SetOptions(opts)
			owner.predict.feed.SetText(opts[0])
			owner.predict.feed.SetPlaceHolder("Feed:")
			owner.predict.feed.Refresh()
		}
	}

	owner.predict.end = dwidget.NewDeroEntry("", 1, 0)
	owner.predict.end.SetPlaceHolder("Closes At:")
	owner.predict.end.AllowFloat = false
	owner.predict.end.Validator = validation.NewRegexp(`^\d{10,}$`, "Unix time required")

	owner.predict.mark = widget.NewEntry()
	owner.predict.mark.SetPlaceHolder("Mark:")
	owner.predict.mark.Validator = validation.NewRegexp(`^\d{1,}$`, "Int required")

	owner.predict.amt = dwidget.NewDeroEntry("", 0.1, 1)
	owner.predict.amt.SetPlaceHolder("Minimum Amount:")
	owner.predict.amt.AllowFloat = true
	owner.predict.amt.Wrapping = fyne.TextTruncate
	owner.predict.amt.Validator = validation.NewRegexp(`^\d{1,}\.\d{1,5}$|^[^0]\d{0,}$`, "Int or float required")

	feeds := []string{"dReams Client"}
	owner.predict.feed = widget.NewSelectEntry(feeds)
	owner.predict.feed.SetPlaceHolder("Feed:")

	owner.predict.deposit = dwidget.NewDeroEntry("", 0.1, 1)
	owner.predict.deposit.SetPlaceHolder("Deposit Amount:")
	owner.predict.deposit.AllowFloat = true
	owner.predict.deposit.Wrapping = fyne.TextTruncate
	owner.predict.deposit.Validator = validation.NewRegexp(`^\d{1,}\.\d{1,5}$|^[^0]\d{0,}$`, "Int or float required")

	reset := window.Content().(*fyne.Container).Objects[2]

	owner.predict.set = widget.NewButton("Set Prediction", func() {
		if owner.predict.deposit.Validate() == nil && owner.predict.amt.Validate() == nil && owner.predict.end.Validate() == nil && owner.predict.mark.Validate() == nil {
			if len(Predict.Contract.SCID) == 64 {
				window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(2, 100, window, reset))
				window.Content().(*fyne.Container).Objects[2].Refresh()
				return
			}

			dialog.NewInformation("Prediction", "Select a valid contract", window).Show()
		}
	})

	owner.predict.cancel = widget.NewButton("Cancel", func() {
		window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(8, 0, window, reset))
		window.Content().(*fyne.Container).Objects[2].Refresh()
	})

	owner.predict.cancel.Hide()

	owner.predict.post = widget.NewButton("Post", func() {
		go SetPredictionPrices(rpc.Daemon.Connect)
		var a float64
		prediction := Predict.prediction
		if isOnChainPrediction(prediction) {
			switch onChainPrediction(prediction) {
			case 1:
				a = rpc.GetDifficulty(Predict.feed)
				window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(6, a, window, reset))
			case 2:
				a = rpc.GetBlockTime(Predict.feed)
				window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(6, a, window, reset))
			case 3:
				d := rpc.DaemonHeight("Prediction", Predict.feed)
				a = float64(d)
				window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(6, a, window, reset))
			default:

			}

		} else {
			a, _ = menu.GetPrice(prediction, "Prediction")
			window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(4, a, window, reset))
		}

		window.Content().(*fyne.Container).Objects[2].Refresh()

	})

	owner.predict.post.Hide()

	owner.predict.pay = widget.NewButton("Prediction Payout", func() {
		go SetPredictionPrices(rpc.Daemon.Connect)
		var a float64
		prediction := Predict.prediction
		if isOnChainPrediction(prediction) {
			switch onChainPrediction(prediction) {
			case 1:
				a = rpc.GetDifficulty(Predict.feed)
				window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(7, a, window, reset))
			case 2:
				a = rpc.GetBlockTime(Predict.feed)
				window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(7, a, window, reset))
			case 3:
				d := rpc.DaemonHeight("Prediction", Predict.feed)
				a = float64(d)
				window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(7, a, window, reset))
			default:

			}

		} else {
			a, _ = menu.GetPrice(prediction, "Prediction")
			window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(5, a, window, reset))
		}

		window.Content().(*fyne.Container).Objects[2].Refresh()
	})

	owner.predict.pay.Hide()

	owner_p := container.NewVBox(
		humanTimeConvert(),
		layout.NewSpacer(),
		owner.predict.name,
		owner.predict.end,
		owner.predict.mark,
		owner.predict.amt,
		owner.predict.feed,
		owner.predict.deposit,
		owner.predict.set,
		layout.NewSpacer(),
		owner.predict.cancel,
		layout.NewSpacer(),
		owner.predict.post,
		layout.NewSpacer(),
		owner.predict.pay,
		layout.NewSpacer(),
	)

	return owner_p
}

// dSports owner control objects for side menu
//   - Pass side menu window to reset to
func sportsOpts(window fyne.Window) fyne.CanvasObject {
	options := []string{}
	owner.sports.game = widget.NewSelect(options, func(s string) {
		var date string
		game := strings.Split(s, "   ")
		for i := range s {
			if i > 3 {
				date = s[0:10]
			}
		}
		comp := date[0:4] + date[5:7] + date[8:10]
		GetGameEnd(comp, game[1], owner.sports.league.Text)
	})
	owner.sports.game.PlaceHolder = "Game:"

	leagues := []string{"EPL", "MLS", "NBA", "NFL", "NHL", "MLB", "Bellator", "UFC"}
	owner.sports.league = widget.NewSelectEntry(leagues)
	owner.sports.league.OnChanged = func(s string) {
		owner.sports.game.Options = []string{}
		owner.sports.game.Selected = ""
		if s == "Bellator" || s == "UFC" {
			owner.sports.game.PlaceHolder = "Fight:"
		} else {
			owner.sports.game.PlaceHolder = "Game:"
		}
		owner.sports.game.Refresh()
		switch s {
		case "EPL":
			go GetCurrentWeek("EPL")
		case "MLS":
			go GetCurrentWeek("MLS")
		case "NBA":
			go GetCurrentWeek("NBA")
		case "NFL":
			go GetCurrentWeek("NFL")
		case "NHL":
			go GetCurrentWeek("NHL")
		case "MLB":
			go GetCurrentWeek("MLB")
		case "UFC":
			go GetCurrentMonth("UFC")
		case "Bellator":
			go GetCurrentMonth("Bellator")
		default:

		}

		owner.sports.feed.SetText("dReams Client")
		owner.sports.feed.Refresh()
	}
	owner.sports.league.SetPlaceHolder("League:")

	owner.sports.end = dwidget.NewDeroEntry("", 1, 0)
	owner.sports.end.SetPlaceHolder("Closes At:")
	owner.sports.end.Validator = validation.NewRegexp(`^\d{10,}$`, "Unix time required")

	owner.sports.amt = dwidget.NewDeroEntry("", 0.1, 1)
	owner.sports.amt.SetPlaceHolder("Minimum Amount:")
	owner.sports.amt.AllowFloat = true
	owner.sports.amt.Wrapping = fyne.TextTruncate
	owner.sports.amt.Validator = validation.NewRegexp(`^\d{1,}\.\d{1,5}$|^[^0]\d{0,}$`, "Int or float required")

	feeds := []string{"dReams Client"}
	owner.sports.feed = widget.NewSelectEntry(feeds)
	owner.sports.feed.SetPlaceHolder("Feed:")

	owner.sports.deposit = dwidget.NewDeroEntry("", 0.1, 1)
	owner.sports.deposit.SetPlaceHolder("Deposit Amount:")
	owner.sports.deposit.AllowFloat = true
	owner.sports.deposit.Wrapping = fyne.TextTruncate
	owner.sports.deposit.Validator = validation.NewRegexp(`^\d{1,}\.\d{1,5}$|^[^0]\d{0,}$`, "Int or float required")

	reset := window.Content().(*fyne.Container).Objects[2]

	owner.sports.set = widget.NewButton("Set Game", func() {
		if owner.sports.deposit.Validate() == nil && owner.sports.amt.Validate() == nil && owner.sports.end.Validate() == nil {
			if len(Sports.Contract.SCID) == 64 {
				window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(1, 100, window, reset))
				window.Content().(*fyne.Container).Objects[2].Refresh()
				return
			}

			dialog.NewInformation("Sports", "Select a valid contract", window).Show()
		}
	})

	owner.sports.cancel = widget.NewButton("Cancel", func() {
		window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(9, 0, window, reset))
		window.Content().(*fyne.Container).Objects[2].Refresh()
	})

	owner.sports.cancel.Hide()

	owner.sports.payout = widget.NewSelectEntry([]string{})
	owner.sports.payout.SetPlaceHolder("Game #")

	sports_confirm := widget.NewButton("Sports Payout", func() {
		if len(Sports.Contract.SCID) == 64 {
			window.Content().(*fyne.Container).Objects[2] = container.NewMax(ownerConfirmAction(3, 100, window, reset))
			window.Content().(*fyne.Container).Objects[2].Refresh()
			return
		}

		dialog.NewInformation("Sports", "Select a valid contract", window).Show()
	})

	sports := container.NewVBox(
		humanTimeConvert(),
		layout.NewSpacer(),
		owner.sports.league,
		owner.sports.game,
		owner.sports.end,
		owner.sports.amt,
		owner.sports.feed,
		owner.sports.deposit,
		owner.sports.set,
		layout.NewSpacer(),
		owner.sports.cancel,
		layout.NewSpacer(),
		owner.sports.payout,
		sports_confirm,
		layout.NewSpacer())

	return sports
}

// dService control objects for side menu
//   - Pass side menu window to reset to
func serviceOpts(window fyne.Window) fyne.CanvasObject {
	get_addr := widget.NewButton("Integrated Address", func() {
		go MakeIntegratedAddr(true)
	})

	txid := widget.NewMultiLineEntry()
	txid.SetPlaceHolder("TXID:")
	txid.Wrapping = fyne.TextWrapWord
	txid.Validator = validation.NewRegexp(`^\w{64,64}$`, "Invalid TXID")

	process := widget.NewButton("Process Tx", func() {
		if !Service.IsProcessing() && !Service.IsRunning() {
			if txid.Validate() == nil {
				processSingleTx(txid.Text)
			}
		} else {
			logger.Warnln("[dService] Stop service to manually process Tx")
		}
	})

	delete := widget.NewButton("Delete Tx", func() {
		if !Service.IsProcessing() && !Service.IsRunning() {
			if txid.Validate() == nil {
				e := rpc.GetWalletTx(txid.Text)
				if e != nil {
					if db := boltDB(); db != nil {
						defer db.Close()
						deleteTx("BET", db, *e)
					}
				}
			}
		} else {
			logger.Warnln("[dService] Stop service to delete Tx")
		}
	})

	store := widget.NewButton("Store Tx", func() {
		if !Service.IsProcessing() && !Service.IsRunning() {
			if txid.Validate() == nil {
				e := rpc.GetWalletTx(txid.Text)
				if e != nil {
					if db := boltDB(); db != nil {
						defer db.Close()
						storeTx("BET", "done", db, *e)
					}
				}
			}
		} else {
			logger.Warnln("[dService] Stop service to store Tx")
		}
	})

	entry := dwidget.NewDeroEntry("", 1, 0)
	entry.SetPlaceHolder("Block #:")
	entry.AllowFloat = false
	entry.Wrapping = fyne.TextTruncate
	entry.Validator = validation.NewRegexp(`^[^0]\d{0,}$`, "Int required")

	var start uint64
	height := widget.NewCheck("Start from current height", func(b bool) {
		if b {
			start = rpc.DaemonHeight("Prediction", rpc.Daemon.Rpc)
			entry.SetText(strconv.Itoa(int(start)))
			entry.Disable()
		} else {
			entry.SetText("")
			entry.Enable()
		}
	})
	height.SetChecked(true)

	debug := widget.NewCheck("Debug", func(b bool) {
		if b {
			Service.Debug = true
		} else {
			Service.Debug = false
		}
	})

	view := widget.NewButton("View Tx History", func() {
		if !Service.IsProcessing() && !Service.IsRunning() {
			if !height.Checked {
				start = uint64(rpc.StringToInt(entry.Text))
			}
			viewProcessedTx(start)
		} else {
			logger.Warnln("[dService] Stop service to view Tx history")
		}
	})

	owner.service.payouts.check = widget.NewCheck("Payouts", func(b bool) {
		if b {
			owner.service.payouts.enabled = true
		} else {
			owner.service.payouts.enabled = false
		}
	})

	if owner.service.payouts.enabled {
		owner.service.payouts.check.SetChecked(true)
		owner.service.payouts.check.Disable()
	}

	owner.service.transactions.check = widget.NewCheck("Transactions", func(b bool) {
		if b {
			owner.service.transactions.enabled = true
		} else {
			owner.service.transactions.enabled = false
		}
	})

	if owner.service.transactions.enabled {
		owner.service.transactions.check.SetChecked(true)
		owner.service.transactions.check.Disable()
	}

	reset := window.Content().(*fyne.Container).Objects[2]

	owner.service.run = widget.NewButton("Run Service", func() {
		if !Service.IsRunning() {
			if entry.Validate() == nil {
				if !height.Checked {
					start = uint64(rpc.StringToInt(entry.Text))
					if start < PAYLOAD_FORMAT {
						start = PAYLOAD_FORMAT
					}
				}

				if owner.service.payouts.check.Checked || owner.service.transactions.check.Checked {
					go func() {
						Service.Start()
						owner.service.run.Hide()
						window.Content().(*fyne.Container).Objects[2] = serviceRunConfirm(start, owner.service.payouts.check.Checked, owner.service.transactions.check.Checked, window, reset)
						window.Content().(*fyne.Container).Objects[2].Refresh()
					}()
				} else {
					logger.Warnln("[dService] Select which services to run")
				}
			} else {
				logger.Warnln("[dService] Enter service starting height")
			}
		} else {
			logger.Warnln("[dService] Service already running")
		}
	})

	if Service.IsRunning() || Service.IsProcessing() {
		owner.service.run.Hide()
	}

	stop := widget.NewButton("Stop Service", func() {
		if Service.IsRunning() {
			logger.Println("[dService] Stopping service")
		}
		Service.Stop()

	})

	box := container.NewVBox(
		layout.NewSpacer(),
		view,
		layout.NewSpacer(),
		txid,
		container.NewAdaptiveGrid(3,
			process,
			delete,
			store),
		layout.NewSpacer(),
		get_addr,
		layout.NewSpacer(),
		height,
		entry,
		owner.service.payouts.check,
		owner.service.transactions.check,
		debug,
		container.NewAdaptiveGrid(2,
			stop,
			owner.service.run,
		))

	return box
}

// SCID update objects for side menu
func updateOpts() fyne.CanvasObject {
	a_label := widget.NewLabel("Time A         ")
	a := dwidget.NewDeroEntry("", 1, 0)
	a.SetPlaceHolder("Time A:")
	a.AllowFloat = false
	a.Wrapping = fyne.TextTruncate
	a.Validator = validation.NewRegexp(`[^0]\d{1,}$`, "Int required")

	b_label := widget.NewLabel("Time B         ")
	b := dwidget.NewDeroEntry("", 1, 0)
	b.SetPlaceHolder("Time B:")
	b.AllowFloat = false
	b.Wrapping = fyne.TextTruncate
	b.Validator = validation.NewRegexp(`[^0]\d{1,}$`, "Int required")

	c_label := widget.NewLabel("Time C         ")
	c := dwidget.NewDeroEntry("", 1, 0)
	c.SetPlaceHolder("Time C:")
	c.AllowFloat = false
	c.Wrapping = fyne.TextTruncate
	c.Validator = validation.NewRegexp(`[^0]\d{1,}$`, "Int required")

	hl_label := widget.NewLabel("Max Games")
	hl := dwidget.NewDeroEntry("", 1, 0)
	hl.SetPlaceHolder("Max Games:")
	hl.AllowFloat = false
	hl.Wrapping = fyne.TextTruncate
	hl.Validator = validation.NewRegexp(`^[^0]\d{0,}$`, "Int required")

	hl_box := container.NewBorder(nil, nil, hl_label, nil, hl)
	hl_box.Hide()

	// l := dwidget.WholeAmtEntry()
	// l.PlaceHolder = "L:"
	// l.Validator = validation.NewRegexp(`^\d{2,}$`, "Format Not Valid")

	sc := widget.NewSelect([]string{"Prediction", "Sports"}, func(s string) {
		if s == "Sports" {
			c_label.SetText("Delete         ")
			c.Validator = validation.NewRegexp(`[^0]\d{0,}$`, "Int required")
			hl_box.Show()
		} else {
			c_label.SetText("Time C         ")
			c.Validator = validation.NewRegexp(`[^0]\d{1,}$`, "Int required")
			hl_box.Hide()
		}
	})
	sc.PlaceHolder = "Select Contract"

	new_owner := widget.NewMultiLineEntry()
	new_owner.Validator = validation.NewRegexp(`^(dero)\w{62}$`, "Invalid Address")
	new_owner.Wrapping = fyne.TextWrapWord
	new_owner.SetPlaceHolder("New owner address:")
	add_owner := widget.NewButton("Add Owner", func() {
		if new_owner.Validate() == nil {
			switch sc.Selected {
			case "Prediction":
				AddOwner(Predict.Contract.SCID, new_owner.Text)
			case "Sports":
				AddOwner(Sports.Contract.SCID, new_owner.Text)
			default:
				logger.Warnln("[dService] Select contract")
			}
		}
	})

	owner_num := dwidget.NewDeroEntry("", 1, 0)
	owner_num.SetPlaceHolder("Owner #:")
	owner_num.AllowFloat = false
	owner_num.Validator = validation.NewRegexp(`^[^0]\d{0,0}$`, "Int required")
	owner_num.Wrapping = fyne.TextTruncate

	remove_owner := widget.NewButton("Remove Owner", func() {
		switch sc.Selected {
		case "Prediction":
			RemoveOwner(Predict.Contract.SCID, rpc.StringToInt(owner_num.Text))
		case "Sports":
			RemoveOwner(Sports.Contract.SCID, rpc.StringToInt(owner_num.Text))
		default:
			logger.Warnln("[dService] Select contract")
		}
	})

	update_var := widget.NewButton("Update Variables", func() {
		if a.Validate() == nil && b.Validate() == nil && c.Validate() == nil {
			switch sc.Selected {
			case "Prediction":
				VarUpdate(Predict.Contract.SCID, rpc.StringToInt(a.Text), rpc.StringToInt(b.Text), rpc.StringToInt(c.Text), 30, 0)
			case "Sports":
				if hl.Validate() == nil {
					VarUpdate(Sports.Contract.SCID, rpc.StringToInt(a.Text), rpc.StringToInt(b.Text), rpc.StringToInt(c.Text), 30, rpc.StringToInt(hl.Text))
				}
			default:
				logger.Warnln("[dService] Select contract")
			}
		}
	})

	return container.NewVBox(
		sc,
		container.NewBorder(nil, nil, a_label, nil, a),
		container.NewBorder(nil, nil, b_label, nil, b),
		container.NewBorder(nil, nil, c_label, nil, c),
		hl_box,
		update_var,
		layout.NewSpacer(),
		new_owner,
		add_owner,
		layout.NewSpacer(),
		container.NewBorder(nil, nil, nil, remove_owner, owner_num),
		layout.NewSpacer())

}

// dSports and dPrediction action confirmation
//   - i defines the action to be confirmed
//   - teamA, teamB needed only for dSports confirmations
//   - Pass main window obj and tabs to reset to
func ConfirmAction(i int, teamA, teamB string, obj []fyne.CanvasObject, tabs *container.AppTabs) fyne.CanvasObject {
	var confirm_display = widget.NewLabel("")
	confirm_display.Wrapping = fyne.TextWrapWord
	confirm_display.Alignment = fyne.TextAlignCenter

	p_scid := Predict.Contract.SCID

	s_scid := Sports.Contract.SCID
	split := strings.Split(Sports.gameSelect.Selected, "   ")
	multi := Sports.multi.Selected

	switch i {
	case 1:
		float := float64(Predict.amount)
		amt := float / 100000
		confirm_display.SetText(fmt.Sprintf("SCID:\n\n%s\n\nLower prediction for %.5f Dero", p_scid, amt))
	case 2:
		float := float64(Predict.amount)
		amt := float / 100000
		confirm_display.SetText(fmt.Sprintf("SCID:\n\n%s\n\nHigher prediction for %.5f Dero", p_scid, amt))
	case 3:
		game := Sports.gameSelect.Selected
		val := float64(GetSportsAmt(s_scid, split[0]))
		var x string

		switch multi {
		case "3x":
			x = fmt.Sprintf("%.5f", val*3/100000)
		case "5x":
			x = fmt.Sprintf("%.5f", val*5/100000)
		default:
			x = fmt.Sprintf("%.5f", val/100000)
		}

		confirm_display.SetText(fmt.Sprintf("SCID:\n\n%s\n\nBetting on Game # %s\n\n%s for %s Dero", s_scid, game, teamA, x))
	case 4:
		game := Sports.gameSelect.Selected
		val := float64(GetSportsAmt(s_scid, split[0]))
		var x string

		switch multi {
		case "3x":
			x = fmt.Sprintf("%.5f", val*3/100000)
		case "5x":
			x = fmt.Sprintf("%.5f", val*5/100000)
		default:
			x = fmt.Sprintf("%.5f", val/100000)
		}

		confirm_display.SetText(fmt.Sprintf("SCID:\n\n%s\n\nBetting on Game # %s\n\n%s for %s Dero", s_scid, game, teamB, x))
	default:
		logger.Errorln("[dService] No Confirm Input")
		confirm_display.SetText("Error")
	}

	cancel_button := widget.NewButton("Cancel", func() {
		obj[1] = tabs
		obj[1].Refresh()
	})

	confirm_button := widget.NewButton("Confirm", func() {
		switch i {
		case 1:
			PredictLower(p_scid, "")
		case 2:
			PredictHigher(p_scid, "")
		case 3:
			PickTeam(s_scid, multi, split[0], GetSportsAmt(s_scid, split[0]), 0)
		case 4:
			PickTeam(s_scid, multi, split[0], GetSportsAmt(s_scid, split[0]), 1)
		default:

		}

		obj[1] = tabs
		obj[1].Refresh()
	})

	display := container.NewVBox(layout.NewSpacer(), confirm_display, layout.NewSpacer())
	options := container.NewAdaptiveGrid(2, confirm_button, cancel_button)
	content := container.NewBorder(nil, options, nil, nil, display)

	go func() {
		for rpc.IsReady() {
			time.Sleep(time.Second)
		}

		obj[1] = tabs
		obj[1].Refresh()
	}()

	return container.NewMax(bundle.Alpha120, content)
}

// dReam Service start confirmation
//   - start is starting height to run service
//   - payout and transfers, params for service
//   - Pass side window to reset to
func serviceRunConfirm(start uint64, payout, transfers bool, window fyne.Window, reset fyne.CanvasObject) fyne.CanvasObject {
	var pay, transac string
	if transfers {
		transac = "process transactions sent to your integrated address"
		if payout {
			transac = transac + " "
		}
	}

	if payout {
		if transfers {
			pay = "and "
		}
		pay = pay + "process payouts to contracts"
	}

	str := fmt.Sprintf("This will automatically %s%s.\n\nStarting service from height %d", transac, pay, start)
	confirm_display := widget.NewLabel(str)
	confirm_display.Wrapping = fyne.TextWrapWord
	confirm_display.Alignment = fyne.TextAlignCenter

	cancel_button := widget.NewButton("Cancel", func() {
		Service.Stop()
		window.Content().(*fyne.Container).Objects[2] = reset
		window.Content().(*fyne.Container).Objects[2].Refresh()
	})

	confirm_button := widget.NewButton("Confirm", func() {
		go RunService(start, payout, transfers)
		window.Content().(*fyne.Container).Objects[2] = reset
		window.Content().(*fyne.Container).Objects[2].Refresh()
	})

	display := container.NewVBox(layout.NewSpacer(), confirm_display, layout.NewSpacer())
	options := container.NewAdaptiveGrid(2, confirm_button, cancel_button)
	content := container.NewBorder(nil, options, nil, nil, display)

	return container.NewMax(content)
}

// Convert unix time to human readable time
func humanTimeConvert() fyne.CanvasObject {
	entry := dwidget.NewDeroEntry("", 1, 0)
	entry.AllowFloat = false
	entry.SetPlaceHolder("Unix time:")
	entry.Validator = validation.NewRegexp(`^\d{10,}$`, "Unix time required")
	res := widget.NewEntry()
	res.Disable()
	button := widget.NewButton("Human Time", func() {
		if entry.Validate() == nil {
			time := time.Unix(int64(rpc.StringToInt(entry.Text)), 0).String()
			res.SetText(time)
		}
	})

	split := container.NewHSplit(entry, button)
	box := container.NewVBox(res, split)

	return box
}

// Check dPrediction SCID for live status
func CheckPredictionStatus() {
	if rpc.Daemon.IsConnected() && menu.Gnomes.IsReady() {
		_, ends := menu.Gnomes.GetSCIDValuesByKey(Predict.Contract.SCID, "p_end_at")
		_, time_a := menu.Gnomes.GetSCIDValuesByKey(Predict.Contract.SCID, "time_a")
		_, time_c := menu.Gnomes.GetSCIDValuesByKey(Predict.Contract.SCID, "time_c")
		_, mark := menu.Gnomes.GetSCIDValuesByKey(Predict.Contract.SCID, "mark")
		if ends != nil && time_a != nil && time_c != nil {
			now := uint64(time.Now().Unix())
			if now >= ends[0] && now <= ends[0]+time_a[0] && mark == nil {
				owner.predict.post.Show()
			} else {
				owner.predict.post.Hide()
			}

			if now >= ends[0]+time_c[0] {
				owner.predict.pay.Show()
			} else {
				owner.predict.pay.Hide()
			}
		}

		if ends == nil {
			owner.predict.post.Hide()
			owner.predict.pay.Hide()
		}
	}
}

// Check dSports SCID for active games
func GetActiveGames() {
	if rpc.Daemon.IsConnected() && menu.Gnomes.IsReady() {
		options := []string{}
		contracts := menu.Gnomes.GetAllOwnersAndSCIDs()
		for sc := range contracts {
			owner, _ := menu.Gnomes.GetSCIDValuesByKey(sc, "owner")
			if (owner != nil && owner[0] == rpc.Wallet.Address) || VerifyBetSigner(sc) {
				if len(sc) == 64 {
					_, init := menu.Gnomes.GetSCIDValuesByKey(sc, "s_init")
					if init != nil {
						for ic := uint64(1); ic <= init[0]; ic++ {
							num := strconv.Itoa(int(ic))
							if game, _ := menu.Gnomes.GetSCIDValuesByKey(sc, "game_"+num); game != nil {
								league, _ := menu.Gnomes.GetSCIDValuesByKey(sc, "league_"+num)
								_, end := menu.Gnomes.GetSCIDValuesByKey(sc, "s_end_at_"+num)
								_, add := menu.Gnomes.GetSCIDValuesByKey(sc, "time_a")
								if league != nil && end != nil && add != nil {
									if end[0]+add[0] < uint64(time.Now().Unix()) {
										options = append(options, num+"   "+league[0]+"   "+game[0])
									}
								}
							}
						}
					}
				}
			}
		}
		owner.sports.payout.SetOptions(options)
	}
}

// Bet contract owner control menu
func ownersMenu() {
	ow := fyne.CurrentApp().NewWindow("Bet Contracts")
	ow.Resize(fyne.NewSize(330, 700))
	ow.SetIcon(bundle.ResourceDReamsIconAltPng)
	Predict.Contract.menu.Hide()
	Sports.Contract.menu.Hide()
	quit := make(chan struct{})
	ow.SetCloseIntercept(func() {
		Predict.Contract.menu.Show()
		Sports.Contract.menu.Show()
		quit <- struct{}{}
		ow.Close()
	})
	ow.SetFixedSize(true)

	owner_tabs := container.NewAppTabs(
		container.NewTabItem("Predict", layout.NewSpacer()),
		container.NewTabItem("Sports", layout.NewSpacer()),
		container.NewTabItem("Service", layout.NewSpacer()),
		container.NewTabItem("Update", updateOpts()),
	)
	owner_tabs.SetTabLocation(container.TabLocationTop)
	owner_tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Sports":
			go GetActiveGames()
		case "Service":
			go MakeIntegratedAddr(false)
		}
	}

	var utime string
	clock := widget.NewEntry()
	clock.Disable()

	entry := dwidget.NewDeroEntry("", 1, 0)
	entry.AllowFloat = false
	entry.SetPlaceHolder("Hours to close:")
	entry.Validator = validation.NewRegexp(`^\d{1,}$`, "Int required")
	button := widget.NewButton("Add Hours", func() {
		if entry.Validate() == nil {
			i := rpc.StringToInt(entry.Text)
			u := rpc.StringToInt(utime)
			r := u + (i * 3600)

			switch owner_tabs.SelectedIndex() {
			case 0:
				owner.predict.end.SetText(strconv.Itoa(r))
			case 1:
				owner.sports.end.SetText(strconv.Itoa(r))
			}
		}
	})

	go func() {
		var ticker = time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				if !rpc.Wallet.IsConnected() {
					ticker.Stop()
					ow.Close()
				}

				if Service.IsRunning() {
					owner.service.run.Hide()
					owner.service.payouts.check.Disable()
					owner.service.transactions.check.Disable()
				}

				if !Service.IsRunning() && !Service.IsProcessing() {
					owner.service.run.Show()
					owner.service.payouts.check.Enable()
					owner.service.transactions.check.Enable()
				}

				CheckPredictionStatus()
				now := time.Now()
				utime = strconv.Itoa(int(now.Unix()))
				clock.SetText("Unix Time: " + utime)
				if now.Unix() < Predict.buffer {
					if Predict.init {
						owner.predict.set.Hide()
						owner.predict.cancel.Show()
					} else {
						owner.predict.set.Show()
						owner.predict.cancel.Hide()
					}
				} else {
					owner.predict.cancel.Hide()
					if Predict.init {
						owner.predict.set.Hide()
					} else {
						owner.predict.set.Show()
					}
				}

				if Sports.buffer {
					owner.sports.cancel.Show()
					owner.sports.set.Hide()
				} else {
					owner.sports.cancel.Hide()
					owner.sports.set.Show()
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	bottom_split := container.NewHSplit(entry, button)
	bottom_box := container.NewVBox(clock, bottom_split)

	border := container.NewBorder(nil, bottom_box, nil, nil, owner_tabs)

	alpha := canvas.NewRectangle(color.RGBA{0, 0, 0, 180})
	if bundle.AppColor == color.White {
		alpha = canvas.NewRectangle(color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x99})
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		ow.SetContent(
			container.New(
				layout.NewMaxLayout(),
				menu.BackgroundRast("ownersMenu"),
				alpha,
				border))

		owner_tabs.SelectIndex(2)
		owner_tabs.Selected().Content = serviceOpts(ow)
		owner_tabs.SelectIndex(1)
		owner_tabs.Selected().Content = sportsOpts(ow)
		owner_tabs.SelectIndex(0)
		owner_tabs.Selected().Content = predictionOpts(ow)

		time.Sleep(time.Second)
		markets := []string{}
		if stored, ok := rpc.FindStringKey(rpc.RatingSCID, "prediction_markets", rpc.Daemon.Rpc).(string); ok {
			if h, err := hex.DecodeString(stored); err == nil {
				if err = json.Unmarshal(h, &markets); err == nil {
					owner.predict.name.SetOptions(markets)
				}
			}
		}

		leagues := []string{}
		if stored, ok := rpc.FindStringKey(rpc.RatingSCID, "sports_leagues", rpc.Daemon.Rpc).(string); ok {
			if h, err := hex.DecodeString(stored); err == nil {
				if err = json.Unmarshal(h, &leagues); err == nil {
					owner.sports.league.SetOptions(leagues)
				}
			}
		}
	}()

	ow.Show()
}

// Bet contract owner action confirmation
//   - i defines action to be confirmed
//   - p for prediction price
//   - Pass side window to reset to
func ownerConfirmAction(i int, p float64, window fyne.Window, reset fyne.CanvasObject) fyne.CanvasObject {
	var confirm_display = widget.NewLabel("")
	confirm_display.Wrapping = fyne.TextWrapWord
	confirm_display.Alignment = fyne.TextAlignCenter

	pre := Predict.prediction
	p_scid := Predict.Contract.SCID
	p_pre := owner.predict.name.Text
	p_amt := owner.predict.amt.Text
	if p_amt_f, err := strconv.ParseFloat(p_amt, 64); err == nil {
		p_amt = fmt.Sprintf("%.5f", p_amt_f)
	}
	p_mark := owner.predict.mark.Text
	p_end := owner.predict.end.Text
	p_end_time, _ := rpc.MsToTime(p_end + "000")
	p_feed := owner.predict.feed.Text
	p_dep := owner.predict.deposit.Text
	if p_dep_f, err := strconv.ParseFloat(p_dep, 64); err == nil {
		p_dep = fmt.Sprintf("%.5f", p_dep_f)
	}
	var price string
	if menu.CoinDecimal(pre) == 8 {
		price = fmt.Sprintf("%.8f", p/100000000)
	} else {
		price = fmt.Sprintf("%.2f", p/100)
	}

	var s_game string
	s_scid := Sports.Contract.SCID
	game_split := strings.Split(owner.sports.game.Selected, "   ")
	if len(game_split) > 1 {
		s_game = game_split[1]
	} else {
		s_game = game_split[0]
	}

	s_league := owner.sports.league.Text
	s_amt := owner.sports.amt.Text
	if s_amt_f, err := strconv.ParseFloat(s_amt, 64); err == nil {
		s_amt = fmt.Sprintf("%.5f", s_amt_f)
	}
	s_end := owner.sports.end.Text
	s_end_time, _ := rpc.MsToTime(s_end + "000")
	s_feed := owner.sports.feed.Text
	n_split := strings.Split(owner.sports.payout.Text, "   ")
	s_pay_n := n_split[0]
	s_dep := owner.sports.deposit.Text
	if s_dep_f, err := strconv.ParseFloat(s_dep, 64); err == nil {
		s_dep = fmt.Sprintf("%.5f", s_dep_f)
	}

	var win, team, a_score, b_score, payout_str, err_string string
	if i == 3 {
		if len(n_split) > 1 {
			if n_split[1] == "Bellator" || n_split[1] == "UFC" {
				win, team = GetMmaWinner(n_split[2], n_split[1])
				payout_str = fmt.Sprintf("SCID:\n\n%s\n\nFight: %s\n\nWinner: %s", s_scid, owner.sports.payout.Text, team)
			} else {
				win, team, a_score, b_score = GetWinner(n_split[2], n_split[1])
				payout_str = fmt.Sprintf("SCID:\n\n%s\n\nGame: %s\n\n%s: %s\n%s: %s\n\nWinner: %s", s_scid, owner.sports.payout.Text, TrimTeamA(n_split[2]), a_score, TrimTeamB(n_split[2]), b_score, team)
			}
		} else {
			logger.Errorln("[dService] Could not format game string")
			i = 100
			if owner.sports.payout.Text == "" {
				err_string = "Select game for payout"
			}
		}
	}

	switch i {
	case 1:
		confirm_display.SetText("SCID:\n\n" + s_scid + "\n\nGame: " + s_game + "\n\nMinimum: " + s_amt + " Dero\n\nCloses At: " + s_end_time.String() + "\n\nFeed: " + s_feed + "\n\nInitial Deposit: " + s_dep + " Dero")
	case 2:
		fn := "Feed: "
		var mark string
		if p_mark == "0" || p_mark == "" {
			mark = "Not Set"
		} else {
			if onChainPrediction(pre) == 2 || onChainPrediction(p_pre) == 2 { /// one decimal place for block time
				fn = "Node: "
				i := rpc.StringToInt(p_mark) * 10000
				x := float64(i) / 100000
				mark = fmt.Sprintf("%.5f", x)
			} else {
				if isOnChainPrediction(pre) || isOnChainPrediction(p_pre) {
					fn = "Node: "
					mark = p_mark
				} else {
					if menu.CoinDecimal(pre) == 8 || menu.CoinDecimal(p_pre) == 8 {
						if f, err := strconv.ParseFloat(p_mark, 32); err == nil { /// eight decimal place for btc
							x := f / 100000000
							mark = fmt.Sprintf("%.8f", x)
						}
					} else {
						if f, err := strconv.ParseFloat(p_mark, 32); err == nil {
							x := f / 100
							mark = fmt.Sprintf("%.2f", x)
						}
					}
				}
			}
		}

		confirm_display.SetText("SCID:\n\n" + p_scid + "\n\nPredicting: " + p_pre + "\n\nMinimum: " + p_amt + " Dero\n\nCloses At: " + p_end_time.String() + "\n\nMark: " + mark + "\n\n" + fn + p_feed + "\n\nInitial Deposit: " + p_dep + " Dero")

	case 3:
		confirm_display.SetText(payout_str)
	case 4:
		confirm_display.SetText("SCID:\n\n" + p_scid + "\n\nFeed from: dReams Client\n\nPost Price: " + price)
	case 5:
		confirm_display.SetText("SCID:\n\n" + p_scid + "\n\nFeed from: dReams Client\n\nFinal Price: " + price)
	case 6:
		switch onChainPrediction(pre) {
		case 1:
			confirm_display.SetText("SCID:\n\n" + p_scid + "\n\n" + pre + ": " + fmt.Sprintf("%.0f", p) + "\n\nNode: " + Predict.feed + "\n\nConfirm Post")
		case 2:
			confirm_display.SetText("SCID:\n\n" + p_scid + "\n\n" + pre + ": " + fmt.Sprintf("%.5f", p) + "\n\nNode: " + Predict.feed + "\n\nConfirm Post")
		case 3:
			confirm_display.SetText("SCID:\n\n" + p_scid + "\n\n" + pre + ": " + fmt.Sprintf("%.0f", p) + "\n\nNode: " + Predict.feed + "\n\nConfirm Post")
		}

	case 7:
		switch onChainPrediction(pre) {
		case 1:
			confirm_display.SetText("SCID:\n\n" + p_scid + "\n\n" + pre + ": " + fmt.Sprintf("%.0f", p) + "\n\nNode: " + Predict.feed + "\n\nConfirm Payout")
		case 2:
			confirm_display.SetText("SCID:\n\n" + p_scid + "\n\n" + pre + ": " + fmt.Sprintf("%.5f", p) + "\n\nNode: " + Predict.feed + "\n\nConfirm Payout")
		case 3:
			confirm_display.SetText("SCID:\n\n" + p_scid + "\n\n" + pre + ": " + fmt.Sprintf("%.0f", p) + "\n\nNode: " + Predict.feed + "\n\nConfirm Payout")
		}

	case 8:
		confirm_display.SetText("SCID:\n\n" + p_scid + "\n\nThis will Cancel the current prediction")
	case 9:
		confirm_display.SetText("SCID:\n\n" + s_scid + "\n\nThis will Cancel the last initiated bet on this contract")
	default:
		logger.Errorln("[dService] No Confirm Input")
		confirm_display.SetText("Error\n\n" + err_string)
	}

	cancel_button := widget.NewButton("Cancel", func() {
		window.Content().(*fyne.Container).Objects[2] = reset
		window.Content().(*fyne.Container).Objects[2].Refresh()
	})

	confirm_button := widget.NewButton("Confirm", func() {
		owner.sports.payout.SetText("")
		switch i {
		case 1:
			SetSports(rpc.StringToInt(s_end), rpc.ToAtomic(s_amt, 5), rpc.ToAtomic(s_dep, 5), s_scid, s_league, s_game, s_feed)
		case 2:
			if onChainPrediction(pre) == 2 || onChainPrediction(p_pre) == 2 { /// decimal of one place for block time
				SetPrediction(rpc.StringToInt(p_end), rpc.StringToInt(p_mark)*10000, rpc.ToAtomic(p_amt, 5), rpc.ToAtomic(p_dep, 5), p_scid, p_pre, p_feed)
			} else {
				SetPrediction(rpc.StringToInt(p_end), rpc.StringToInt(p_mark), rpc.ToAtomic(p_amt, 5), rpc.ToAtomic(p_dep, 5), p_scid, p_pre, p_feed)
			}
		case 3:
			EndSports(s_scid, s_pay_n, win)
		case 4:
			PostPrediction(p_scid, int(p))
		case 5:
			EndPrediction(p_scid, int(p))
		case 6:
			switch onChainPrediction(pre) {
			case 1:
				PostPrediction(p_scid, int(p))
			case 2:
				PostPrediction(p_scid, int(p*100000))
			case 3:
				PostPrediction(p_scid, int(p))
			default:
			}
		case 7:
			switch onChainPrediction(pre) {
			case 1:
				EndPrediction(p_scid, int(p))
			case 2:
				EndPrediction(p_scid, int(p*100000))
			case 3:
				EndPrediction(p_scid, int(p))
			default:
			}
		case 8:
			CancelInitiatedBet(Predict.Contract.SCID, 0)
		case 9:
			CancelInitiatedBet(Sports.Contract.SCID, 1)
		default:

		}

		window.Content().(*fyne.Container).Objects[2] = reset
		window.Content().(*fyne.Container).Objects[2].Refresh()
	})

	display := container.NewVBox(layout.NewSpacer(), confirm_display, layout.NewSpacer())
	options := container.NewAdaptiveGrid(2, confirm_button, cancel_button)
	content := container.NewBorder(nil, options, nil, nil, display)

	return container.NewMax(content)
}

// Confirmation for dPrediction contract installs
func newPredictConfirm(c int, obj []fyne.CanvasObject, tabs *container.AppTabs) fyne.CanvasObject {
	var text string
	gas_fee := 0.125
	unlock_fee := float64(rpc.UnlockFee) / 100000
	switch c {
	case 1:
		text = `You are about to unlock and install your first dPrediction contract 
		
To help support the project, there is a ` + fmt.Sprintf("%.5f", unlock_fee) + ` DERO donation attached to preform this action

Unlocking dPrediction or dSports gives you unlimited access to bet contract uploads and all base level owner features

Total transaction will be ` + fmt.Sprintf("%0.5f", unlock_fee+gas_fee) + ` DERO (0.12500 gas fee for contract install)


Select a public or private contract

Public will show up in indexed list of contracts

Private will not show up in the list`
	case 2:
		text = `You are about to install a new dPrediction contract

Gas fee to install contract is 0.12500 DERO


Select a public or private contract

Public will show up in indexed list of contracts

Private will not show up in the list`
	}

	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	label.Alignment = fyne.TextAlignCenter

	var choice *widget.Select

	pre_button := widget.NewButton("Install", func() {
		if choice.SelectedIndex() < 2 && choice.SelectedIndex() >= 0 {
			UploadBetContract(true, choice.SelectedIndex())
		}

		obj[1] = tabs
		obj[1].Refresh()
	})

	pre_button.Hide()

	options := []string{"Public", "Private"}
	choice = widget.NewSelect(options, func(s string) {
		if s == "Public" || s == "Private" {
			pre_button.Show()
		} else {
			pre_button.Hide()
		}
	})

	cancel_button := widget.NewButton("Cancel", func() {
		obj[1] = tabs
		obj[1].Refresh()
	})

	left := container.NewVBox(pre_button)
	right := container.NewVBox(cancel_button)
	buttons := container.NewAdaptiveGrid(3, left, container.NewVBox(layout.NewSpacer()), right)
	actions := container.NewVBox(choice, buttons)
	info_box := container.NewVBox(layout.NewSpacer(), label, layout.NewSpacer())

	content := container.NewBorder(nil, actions, nil, nil, info_box)

	go func() {
		for rpc.IsReady() {
			time.Sleep(time.Second)
		}

		obj[1] = tabs
		obj[1].Refresh()
	}()

	return container.NewMax(content)
}

// Confirmation for dSports contract installs
func newSportsConfirm(c int, obj []fyne.CanvasObject, tabs *container.AppTabs) fyne.CanvasObject {
	var text string
	gas_fee := 0.14
	unlock_fee := float64(rpc.UnlockFee) / 100000
	switch c {
	case 1:
		text = `You are about to unlock and install your first dSports contract
		
To help support the project, there is a ` + fmt.Sprintf("%.5f", unlock_fee) + ` DERO donation attached to preform this action

Unlocking dPrediction or dSports gives you unlimited access to bet contract uploads and all base level owner features

Total transaction will be ` + fmt.Sprintf("%0.5f", unlock_fee+gas_fee) + ` DERO (0.14000 gas fee for contract install)


Select a public or private contract

Public will show up in indexed list of contracts

Private will not show up in the list`
	case 2:
		text = `You are about to install a new dSports contract

Gas fee to install contract is 0.14000 DERO


Select a public or private contract

Public will show up in indexed list of contracts

Private will not show up in the list`
	}

	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	label.Alignment = fyne.TextAlignCenter

	var choice *widget.Select

	sports_button := widget.NewButton("Install", func() {
		if choice.SelectedIndex() < 2 && choice.SelectedIndex() >= 0 {
			UploadBetContract(false, choice.SelectedIndex())
		}

		obj[1] = tabs
		obj[1].Refresh()
	})

	sports_button.Hide()

	options := []string{"Public", "Private"}
	choice = widget.NewSelect(options, func(s string) {
		if s == "Public" || s == "Private" {
			sports_button.Show()
		} else {
			sports_button.Hide()
		}
	})

	cancel_button := widget.NewButton("Cancel", func() {
		obj[1] = tabs
		obj[1].Refresh()
	})

	left := container.NewVBox(sports_button)
	right := container.NewVBox(cancel_button)
	buttons := container.NewAdaptiveGrid(3, left, container.NewVBox(layout.NewSpacer()), right)
	actions := container.NewVBox(choice, buttons)
	info_box := container.NewVBox(layout.NewSpacer(), label, layout.NewSpacer())

	content := container.NewBorder(nil, actions, nil, nil, info_box)

	go func() {
		for rpc.IsReady() {
			time.Sleep(time.Second)
		}

		obj[1] = tabs
		obj[1].Refresh()
	}()

	return container.NewMax(content)
}
