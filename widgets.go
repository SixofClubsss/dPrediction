package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

// / wallet tab widgets
var rpcLoginInput = widget.NewPasswordEntry()
var rpcWalletInput = widget.NewEntry()
var currentHeight = widget.NewEntry()
var walletBalance = widget.NewEntry()

var daemon_options = []string{"REMOTE", "MAINNET", "SIMULATOR", "CUSTOM"} /// set daemon select menu
var daemon_dropDown = widget.NewSelect(daemon_options, func(s string) {   /// do when select changes
	log.Println("Daemon Set To:", s)
	whichDaemon(s)

})

var daemonCheckBox = widget.NewCheck("Daemon Connected", func(value bool) {
	/// do something on change
})

var walletCheckBox = widget.NewCheck("Wallet Connected", func(value bool) {

})

func rpcLoginEdit() fyne.Widget { /// user:pass password entry
	rpcLoginInput.SetPlaceHolder("Enter user:pass")
	rpcWalletInput.OnChanged = func(input string) {
		enablePrediction(false)
		walletCheckBox.SetChecked(false)
		walletConnectBool = false
		walletAddress = ""
	}

	return rpcLoginInput
}

func rpcWalletEdit() fyne.Widget { /// wallet rpc address entry
	rpcWalletInput.SetPlaceHolder("Wallet RPC Address")
	rpcWalletInput.SetText("127.0.0.1:30000")
	rpcWalletInput.OnChanged = func(input string) {
		enablePrediction(false)
		walletCheckBox.SetChecked(false)
		walletConnectBool = false
		walletAddress = ""
	}

	return rpcWalletInput
}

func rpcConnectButton() fyne.Widget { /// rpc connect button
	button := widget.NewButton("Connect", func() {
		log.Println("Connect Pressed")
		pre := "http://"
		suff := "/json_rpc"
		walletAddress = pre + rpcWalletInput.Text + suff
		getAddress()
		ifPlaying()
		if nameInput.Text != "" {
			nameInput.Disable()
			changeNameButton.Enable()
		} else {
			nameInput.Enable()
			changeNameButton.Disable()
		}
		nameInput.Refresh()
	})
	button.Resize(fyne.NewSize(100, 42))
	button.Move(fyne.NewPos(270, 702))

	return button
}

func daemonSelectOption() fyne.Widget { /// daemon select menu
	daemon_dropDown.SetSelectedIndex(0)

	return daemon_dropDown
}

func daemonConnectBox() fyne.Widget { /// daemon check box
	daemonCheckBox.Disable()

	return daemonCheckBox
}

func walletConnectBox() fyne.Widget { /// wallet check box
	walletCheckBox.Disable()

	return walletCheckBox
}

func heightDisplay() fyne.Widget { /// height display
	currentHeight.SetText("Height:")
	currentHeight.Disable()

	return currentHeight
}

func balanceDisplay() fyne.Widget { /// balance display
	walletBalance.SetText("Balance:")
	walletBalance.Disable()

	return walletBalance
}

var owner_button = widget.NewButton("Owner", func() {
	log.Println("Owner Pressed")
	ownersMenu()
})

func ownerButton() fyne.Widget {
	owner_button.Disable()

	return owner_button
}

// / sports widgets
var s_sc_displayT = widget.NewLabel("")
var game_select = widget.NewEntry()

var a_button = widget.NewButton("TEAM A", func() {
	log.Println("TEAM A Pressed")
	confirmPopUp(3)
})

var b_button = widget.NewButton("TEAM B", func() {
	log.Println("TEAM B Pressed")
	confirmPopUp(4)
})

func teamA() fyne.Widget { /// team A button
	a_button.Disable()

	return a_button
}

func teamB() fyne.Widget { /// team B button
	b_button.Disable()

	return b_button
}

func s_scDisplayTop() fyne.Widget { /// sports display label
	s_sc_displayT.Wrapping = fyne.TextWrapWord
	s_sc_displayT.SetText("SC ID: \n" + S_SC_ID + "\n")

	return s_sc_displayT
}

func gameOptions() fyne.Widget { /// game number entry
	game_select.SetPlaceHolder("Game #:")
	/// TODO add validator
	game_select.Disable()

	return game_select
}

// / prediction widgets
var p_sc_displayT = widget.NewLabel("Loading Data...")
var p_sc_displayB = widget.NewLabel("")
var nameInput = widget.NewEntry()

var changeNameButton = widget.NewButton("Change Name", func() {
	log.Println("Change Pressed")
	if nameInput.Disabled() {
		nameInput.Enable()
	} else {
		namePopUp(1)
	}
})

func p_scDisplayTop() fyne.Widget { /// prediction top label
	p_sc_displayT.Wrapping = fyne.TextWrapWord

	return p_sc_displayT
}

func p_scDisplayBottom() fyne.Widget { /// prediction bottom label
	p_sc_displayB.Wrapping = fyne.TextWrapWord

	return p_sc_displayB
}

func nameEdit() fyne.Widget { /// name entry
	nameInput.SetPlaceHolder("Name")
	nameInput.OnChanged = func(input string) {
		nameInput.Validator = validation.NewRegexp(`\w{3,}`, "Three Letters Minimum")
		nameInput.Validate()
		nameInput.Refresh()
	}
	nameInput.Disable()

	return nameInput
}

func change() fyne.Widget { /// change name button
	changeNameButton.Resize(fyne.NewSize(175, 45))
	changeNameButton.Move(fyne.NewPos(195, 605))
	changeNameButton.Disable()

	return changeNameButton
}

var upButton = widget.NewButton("Higher", func() {
	log.Println("Higher Pressed")
	confirmPopUp(2)
})

var downButton = widget.NewButton("Lower", func() {
	log.Println("Lower Pressed")
	confirmPopUp(1)
})

func higher() fyne.Widget { /// higher prediction button
	upButton.Disable()

	return upButton
}

func lower() fyne.Widget { /// lower prediction button
	downButton.Disable()

	return downButton
}

// / leaderboard widgets
var leadersTotal = []string{"Loading Data..."}

var leaders_list = widget.NewList(
	func() int {
		return len(leadersTotal)
	},
	func() fyne.CanvasObject {
		return widget.NewLabel("")
	},
	func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(leadersTotal[i])
	})

func leadersDisplay() fyne.Widget { /// leaderboard display
	leaders_list.Resize(fyne.NewSize(360, 680))
	leaders_list.Move(fyne.NewPos(5, 10))

	return leaders_list
}

var removeButton = widget.NewButton("Remove", func() {
	log.Println("Remove Pressed")
	namePopUp(1)
})

func remove() fyne.Widget { /// remove address button
	removeButton.Disable()

	return removeButton
}
