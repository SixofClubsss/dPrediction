package prediction

import (
	"fmt"
	"strings"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var S dwidget.ContainerStack

// dSports tab layout
func LayoutSports(d *dreams.AppObject) *fyne.Container {
	S.Left.Label = widget.NewLabel("")

	S.Right.Label = widget.NewLabel("")
	S.Right.SetUpdate(dreams.SetBalanceLabelText)

	Sports.info = widget.NewLabel("SCID:\n\n" + Sports.Contract.SCID + "\n")
	Sports.info.Wrapping = fyne.TextWrapWord

	sports_content := container.NewVBox(Sports.info)
	sports_scroll := container.NewVScroll(sports_content)
	sports_scroll.SetMinSize(fyne.NewSize(180, 500))

	check_box := container.NewVBox(SportsConnectedBox())

	Sports.gameSelect = widget.NewSelect(Sports.games, func(s string) {
		split := strings.Split(s, "   ")
		a, b := GetSportsTeams(Sports.Contract.SCID, split[0])
		if Sports.gameSelect.SelectedIndex() >= 0 {
			Sports.multi.Show()
			Sports.buttonA.Show()
			Sports.buttonB.Show()
			Sports.buttonA.Text = a
			Sports.buttonA.Refresh()
			Sports.buttonB.Text = b
			Sports.buttonB.Refresh()
		} else {
			Sports.multi.Hide()
			Sports.buttonA.Hide()
			Sports.buttonB.Hide()
		}
	})

	Sports.gameSelect.PlaceHolder = "Select Game #"
	Sports.gameSelect.Hide()

	var multi_options = []string{"1x", "3x", "5x"}
	Sports.multi = widget.NewRadioGroup(multi_options, nil)
	Sports.multi.SetSelected("1x")
	Sports.multi.Horizontal = true
	Sports.multi.Required = true
	Sports.multi.Hide()

	Sports.buttonA = widget.NewButton("TEAM A", nil)
	Sports.buttonA.Importance = widget.HighImportance
	Sports.buttonA.Hide()

	Sports.buttonB = widget.NewButton("TEAM B", nil)
	Sports.buttonB.Importance = widget.HighImportance
	Sports.buttonB.Hide()

	sports_multi := container.NewCenter(Sports.multi)
	Sports.Container = container.NewVBox(
		sports_multi,
		Sports.gameSelect,
		Sports.buttonA,
		Sports.buttonB)

	Sports.Container.Hide()

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

	loading_img := canvas.NewImageFromResource(resourceDServiceCirclePng)
	loading_img.SetMinSize(fyne.NewSize(140, 140))

	score_tabs.OnSelected = func(ti *container.TabItem) {
		go func() {
			ti.Content = container.NewStack(container.NewCenter(loading_img, dwidget.NewCanvasText(fmt.Sprintf("Loading %s...", ti.Text), 17, fyne.TextAlignCenter)), widget.NewProgressBarInfinite())
			switch ti.Text {
			case "EPL":
				GetScores(epl, "EPL")
				ti.Content = epl_scroll
			case "MLS":
				GetScores(mls, "MLS")
				ti.Content = mls_scroll
			case "NBA":
				GetScores(nba, "NBA")
				ti.Content = nba_scroll
			case "NFL":
				GetScores(nfl, "NFL")
				ti.Content = nfl_scroll
			case "NHL":
				GetScores(nhl, "NHL")
				ti.Content = nhl_scroll
			case "MLB":
				GetScores(mlb, "MLB")
				ti.Content = mlb_scroll
			case "Bellator":
				GetMmaResults(bellator, "Bellator")
				ti.Content = bellator_scroll
			case "UFC":
				GetMmaResults(ufc, "UFC")
				ti.Content = ufc_scroll
			default:

			}
		}()
	}

	Sports.Contract.unlock = widget.NewButton("Unlock dSports Contracts", nil)
	Sports.Contract.unlock.Importance = widget.HighImportance
	Sports.Contract.unlock.Hide()

	Sports.Contract.new = widget.NewButton("New dSports Contract", nil)
	Sports.Contract.new.Importance = widget.HighImportance
	Sports.Contract.new.Hide()

	unlock_cont := container.NewVBox(
		Sports.Contract.unlock,
		Sports.Contract.new)

	Sports.Contract.menu = widget.NewButton("Owner Options", func() {
		go ownersMenu()
	})
	Sports.Contract.menu.Importance = widget.HighImportance
	Sports.Contract.menu.Hide()

	owner_buttons := container.NewAdaptiveGrid(2, container.NewStack(Sports.Contract.menu), unlock_cont)
	owned_tab := container.NewBorder(nil, owner_buttons, nil, nil, SportsOwned())

	tabs := container.NewAppTabs(
		container.NewTabItem("Contracts", SportsListings(d)),
		container.NewTabItem("Favorites", SportsFavorites()),
		container.NewTabItem("Owned", owned_tab),
		container.NewTabItem("Scores", score_tabs),
		container.NewTabItem("Payouts", SportsPayouts()))

	tabs.SelectIndex(0)

	tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Contracts":
			go PopulateSports(nil)
		default:

		}
	}

	max := container.NewStack(bundle.Alpha120, tabs)

	Sports.buttonA.OnTapped = func() {
		if len(Sports.Contract.SCID) == 64 {
			ConfirmAction(3, Sports.buttonA.Text, Sports.buttonB.Text, d)
		}
	}

	Sports.buttonB.OnTapped = func() {
		if len(Sports.Contract.SCID) == 64 {
			ConfirmAction(4, Sports.buttonA.Text, Sports.buttonB.Text, d)
		}
	}

	Sports.Contract.unlock.OnTapped = func() {
		newSportsConfirm(1, d)
	}

	Sports.Contract.new.OnTapped = func() {
		newSportsConfirm(2, d)
	}

	contract_scroll := container.NewHScroll(SportsContractEntry())
	contract_scroll.SetMinSize(fyne.NewSize(600, 35.1875))
	contract_cont := container.NewHBox(contract_scroll, check_box)

	sports_left := container.NewVBox(
		contract_cont,
		sports_scroll,
		layout.NewSpacer(),
		Sports.Container)

	sports_label := container.NewHBox(S.Left.Label, layout.NewSpacer(), S.Right.Label)
	sports_box := container.NewHSplit(sports_left, max)

	S.DApp = container.NewBorder(
		dwidget.LabelColor(sports_label),
		nil,
		nil,
		nil,
		sports_box)

	go fetch(d)

	return container.NewStack(S.DApp)
}
