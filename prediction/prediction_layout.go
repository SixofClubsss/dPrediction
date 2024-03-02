package prediction

import (
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var P dwidget.ContainerStack

// dPrediction tab layout
func LayoutPredictions(d *dreams.AppObject) *fyne.Container {
	P.Left.Label = widget.NewLabel("")

	P.Right.Label = widget.NewLabel("")
	P.Right.SetUpdate(dreams.SetBalanceLabelText)

	Predict.info = widget.NewLabel("SCID:\n\n" + Predict.Contract.SCID + "\n")
	Predict.info.Wrapping = fyne.TextWrapWord
	Predict.prices = widget.NewLabel("")

	predict_info := container.NewVBox(Predict.info, Predict.prices)
	predict_scroll := container.NewScroll(predict_info)
	predict_scroll.SetMinSize(fyne.NewSize(540, 500))

	check_box := container.NewVBox(PredictConnectedBox())

	Predict.higher = widget.NewButtonWithIcon("Higher", dreams.FyneIcon("arrowUp"), nil)
	Predict.higher.Importance = widget.HighImportance
	Predict.higher.Hide()

	Predict.lower = widget.NewButtonWithIcon("Lower", dreams.FyneIcon("arrowDown"), nil)
	Predict.lower.Importance = widget.HighImportance
	Predict.lower.Hide()

	Predict.Container = container.NewVBox(Predict.higher, Predict.lower)
	Predict.Container.Hide()

	Predict.Contract.unlock = widget.NewButton("Unlock dPrediction Contract", nil)
	Predict.Contract.unlock.Importance = widget.HighImportance
	Predict.Contract.unlock.Hide()

	Predict.Contract.new = widget.NewButton("New dPrediction Contract", nil)
	Predict.Contract.new.Importance = widget.HighImportance
	Predict.Contract.new.Hide()

	unlock_cont := container.NewVBox(Predict.Contract.unlock, Predict.Contract.new)

	Predict.Contract.menu = widget.NewButton("Owner Options", func() {
		go ownersMenu()
	})
	Predict.Contract.menu.Importance = widget.HighImportance
	Predict.Contract.menu.Hide()

	owner_buttons := container.NewAdaptiveGrid(2, container.NewStack(Predict.Contract.menu), unlock_cont)
	owned_tab := container.NewBorder(nil, owner_buttons, nil, nil, PredictionOwned())

	tabs := container.NewAppTabs(
		container.NewTabItem("Contracts", PredictionListings(d)),
		container.NewTabItem("Favorites", PredictionFavorites()),
		container.NewTabItem("Owned", owned_tab))

	tabs.SelectIndex(0)

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
			ConfirmAction(2, "", "", d)
		}
	}

	Predict.lower.OnTapped = func() {
		if len(Predict.Contract.SCID) == 64 {
			ConfirmAction(1, "", "", d)
		}
	}

	Predict.Contract.unlock.OnTapped = func() {
		newPredictConfirm(1, d)
	}

	Predict.Contract.new.OnTapped = func() {
		newPredictConfirm(2, d)
	}

	contract_scroll := container.NewHScroll(PredictionContractEntry())
	contract_scroll.SetMinSize(fyne.NewSize(600, 35.1875))
	contract_cont := container.NewHBox(contract_scroll, check_box)

	predict_content := container.NewVBox(
		contract_cont,
		predict_scroll,
		layout.NewSpacer(),
		Predict.Container)

	predict_label := container.NewHBox(P.Left.Label, layout.NewSpacer(), P.Right.Label)
	predict_box := container.NewHSplit(predict_content, max)

	P.DApp = container.NewBorder(
		dwidget.LabelColor(predict_label),
		nil,
		nil,
		nil,
		predict_box)

	return container.NewStack(P.DApp)
}
