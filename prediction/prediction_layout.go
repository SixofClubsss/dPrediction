package prediction

import (
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var P dreams.ContainerStack

// dPrediction tab layout
func LayoutPredictItems(d *dreams.AppObject) *fyne.Container {
	P.LeftLabel = widget.NewLabel("")
	P.RightLabel = widget.NewLabel("")
	P.RightLabel.SetText("dReams Balance: " + rpc.DisplayBalance("dReams") + "      Dero Balance: " + rpc.DisplayBalance("Dero") + "      Height: " + rpc.Wallet.Display.Height)

	Predict.info = widget.NewLabel("SCID:\n\n" + Predict.Contract.SCID + "\n")
	Predict.info.Wrapping = fyne.TextWrapWord
	Predict.prices = widget.NewLabel("")

	predict_info := container.NewVBox(Predict.info, Predict.prices)
	predict_scroll := container.NewScroll(predict_info)
	predict_scroll.SetMinSize(fyne.NewSize(540, 500))

	check_box := container.NewVBox(PredictConnectedBox())

	Predict.higher = widget.NewButton("Higher", nil)
	Predict.higher.Hide()

	Predict.lower = widget.NewButton("Lower", nil)
	Predict.lower.Hide()

	Predict.Container = container.NewVBox(Predict.higher, Predict.lower)
	Predict.Container.Hide()

	Predict.Contract.unlock = widget.NewButton("Unlock dPrediction Contract", nil)
	Predict.Contract.unlock.Hide()

	Predict.Contract.new = widget.NewButton("New dPrediction Contract", nil)
	Predict.Contract.new.Hide()

	unlock_cont := container.NewVBox(Predict.Contract.unlock, Predict.Contract.new)

	Predict.Contract.menu = widget.NewButton("Owner Options", func() {
		go ownersMenu()
	})
	Predict.Contract.menu.Hide()

	owner_buttons := container.NewAdaptiveGrid(2, container.NewStack(Predict.Contract.menu), unlock_cont)
	owned_tab := container.NewBorder(nil, owner_buttons, nil, nil, PredictionOwned())

	tabs := container.NewAppTabs(
		container.NewTabItem("Contracts", layout.NewSpacer()),
		container.NewTabItem("Favorites", PredictionFavorites()),
		container.NewTabItem("Owned", owned_tab))

	tabs.SelectIndex(0)
	tabs.Selected().Content = PredictionListings(tabs)

	tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Contracts":
			go PopulatePredictions(nil)
		default:
		}
	}

	max := container.NewStack(bundle.Alpha120, tabs)

	Predict.higher.OnTapped = func() {
		if len(Predict.Contract.SCID) == 64 {
			max.Objects[1] = ConfirmAction(2, "", "", max.Objects, tabs)
			max.Objects[1].Refresh()
		}
	}

	Predict.lower.OnTapped = func() {
		if len(Predict.Contract.SCID) == 64 {
			max.Objects[1] = ConfirmAction(1, "", "", max.Objects, tabs)
			max.Objects[1].Refresh()
		}
	}

	Predict.Contract.unlock.OnTapped = func() {
		max.Objects[1] = newPredictConfirm(1, max.Objects, tabs)
		max.Objects[1].Refresh()
	}

	Predict.Contract.new.OnTapped = func() {
		max.Objects[1] = newPredictConfirm(2, max.Objects, tabs)
		max.Objects[1].Refresh()
	}

	contract_scroll := container.NewHScroll(PredictionContractEntry())
	contract_scroll.SetMinSize(fyne.NewSize(600, 35.1875))
	contract_cont := container.NewHBox(contract_scroll, check_box)

	predict_content := container.NewVBox(
		contract_cont,
		predict_scroll,
		layout.NewSpacer(),
		Predict.Container)

	predict_label := container.NewHBox(P.LeftLabel, layout.NewSpacer(), P.RightLabel)
	predict_box := container.NewHSplit(predict_content, max)

	P.DApp = container.NewBorder(
		dwidget.LabelColor(predict_label),
		nil,
		nil,
		nil,
		predict_box)

	return container.NewStack(P.DApp)
}
