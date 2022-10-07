package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var skipCount int
var customDaemonInput = widget.NewEntry()

func customPopUp() { /// pop up for entering custom daemon address
	cw := fyne.CurrentApp().NewWindow("Enter Custom Address")
	cw.SetIcon(resourceDReamTablesIconPng)
	cw.Resize(fyne.NewSize(400, 70))
	cw.SetFixedSize(true)
	content := container.NewWithoutLayout(customDaemonButton(), customeDaemonEdit())
	cw.SetContent(content)
	cw.Show()

}

func customeDaemonEdit() fyne.Widget { /// entry for custom daemon address
	customDaemonInput.SetPlaceHolder("Daemon Address")
	customDaemonInput.Resize(fyne.NewSize(270, 45))
	customDaemonInput.Move(fyne.NewPos(10, 10))

	return customDaemonInput
}

func customDaemonButton() fyne.Widget { /// confirm custom daemon button
	confirmButton := widget.NewButton("Enter", func() {
		log.Println("confirm tapped")
		suff := "/json_rpc"
		pre := "http://"
		daemonAddress = pre + customDaemonInput.Text + suff
		if gnomonInit {
			defaultIndexer.Close()
			gnomonInit = false
			gnomon(customDaemonInput.Text)
		}
		log.Println("Daemon Set To: CUSTOM")
	})
	confirmButton.Resize(fyne.NewSize(100, 42))
	confirmButton.Move(fyne.NewPos(290, 11))

	return confirmButton
}

func ownersMenu() { /// owners menu pop up
	ow := fyne.CurrentApp().NewWindow("Owners Controls")
	ow.Resize(fyne.NewSize(MIN_WIDTH-50, MIN_HEIGHT-100))
	ow.SetIcon(resourceDReamTablesIconPng)
	ow.SetFixedSize(true)

	owner_s := container.NewVBox(layout.NewSpacer(), s_gameEdit(), s_leagueEdit(), s_EndEdit(), s_AmtEdit(), s_feedEdit(), s_confirmButton())
	owner_p := container.NewVBox(layout.NewSpacer(), p_nameEdit(), p_EndEdit(), p_AmtEdit(), p_feedEdit(), p_confirmButton())
	payout := container.NewVBox(layout.NewSpacer(), nEdit(), payout_teamInput, payoutSButton(), layout.NewSpacer(), p_postPriceEdit(), postButton(), layout.NewSpacer(), p_finalPriceEdit(), payoutPButton())

	owner_tabs := container.NewAppTabs(
		container.NewTabItem("Sports", owner_s),
		container.NewTabItem("Predict", owner_p),
		container.NewTabItem("Payout", payout),
	)
	owner_tabs.SetTabLocation(container.TabLocationTop)
	ow.SetContent(owner_tabs)
	ow.Show()

}

// / owners menu widgets
var s_endInput = widget.NewEntry()
var s_amtInput = widget.NewEntry()
var s_gameInput = widget.NewEntry()
var s_leagueInput = widget.NewEntry()
var s_feedInput = widget.NewEntry()

var p_endInput = widget.NewEntry()
var p_amtInput = widget.NewEntry()
var p_nameInput = widget.NewEntry()
var p_feedInput = widget.NewEntry()

var post_priceInput = widget.NewEntry()

var payout_priceInput = widget.NewEntry()
var payout_nInput = widget.NewEntry()

func s_EndEdit() fyne.Widget { /// sports end entry
	s_endInput.SetPlaceHolder("Ends At:")
	s_endInput.OnChanged = func(input string) {

	}

	return s_endInput
}

func s_AmtEdit() fyne.Widget { /// sports amount entry
	s_amtInput.SetPlaceHolder("Minimum Amount:")
	s_amtInput.OnChanged = func(input string) {

	}

	return s_amtInput
}

func s_gameEdit() fyne.Widget { /// game entry
	s_gameInput.SetPlaceHolder("Game:")

	return s_gameInput
}

func s_leagueEdit() fyne.Widget { /// league entry
	s_leagueInput.SetPlaceHolder("League:")

	return s_leagueInput
}

func s_feedEdit() fyne.Widget { /// sports feed entry
	s_feedInput.SetPlaceHolder("Feed:")
	s_feedInput.OnChanged = func(input string) {
		/// TODO validate url
	}

	return s_feedInput
}

func s_confirmButton() fyne.Widget { /// confirm sports button
	confirmButton := widget.NewButton("Set Game", func() {
		ownerConfirmPopUp(1)
	})
	confirmButton.Resize(fyne.NewSize(100, 42))
	confirmButton.Move(fyne.NewPos(290, 11))

	return confirmButton
}

func p_EndEdit() fyne.Widget { /// prediction end entry
	p_endInput.SetPlaceHolder("Ends At:")
	p_endInput.OnChanged = func(input string) {
		/// TODO validate min entry
	}

	return p_endInput
}

func p_AmtEdit() fyne.Widget { /// prediction amount entry
	p_amtInput.SetPlaceHolder("Minimum Amount:")

	return p_amtInput
}

func p_nameEdit() fyne.Widget { /// prediction name entry
	p_nameInput.SetPlaceHolder("Name:")
	p_nameInput.OnChanged = func(input string) {
		/// TODO validate ticker
	}

	return p_nameInput
}

func p_feedEdit() fyne.Widget { /// prediction feed entry
	p_feedInput.SetPlaceHolder("Feed:")
	p_feedInput.OnChanged = func(input string) {
		/// TODO validate url
	}

	return p_feedInput
}

func p_confirmButton() fyne.Widget { /// prediction confirm button
	confirmButton := widget.NewButton("Set Prediction", func() {
		ownerConfirmPopUp(2)
	})

	return confirmButton
}

func p_postPriceEdit() fyne.Widget { /// prediction post entry
	post_priceInput.SetPlaceHolder("Post Price:")
	post_priceInput.OnChanged = func(input string) {
		/// TODO validate numerical
	}

	return post_priceInput
}

func postButton() fyne.Widget { /// prediction post button
	postButton := widget.NewButton("Post Price", func() {
		ownerConfirmPopUp(4)

	})

	return postButton
}

func p_finalPriceEdit() fyne.Widget { /// prediction price entry
	payout_priceInput.SetPlaceHolder("Final Price:")
	payout_priceInput.OnChanged = func(input string) {
		/// TODO validate numerical
	}

	return payout_priceInput
}

func nEdit() fyne.Widget { /// sports game number entry
	payout_nInput.SetPlaceHolder("Game#:")
	payout_nInput.OnChanged = func(input string) {

	}

	return payout_nInput
}

var payout_teamOptions = []string{"Team A", "Team B"}
var payout_teamInput = widget.NewSelect(payout_teamOptions, func(s string) {

})

func payoutSButton() fyne.Widget { /// sports payout button
	confirmButton := widget.NewButton("Sports Payout", func() {
		ownerConfirmPopUp(3)
	})

	return confirmButton
}

func payoutPButton() fyne.Widget { /// prediction payout button
	confirmButton := widget.NewButton("Prediction Payout", func() {
		ownerConfirmPopUp(5)
	})

	return confirmButton
}

func confirmPopUp(i int) { /// action confirmation pop up
	ocw := fyne.CurrentApp().NewWindow("Confirm")
	ocw.SetIcon(resourceDReamTablesIconPng)
	ocw.Resize(fyne.NewSize(MIN_WIDTH-50, 150))
	ocw.SetFixedSize(true)
	var confirm_display = widget.NewLabel("")
	confirm_display.Wrapping = fyne.TextWrapWord

	switch i {
	case 1:
		float := float64(amount)
		amt := float / 100000
		a := fmt.Sprintf("%f", amt)

		confirm_display.SetText("Lower prediction of " + a + " Dero\n\nConfirm")
	case 2:
		float := float64(amount)
		amt := float / 100000
		a := fmt.Sprintf("%f", amt)

		confirm_display.SetText("Higer prediction of " + a + " Dero\n\nConfirm")
	case 3:
		game := game_select.Text
		val := float64(getSportsAmt(game))
		var x string

		switch s_multi.Selected {
		case "3x":
			x = fmt.Sprint(val * 3 / 100000)
		case "5x":
			x = fmt.Sprint(val * 5 / 100000)
		default:
			x = fmt.Sprint(val / 100000)
		}
		confirm_display.SetText("Betting on Game#: " + game + "\nTeam " + "A" + " for " + x + " Dero\n\nConfirm")
	case 4:
		game := game_select.Text
		val := float64(getSportsAmt(game))
		var x string

		switch s_multi.Selected {
		case "3x":
			x = fmt.Sprint(val * 3 / 100000)
		case "5x":
			x = fmt.Sprint(val * 5 / 100000)
		default:
			x = fmt.Sprint(val / 100000)
		}
		confirm_display.SetText("Betting on Game#: " + game + "\nTeam " + "B" + " for " + x + " Dero\n\nConfirm")
	default:
		log.Println("No Confirm Input")
		confirm_display.SetText("Error")
	}

	n_confirmButton := widget.NewButton("No", func() {
		ocw.Close()
	})
	y_confirmButton := widget.NewButton("Yes", func() {
		switch i {
		case 1:
			predictLower()
		case 2:
			predictHigher()
		case 3:
			pickTeam(0)
		case 4:
			pickTeam(1)
		default:
			log.Println("Confirm RPC Error")
		}
		ocw.Close()
	})

	display := container.NewVBox(confirm_display, layout.NewSpacer())
	options := container.NewHBox(layout.NewSpacer(), n_confirmButton, y_confirmButton)
	content := container.NewVBox(display, options)

	ocw.SetContent(content)
	ocw.Show()

}

func namePopUp(i int) { /// name change confirmation pop up
	ncw := fyne.CurrentApp().NewWindow("Confirm")
	ncw.SetIcon(resourceDReamTablesIconPng)
	ncw.Resize(fyne.NewSize(MIN_WIDTH-50, 150))
	ncw.SetFixedSize(true)
	var confirm_display = widget.NewLabel("")
	confirm_display.Wrapping = fyne.TextWrapWord

	switch i {
	case 1:
		confirm_display.SetText("0.1 Dero Fee to Change Name\n\nConfirm")
	case 2:
		confirm_display.SetText("0.05 Dero Fee to Remove Address from contract\n\nConfirm")
	default:
		confirm_display.SetText("Error")
	}

	n_confirmButton := widget.NewButton("No", func() {
		ncw.Close()
	})
	y_confirmButton := widget.NewButton("Yes", func() {

		switch i {
		case 1:
			nameChange()
		case 2:
			removeAddress()
		default:
			log.Println("RPC Error")
		}

		ncw.Close()
	})

	display := container.NewVBox(confirm_display, layout.NewSpacer())
	options := container.NewHBox(layout.NewSpacer(), n_confirmButton, y_confirmButton)
	content := container.NewVBox(display, options)

	ncw.SetContent(content)
	ncw.Show()

}

func ownerConfirmPopUp(i int) { /// owner action confirmation
	ocw := fyne.CurrentApp().NewWindow("Confirm")
	ocw.SetIcon(resourceDReamTablesIconPng)
	ocw.Resize(fyne.NewSize(MIN_WIDTH-50, 150))
	ocw.SetFixedSize(true)
	var confirm_display = widget.NewLabel("")
	confirm_display.Wrapping = fyne.TextWrapWord
	s_game := s_gameInput.Text
	s_amt := s_amtInput.Text
	s_end := s_endInput.Text
	s_feed := s_feedInput.Text
	p_pre := p_nameInput.Text
	p_amt := p_amtInput.Text
	p_end := p_endInput.Text
	p_feed := p_feedInput.Text
	s_pay_n := payout_nInput.Text
	s_pay_team := strconv.Itoa(payout_teamInput.SelectedIndex())
	p_pay_price := payout_priceInput.Text
	p_post_price := post_priceInput.Text

	switch i {
	case 1:
		confirm_display.SetText("Game: " + s_game + "\nMinimum: " + s_amt + "\nEnds At: " + s_end + "\nFeed: " + s_feed + "\n\nConfirm")
	case 2:
		confirm_display.SetText("Predicing: " + p_pre + "\nMinimum: " + p_amt + "\nEnds At: " + p_end + "\nFeed: " + p_feed + "\n\nConfirm")
	case 3:
		confirm_display.SetText("Game: " + s_pay_n + "\nTeam: " + teamReturn(s_pay_team) + "\n\nConfirm")
	case 4:
		confirm_display.SetText("Post Price: " + p_post_price + "\n\nConfirm")
	case 5:
		confirm_display.SetText("Final Price: " + p_pay_price + "\n\nConfirm")
	default:
		log.Println("No Confirm Input")
		confirm_display.SetText("Error")
	}

	n_confirmButton := widget.NewButton("No", func() {
		ocw.Close()
	})
	y_confirmButton := widget.NewButton("Yes", func() {
		switch i {
		case 1:
			setSports(convertString_Int(s_end), convertString_Int(s_amt))
		case 2:
			setPrediction(convertString_Int(p_end), convertString_Int(p_amt))
		case 3:
			endSports(payout_nInput.Text, teamReturn(s_pay_team))
		case 4:
			postPrediction(convertString_Int(p_post_price))
		case 5:
			endPredition(convertString_Int(p_pay_price))
		default:
			log.Println("Owner RPC Error")
		}
		ocw.Close()
	})

	display := container.NewVBox(confirm_display, layout.NewSpacer())
	options := container.NewHBox(layout.NewSpacer(), n_confirmButton, y_confirmButton)
	content := container.NewVBox(display, options)

	ocw.SetContent(content)
	ocw.Show()

}

func s_Results(g, gN, l, min, eA, c, tA, tB, tAV, tBV, total string) { /// sports info label

	result, err := strconv.ParseFloat(total, 32)

	if err != nil {
		log.Println("Float Conversion Error", err)
	}

	s := fmt.Sprintf("%f", result/100000)

	human_time, _ := msToTime(eA)
	utc := human_time.String()

	s_sc_displayT.SetText(s_sc_displayT.Text + "\nGame " + gN + " - " + g + "\nLeague: " + l + "\nMinimum: " + min +
		"\nAccepting Until: " + utc + "\nPot Total: " + s + "\nPicks: " + c + "\n" + tA + " Picks: " + tAV + "\n" + tB + " Picks: " + tBV + "\n")

}

func p_initResults(p, amt, eA, c, t, u, d, r, f, m string, post bool) { /// prediction info label, initialized

	skipCount++
	human_time, _ := msToTime(eA)
	utc := human_time.String()
	add := convertString_Int(eA)
	end := strconv.Itoa(add + 86400000)
	utc_end, _ := msToTime(end)

	result, err := strconv.ParseFloat(t, 32)

	if err != nil {
		log.Println("Float Conversion Error", err)
	}

	s := fmt.Sprintf("%f", result/100000)
	if post {
		p_sc_displayT.SetText("SC ID: \n" + P_SC_ID + "\n" + "\n" + p + " Price Posted" +
			"\nPrediction Amount: " + amt + "\nAccepting Until: " + utc + "\nPredictions: " + c +
			"\nRound Pot: " + s + "\nUp Predictions: " + u + "\nDown Predictions: " + d + "\nPosted Price: " + m + "\nPayout After: " + utc_end.String() + "\nTotal Rounds Played: " + r)
	} else {
		p_sc_displayT.SetText("SC ID: \n" + P_SC_ID + "\n" + "\nAccepting " + p + " Predictions " +
			"\nPrediction Amount: " + amt + "\nAccepting Until: " + utc + "\nPredictions: " + c +
			"\nRound Pot: " + s + "\nUp Predictions: " + u + "\nDown Predictions: " + d + "\nPayout After: " + utc_end.String() + "\nTotal Rounds Played: " + r)
	}

	if skipCount == 1 {
		p_sc_displayB.SetText(f + "\nBTC: " + getPrice("BTC") + "\nDero: " + getPrice("DERO") + "\nXMR: " + getPrice("XMR"))
	} else if skipCount == 11 {
		skipCount = 0
	}
}

func p_no_initResults(fr, tx, r, m string) { /// prediction info label, not initialized
	p_sc_displayT.SetText("SC ID: \n" + P_SC_ID + "\n" + "\nNot Accepting Predictions\n\nLast Round Mark: " + m +
		"\nLast Round Results: " + fr + "\nLast Round TXID: " + tx + "\n\nTotal Rounds Played: " + r)

	p_sc_displayB.SetText("")
}

func printLeaders() { /// make leaderboard
	keys := make([]string, 0, len(leaders))

	for key := range leaders {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return leaders[keys[i]] > leaders[keys[j]]
	})

	for _, k := range keys {
		leadersTotal = append(leadersTotal, k+": "+strconv.FormatUint(leaders[k], 10))
	}

	leaders_list.Refresh()
}
