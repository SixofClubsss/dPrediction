package prediction

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type predictObjects struct {
	*fyne.Container
	init       bool
	owner      bool
	amount     uint64
	buffer     int64
	prediction string
	feed       string
	info       *widget.Label
	prices     *widget.Label
	Public     dwidget.Lists
	Favorites  dwidget.Lists
	Owned      dwidget.Lists
	higher     *widget.Button
	lower      *widget.Button
	Contract   struct {
		SCID   string
		unlock *widget.Button
		new    *widget.Button
		menu   *widget.Button
		check  *widget.Check
		entry  *widget.SelectEntry
	}
}

var Predict predictObjects

// Disable dPrediction objects
func disablePredictions(d bool) {
	if d {
		Predict.Container.Hide()
	} else {
		Predict.Container.Show()
	}
	Predict.Container.Refresh()
}

// Check box for dPrediction SCID
//   - Hides prediction controls on disconnect
func PredictConnectedBox() fyne.Widget {
	Predict.Contract.check = widget.NewCheck("", func(b bool) {
		if !b {
			Predict.higher.Hide()
			Predict.lower.Hide()

		}
	})
	Predict.Contract.check.Disable()

	return Predict.Contract.check
}

// Entry for dPrediction SCID
//   - Bound to Predict.Settings.SCID
//   - Checks for valid SCID on changed
func PredictionContractEntry() fyne.Widget {
	options := []string{}
	Predict.Contract.entry = widget.NewSelectEntry(options)
	Predict.Contract.entry.PlaceHolder = "Contract Address: "
	Predict.Contract.entry.OnCursorChanged = func() {
		if rpc.Daemon.IsConnected() {
			go func() {
				if len(Predict.Contract.SCID) == 64 {
					yes := ValidBetContract(Predict.Contract.SCID)
					if yes {
						Predict.Contract.check.SetChecked(true)
					} else {
						Predict.Owned.List.UnselectAll()
						Predict.Public.List.UnselectAll()
						Predict.Favorites.List.UnselectAll()
						Predict.Contract.check.SetChecked(false)
					}
				} else {
					Predict.Owned.List.UnselectAll()
					Predict.Public.List.UnselectAll()
					Predict.Favorites.List.UnselectAll()
					Predict.Contract.check.SetChecked(false)
				}
			}()
		}
	}

	this := binding.BindString(&Predict.Contract.SCID)
	Predict.Contract.entry.Bind(this)

	return Predict.Contract.entry
}

// When called, enable and show dPrediction controls
func ShowPredictionControls() {
	disablePredictions(false)
	Predict.higher.Show()
	Predict.lower.Show()
}

// Routine when dPrediction SCID is clicked
//   - Sets label info and controls
//   - item returned for adding and removing favorites
func setPredictionControls(str string) (item string) {
	split := strings.Split(str, "   ")
	if len(split) >= 3 {
		trimmed := strings.Trim(split[2], " ")
		if len(trimmed) == 64 {
			item = str
			Predict.Contract.entry.SetText(trimmed)
			go SetPredictionInfo(trimmed)
			if CheckActivePrediction(trimmed) {
				ShowPredictionControls()
			} else {
				disablePredictions(true)
			}
		}
	}

	return
}

// Sets dPrediction info label
func SetPredictionInfo(scid string) {
	info := GetPrediction(scid)
	if info != "" {
		Predict.info.SetText(info)
		Predict.info.Refresh()
	}
}

// Update price feed for dPrediction display
func SetPredictionPrices() {
	if rpc.Daemon.IsConnected() {
		_, btc := menu.GetPrice("BTC-USDT", "Prediction")
		_, dero := menu.GetPrice("DERO-USDT", "Prediction")
		_, xmr := menu.GetPrice("XMR-USDT", "Prediction")
		/// custom feed with rpc.Display.P_feed
		prices := "Current Price feed from dReams Client\nBTC: " + btc + "\nDERO: " + dero + "\nXMR: " + xmr

		Predict.prices.SetText(prices)
	}
}

// List object for populating public dPrediction contracts, with rating and add favorite controls
//   - Pass tab for action confirmation reset
func PredictionListings(d *dreams.AppObject) fyne.CanvasObject {
	Predict.Public.List = widget.NewList(
		func() int {
			return len(Predict.Public.SCIDs)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(container.NewStack(canvas.NewImageFromImage(nil)), widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(Predict.Public.SCIDs[i])
			if Predict.Public.SCIDs[i][0:2] != "  " {
				var key string
				split := strings.Split(Predict.Public.SCIDs[i], "   ")
				if len(split) >= 3 {
					trimmed := strings.Trim(split[2], " ")
					if len(trimmed) == 64 {
						key = trimmed
					}
				}

				badge := canvas.NewImageFromResource(menu.DisplayRating(menu.Control.Ratings[key]))
				badge.SetMinSize(fyne.NewSize(35, 35))
				o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = badge
			}
		})

	var item string

	Predict.Public.List.OnSelected = func(id widget.ListItemID) {
		if id != 0 && gnomes.IsConnected() {
			item = setPredictionControls(Predict.Public.SCIDs[id])
			Predict.Favorites.List.UnselectAll()
			Predict.Owned.List.UnselectAll()
		} else {
			disablePredictions(true)
		}
	}

	save := widget.NewButton("Favorite", func() {
		var have bool
		for _, f := range Predict.Favorites.SCIDs {
			if item == f {
				have = true
			}
		}

		if !have {
			Predict.Favorites.SCIDs = append(Predict.Favorites.SCIDs, item)
			sort.Strings(Predict.Favorites.SCIDs)
		}
	})
	save.Importance = widget.LowImportance

	rate := widget.NewButton("Rate", func() {
		if len(Predict.Contract.SCID) == 64 {
			if !gnomes.CheckOwner(Predict.Contract.SCID) {
				menu.RateConfirm(Predict.Contract.SCID, d)
			} else {
				dialog.NewInformation("Can't rate", "You are the owner of this SCID", d.Window).Show()
				logger.Warnln("[Predictions] Can't rate, you own this contract")
			}
		}
	})
	rate.Importance = widget.LowImportance

	return container.NewBorder(
		nil,
		container.NewBorder(nil, nil, save, rate, layout.NewSpacer()),
		nil,
		nil,
		Predict.Public.List)
}

// List object for populating favorite dPrediction contracts, with remove favorite control
func PredictionFavorites() fyne.CanvasObject {
	Predict.Favorites.List = widget.NewList(
		func() int {
			return len(Predict.Favorites.SCIDs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(Predict.Favorites.SCIDs[i])
		})

	var item string

	Predict.Favorites.List.OnSelected = func(id widget.ListItemID) {
		if gnomes.IsConnected() {
			item = setPredictionControls(Predict.Favorites.SCIDs[id])
			Predict.Public.List.UnselectAll()
			Predict.Owned.List.UnselectAll()
		} else {
			disablePredictions(true)
		}
	}

	remove := widget.NewButton("Remove", func() {
		if len(Predict.Favorites.SCIDs) > 0 {
			Predict.Favorites.List.UnselectAll()
			new := Predict.Favorites.SCIDs
			for i := range new {
				if new[i] == item {
					copy(new[i:], new[i+1:])
					new[len(new)-1] = ""
					new = new[:len(new)-1]
					Predict.Favorites.SCIDs = new
					break
				}
			}
		}
		Predict.Favorites.List.Refresh()
		sort.Strings(Predict.Favorites.SCIDs)
	})
	remove.Importance = widget.LowImportance

	return container.NewBorder(
		nil,
		container.NewBorder(nil, nil, nil, remove, layout.NewSpacer()),
		nil,
		nil,
		Predict.Favorites.List)
}

// List object for populating owned dPrediction contracts
func PredictionOwned() fyne.CanvasObject {
	Predict.Owned.List = widget.NewList(
		func() int {
			return len(Predict.Owned.SCIDs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(Predict.Owned.SCIDs[i])
		})

	Predict.Owned.List.OnSelected = func(id widget.ListItemID) {
		if gnomes.IsConnected() {
			setPredictionControls(Predict.Owned.SCIDs[id])
			Predict.Public.List.UnselectAll()
			Predict.Favorites.List.UnselectAll()
		} else {
			disablePredictions(true)
		}
	}

	return Predict.Owned.List
}

// Formats initialized dPrediction info string
//   - p defines prediction
//   - amt is Dero value
//   - eA is prediction end time
//   - c is number of current predictions
//   - to is current total prediction Dero pot value
//   - u is higher predictions
//   - d is lower predictions
//   - r is total completed prediction rounds
//   - f is prediction feed
//   - m is prediction mark
//   - ta, tb, tc are current contract time frames
func P_initResults(p, amt, eA, c, to, u, d, r, f, m string, ta, tb, tc int) (info string) {
	end_time, _ := rpc.MsToTime(eA)
	utc := end_time.String()
	add := rpc.StringToInt(eA)
	end := strconv.Itoa(add + (tc * 1000))
	end_pay, _ := rpc.MsToTime(end)
	rf := strconv.Itoa(tb / 60)

	result, err := strconv.ParseFloat(to, 32)

	if err != nil {
		logger.Errorln("[Predictions]", err)
	}

	s := fmt.Sprintf("%.5f", result/100000)

	now := time.Now().Unix()
	done := now > end_time.Unix()

	if done {
		mark := m
		if mark == "0" {
			mark = ""
		}

		var wfp string
		if mark == "" {
			wfp = "   Waiting for Mark"
		} else {
			wfp = "   Waiting for Payout"
		}

		if isOnChainPrediction(p) {
			info = "SCID:\n\n" + Predict.Contract.SCID + "\n\n" + p + wfp + "\n\nNode: " + f + "\n\nMark: " + mark + "\nRound Pot: " + s +
				" Dero\n\nPredictions: " + c + "\nHigher Predictions: " + u + "\nLower Predictions: " + d +
				"\n\nPayout After: " + end_pay.String() + "\nRefund if not paid within " + rf + " minutes\n\nRounds Completed: " + r
		} else {
			info = "SCID:\n\n" + Predict.Contract.SCID + "\n\n" + p + wfp + "\n\nMark: " + mark + "\nRound Pot: " + s +
				" Dero\n\nPredictions: " + c + "\nHigher Predictions: " + u + "\nLower Predictions: " + d +
				"\n\nPayout After: " + end_pay.String() + "\nRefund if not paid within " + rf + " minutes\n\nRounds Completed: " + r
		}

	} else {
		var live string
		if now > Predict.buffer {
			live = "\n\nAccepting " + p + " Predictions "
		} else {
			left := Predict.buffer - now
			live = "\n\n" + p + "\nBuffer ends in " + strconv.Itoa(int(left)) + " seconds"
		}

		if amt_f, err := strconv.ParseFloat(amt, 64); err == nil {
			amt = fmt.Sprintf("%.5f", amt_f)
		}

		node := ""
		if isOnChainPrediction(p) {
			node = "\n\nNode: " + f
		}

		if m == "0" {
			pw := strconv.Itoa(ta / 60)
			info = "SCID:\n\n" + Predict.Contract.SCID + live + node +
				"\n\nCloses at: " + utc + "\nMark posted with in " + pw + " minutes of close\n\nPrediction Amount: " + amt + " Dero\nRound Pot: " + s + " Dero\n\nPredictions: " + c +
				"\nHigher Predictions: " + u + "\nLower Predictions: " + d + "\n\nPayout After: " + end_pay.String() + "\nRefund if not paid within " + rf + " minutes\n\nRounds Completed: " + r
		} else {
			info = "SCID:\n\n" + Predict.Contract.SCID + live + node +
				"\n\nCloses at: " + utc + "\nMark: " + m + "\n\nPrediction Amount: " + amt + " Dero\nRound Pot: " + s + " Dero\n\nPredictions: " + c +
				"\nHigher Predictions: " + u + "\nLower Predictions: " + d + "\n\nPayout After: " + end_pay.String() + "\nRefund if not paid within " + rf + " minutes\n\nRounds Completed: " + r
		}
	}

	return
}

// Format dPrediction end result text
//   - fr is the un-split result string
//   - m is prediction mark
func roundResults(fr, m string) string {
	if len(Predict.Contract.SCID) == 64 && fr != "" {
		split := strings.Split(fr, "_")
		var res string
		var def string
		var x float64

		if isOnChainPrediction(split[0]) {
			switch onChainPrediction(split[0]) {
			case 1:
				x = 1
			case 2:
				x = 100000
			case 3:
				x = 1
			default:
				x = 1
			}
		} else {
			if menu.CoinDecimal(split[0]) == 8 {
				x = 100000000
			} else {
				x = 100
			}
		}

		if mark, err := strconv.ParseFloat(m, 64); err == nil {
			if rpc.StringToInt(split[1]) > int(mark*x) {
				res = "Higher "
				def = " > "
			} else if rpc.StringToInt(split[1]) == int(mark*x) {
				res = "Equal "
				def = " == "
			} else {
				res = "Lower "
				def = " < "
			}
		}

		if final, err := strconv.ParseFloat(split[1], 64); err == nil {
			var fStr string
			if isOnChainPrediction(split[0]) {
				switch onChainPrediction(split[0]) {
				case 1:
					fStr = fmt.Sprintf("%.0f", final/x)
				case 2:
					fStr = fmt.Sprintf("%.5f", final/x)
				case 3:
					fStr = fmt.Sprintf("%.0f", final/x)
				default:

				}
			} else {
				if menu.CoinDecimal(split[0]) == 8 {
					fStr = fmt.Sprintf("%.8f", final/x)
				} else {
					fStr = fmt.Sprintf("%.2f", final/x)
				}
			}

			return split[0] + " " + res + fStr + def + m
		}

	}
	return ""
}

// Formats non-initialized dPrediction info string
//   - fr is the un-split result string
//   - tx is the previous payout TXID
//   - r is total completed prediction rounds
//   - m is prediction mark
func P_no_initResults(fr, tx, r, m string) (info string) {
	info = "SCID:\n\n" + Predict.Contract.SCID + "\n" + "\nRound Completed\n\nRound Mark: " + m +
		"\nRound Results: " + roundResults(fr, m) + "\n\nPayout TXID: " + tx + "\n\nRounds Completed: " + r

	Predict.prediction = ""

	return
}

// Populate all dReams dPrediction contracts
//   - Pass contracts from db store, can be nil arg
func PopulatePredictions(contracts map[string]string) {
	if rpc.Daemon.IsConnected() && gnomon.IsReady() {
		list := []string{}
		owned := []string{}
		if contracts == nil {
			contracts = gnomon.GetAllOwnersAndSCIDs()
		}

		for sc := range contracts {
			list, owned = checkBetContract(sc, "p", list, owned)
		}

		t := len(list)
		list = append(list, " Contracts: "+strconv.Itoa(t))
		sort.Strings(list)
		Predict.Public.SCIDs = list

		sort.Strings(owned)
		Predict.Owned.SCIDs = owned

	}
}

// Check if dPrediction is live on SCID
func CheckActivePrediction(scid string) bool {
	if len(scid) == 64 && gnomon.IsReady() {
		_, ends := gnomon.GetSCIDValuesByKey(scid, "p_end_at")
		_, buff := gnomon.GetSCIDValuesByKey(scid, "buffer")
		if ends != nil && buff != nil {
			now := time.Now().Unix()
			if now < int64(ends[0]) && now > int64(buff[0]) {
				return true
			}
		}
	}
	return false
}

// Gets dPrediction data from SCID and return formatted info string
func GetPrediction(scid string) (info string) {
	if rpc.Daemon.IsConnected() && gnomon.IsReady() {
		predicting, _ := gnomon.GetSCIDValuesByKey(scid, "predicting")
		url, _ := gnomon.GetSCIDValuesByKey(scid, "p_url")
		final, _ := gnomon.GetSCIDValuesByKey(scid, "p_final")
		//final_tx, _ := gnomon.GetSCIDValuesByKey(scid, "p_final_txid")
		_, amt := gnomon.GetSCIDValuesByKey(scid, "p_amount")
		_, init := gnomon.GetSCIDValuesByKey(scid, "p_init")
		_, up := gnomon.GetSCIDValuesByKey(scid, "p_up")
		_, down := gnomon.GetSCIDValuesByKey(scid, "p_down")
		_, count := gnomon.GetSCIDValuesByKey(scid, "p_#")
		_, end := gnomon.GetSCIDValuesByKey(scid, "p_end_at")
		_, buffer := gnomon.GetSCIDValuesByKey(scid, "buffer")
		_, pot := gnomon.GetSCIDValuesByKey(scid, "p_total")
		_, rounds := gnomon.GetSCIDValuesByKey(scid, "p_played")
		_, mark := gnomon.GetSCIDValuesByKey(scid, "mark")
		_, time_a := gnomon.GetSCIDValuesByKey(scid, "time_a")
		_, time_b := gnomon.GetSCIDValuesByKey(scid, "time_b")
		_, time_c := gnomon.GetSCIDValuesByKey(scid, "time_c")

		var pre, p_played, p_final, p_mark string
		if init != nil {
			if init[0] == 1 {
				Predict.init = true

				if buffer != nil {
					Predict.buffer = int64(buffer[0])
				}

				Predict.amount = amt[0]
				if predicting != nil {
					pre = predicting[0]
				}
				Predict.prediction = pre

				p_amt := fmt.Sprint(float64(Predict.amount) / 100000)

				p_down := fmt.Sprint(down[0])
				p_up := fmt.Sprint(up[0])
				p_count := fmt.Sprint(count[0])

				p_pot := fmt.Sprint(pot[0])
				p_played := fmt.Sprint(rounds[0])

				var p_feed string
				if url != nil {
					Predict.feed = url[0]
					p_feed = url[0]
				}

				if mark != nil {
					if predicting != nil {
						if isOnChainPrediction(predicting[0]) {
							i := onChainPrediction(predicting[0])
							switch i {
							case 1:
								p_mark = fmt.Sprintf("%d", mark[0])
							case 2:
								p_mark = fmt.Sprintf("%.5f", float64(mark[0])/100000)
							case 3:
								p_mark = fmt.Sprintf("%d", mark[0])
							}
						} else {
							if menu.CoinDecimal(pre) == 8 {
								p_mark = fmt.Sprintf("%.8f", float64(mark[0])/100000000)
							} else {
								p_mark = fmt.Sprintf("%.2f", float64(mark[0])/100)
							}
						}
					}
				} else {
					p_mark = "0"
				}

				var p_end string
				if init[0] == 1 {
					end_at := uint(end[0])
					p_end = fmt.Sprint(end_at * 1000)
				}

				info = P_initResults(pre, p_amt, p_end, p_count, p_pot, p_up, p_down, p_played, p_feed, p_mark, int(time_a[0]), int(time_b[0]), int(time_c[0]))

			} else {
				if final != nil {
					p_final = final[0]
				}

				Predict.init = false
				txid := FetchPredictionFinal(scid)

				if mark != nil {
					split := strings.Split(p_final, "_")
					if isOnChainPrediction(split[0]) {
						i := onChainPrediction(split[0])

						switch i {
						case 1:
							p_mark = fmt.Sprintf("%d", mark[0])
						case 2:
							p_mark = fmt.Sprintf("%.5f", float64(mark[0])/100000)
						case 3:
							p_mark = fmt.Sprintf("%d", mark[0])
						}
					} else {
						if menu.CoinDecimal(split[0]) == 8 {
							p_mark = fmt.Sprintf("%.8f", float64(mark[0])/100000000)
						} else {
							p_mark = fmt.Sprintf("%.2f", float64(mark[0])/100)
						}
					}

				} else {
					p_mark = "0"
				}

				if rounds != nil {
					p_played = fmt.Sprint(rounds[0])
				}

				info = P_no_initResults(p_final, txid, p_played, p_mark)
			}
		}
	}

	return
}
