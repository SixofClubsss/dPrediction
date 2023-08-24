package prediction

import (
	"strings"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var S dreams.ContainerStack

// dSports tab layout
func LayoutSportsItems(d *dreams.AppObject) *fyne.Container {
	S.LeftLabel = widget.NewLabel("")
	S.RightLabel = widget.NewLabel("")
	S.RightLabel.SetText("dReams Balance: " + rpc.DisplayBalance("dReams") + "      Dero Balance: " + rpc.DisplayBalance("Dero") + "      Height: " + rpc.Wallet.Display.Height)

	Sports.Info = widget.NewLabel("SCID:\n\n" + Sports.Contract + "\n")
	Sports.Info.Wrapping = fyne.TextWrapWord

	sports_content := container.NewVBox(Sports.Info)
	sports_scroll := container.NewVScroll(sports_content)
	sports_scroll.SetMinSize(fyne.NewSize(180, 500))

	check_box := container.NewVBox(SportsConnectedBox())

	contract_scroll := container.NewHScroll(SportsContractEntry())
	contract_scroll.SetMinSize(fyne.NewSize(600, 35.1875))
	contract_cont := container.NewHBox(contract_scroll, check_box)

	Sports.Game_select = widget.NewSelect(Sports.Game_options, func(s string) {
		split := strings.Split(s, "   ")
		a, b := GetSportsTeams(Sports.Contract, split[0])
		if Sports.Game_select.SelectedIndex() >= 0 {
			Sports.Multi.Show()
			Sports.ButtonA.Show()
			Sports.ButtonB.Show()
			Sports.ButtonA.Text = a
			Sports.ButtonA.Refresh()
			Sports.ButtonB.Text = b
			Sports.ButtonB.Refresh()
		} else {
			Sports.Multi.Hide()
			Sports.ButtonA.Hide()
			Sports.ButtonB.Hide()
		}
	})

	Sports.Game_select.PlaceHolder = "Select Game #"
	Sports.Game_select.Hide()

	var multi_options = []string{"1x", "3x", "5x"}
	Sports.Multi = widget.NewRadioGroup(multi_options, nil)
	Sports.Multi.SetSelected("1x")
	Sports.Multi.Horizontal = true
	Sports.Multi.Required = true
	Sports.Multi.Hide()

	Sports.ButtonA = widget.NewButton("TEAM A", nil)
	Sports.ButtonA.Hide()

	Sports.ButtonB = widget.NewButton("TEAM B", nil)
	Sports.ButtonB.Hide()

	sports_multi := container.NewCenter(Sports.Multi)
	Sports.Container = container.NewVBox(
		sports_multi,
		Sports.Game_select,
		Sports.ButtonA,
		Sports.ButtonB)

	Sports.Container.Hide()

	sports_left := container.NewVBox(
		contract_cont,
		sports_scroll,
		layout.NewSpacer(),
		Sports.Container)

	epl := widget.NewLabel("")
	epl.Wrapping = fyne.TextWrapWord
	epl_scroll := container.NewVScroll(epl)
	mls := widget.NewLabel("")
	mls.Wrapping = fyne.TextWrapWord
	mls_scroll := container.NewVScroll(mls)
	nba := widget.NewLabel("")
	nba.Wrapping = fyne.TextWrapWord
	nba_scroll := container.NewVScroll(nba)
	nfl := widget.NewLabel("")
	nfl.Wrapping = fyne.TextWrapWord
	nfl_scroll := container.NewVScroll(nfl)
	nhl := widget.NewLabel("")
	nhl.Wrapping = fyne.TextWrapWord
	nhl_scroll := container.NewVScroll(nhl)
	mlb := widget.NewLabel("")
	mlb.Wrapping = fyne.TextWrapWord
	mlb_scroll := container.NewVScroll(mlb)
	bellator := widget.NewLabel("")
	bellator.Wrapping = fyne.TextWrapWord
	bellator_scroll := container.NewVScroll(bellator)
	ufc := widget.NewLabel("")
	ufc.Wrapping = fyne.TextWrapWord
	ufc_scroll := container.NewVScroll(ufc)
	score_tabs := container.NewAppTabs(
		container.NewTabItem("EPL", epl_scroll),
		container.NewTabItem("MLS", mls_scroll),
		container.NewTabItem("NBA", nba_scroll),
		container.NewTabItem("NFL", nfl_scroll),
		container.NewTabItem("NHL", nhl_scroll),
		container.NewTabItem("MLB", mlb_scroll),
		container.NewTabItem("Bellator", bellator_scroll),
		container.NewTabItem("UFC", ufc_scroll))

	score_tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "EPL":
			go GetScores(epl, "EPL")
		case "MLS":
			go GetScores(mls, "MLS")
		case "NBA":
			go GetScores(nba, "NBA")
		case "NFL":
			go GetScores(nfl, "NFL")
		case "NHL":
			go GetScores(nhl, "NHL")
		case "MLB":
			go GetScores(mlb, "MLB")
		case "Bellator":
			go GetMmaResults(bellator, "Bellator")
		case "UFC":
			go GetMmaResults(ufc, "UFC")
		default:
		}
	}

	Sports.Settings.Unlock = widget.NewButton("Unlock dSports Contracts", nil)
	Sports.Settings.Unlock.Hide()

	Sports.Settings.New = widget.NewButton("New dSports Contract", nil)
	Sports.Settings.New.Hide()

	unlock_cont := container.NewVBox(
		Sports.Settings.Unlock,
		Sports.Settings.New)

	Sports.Settings.Menu = widget.NewButton("Owner Options", func() {
		go ownersMenu()
	})
	Sports.Settings.Menu.Hide()

	owner_buttons := container.NewAdaptiveGrid(2, container.NewMax(Sports.Settings.Menu), unlock_cont)
	owned_tab := container.NewBorder(nil, owner_buttons, nil, nil, SportsOwned())

	tabs := container.NewAppTabs(
		container.NewTabItem("Contracts", layout.NewSpacer()),
		container.NewTabItem("Favorites", SportsFavorites()),
		container.NewTabItem("Owned", owned_tab),
		container.NewTabItem("Scores", score_tabs),
		container.NewTabItem("Payouts", SportsPayouts()))

	tabs.SelectIndex(0)
	tabs.Selected().Content = SportsListings(tabs)

	tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Contracts":
			go PopulateSports(nil)
		default:

		}
	}

	max := container.NewMax(bundle.Alpha120, tabs)

	Sports.ButtonA.OnTapped = func() {
		if len(Sports.Contract) == 64 {
			max.Objects[1] = ConfirmAction(3, Sports.ButtonA.Text, Sports.ButtonB.Text, max.Objects, tabs)
			max.Objects[1].Refresh()
		}
	}
	Sports.ButtonA.Hide()

	Sports.ButtonB.OnTapped = func() {
		if len(Sports.Contract) == 64 {
			max.Objects[1] = ConfirmAction(4, Sports.ButtonA.Text, Sports.ButtonB.Text, max.Objects, tabs)
			max.Objects[1].Refresh()
		}
	}

	Sports.Settings.Unlock.OnTapped = func() {
		max.Objects[1] = newSportsConfirm(1, max.Objects, tabs)
		max.Objects[1].Refresh()
	}

	Sports.Settings.New.OnTapped = func() {
		max.Objects[1] = newSportsConfirm(2, max.Objects, tabs)
		max.Objects[1].Refresh()
	}

	sports_label := container.NewHBox(S.LeftLabel, layout.NewSpacer(), S.RightLabel)
	sports_box := container.NewHSplit(sports_left, max)

	S.DApp = container.NewBorder(
		dwidget.LabelColor(sports_label),
		nil,
		nil,
		nil,
		sports_box)

	go fetch(d)

	return container.NewMax(S.DApp)
}
