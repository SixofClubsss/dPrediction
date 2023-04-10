package prediction

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/SixofClubsss/dReams/holdero"
	"github.com/SixofClubsss/dReams/menu"
	"github.com/SixofClubsss/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type predictObjects struct {
	Contract    string
	Leaders_map map[string]uint64
	//Leaders_display []string
	Info          *widget.Label
	Prices        *widget.Label
	Predict_list  *widget.List
	Favorite_list *widget.List
	Owned_list    *widget.List
	Leaders_list  *widget.List
	Remove_button *widget.Button

	Higher *widget.Button
	Lower  *widget.Button
	// prediction leaderboard
	//Change         *widget.Button
	//Remove         *widget.Button
	//NameEntry      *widget.Entry
	Prediction_box *fyne.Container
	//P_contract     *widget.SelectEntry
}

var Predict predictObjects

// Disable dPrediction objects
func DisablePreditions(d bool) {
	if d {
		Predict.Prediction_box.Hide()
	} else {
		Predict.Prediction_box.Show()
	}
	Predict.Prediction_box.Refresh()
}

// Disable dSports objects
func DisableSports(d bool) {
	if d {
		Sports.Sports_box.Hide()
		menu.Control.Sports_check.SetChecked(false)
	}

	Sports.Sports_box.Refresh()
}

func PredictConnectedBox() fyne.Widget {
	menu.Control.Predict_check = widget.NewCheck("", func(b bool) {
		if !b {
			// prediction leaderboard
			// holdero.Actions.NameEntry.Hide()
			// holdero.Actions.Change.Hide()
			//Predict.Leaders_display = []string{}
			Predict.Higher.Hide()
			Predict.Lower.Hide()

		}
	})
	menu.Control.Predict_check.Disable()

	return menu.Control.Predict_check
}

func PreictionContractEntry() fyne.Widget {
	options := []string{}
	menu.Control.P_contract = widget.NewSelectEntry(options)
	menu.Control.P_contract.PlaceHolder = "Contract Address: "
	menu.Control.P_contract.OnCursorChanged = func() {
		if rpc.Daemon.Connect {
			go func() {
				if len(Predict.Contract) == 64 {
					yes := rpc.ValidBetContract(Predict.Contract)
					if yes {
						menu.Control.Predict_check.SetChecked(true)
					} else {
						menu.Control.Predict_check.SetChecked(false)
					}
				} else {
					menu.Control.Predict_check.SetChecked(false)
				}
			}()
		}
	}

	this := binding.BindString(&Predict.Contract)
	menu.Control.P_contract.Bind(this)

	return menu.Control.P_contract
}

// prediction leaderboard
// func LeadersDisplay() fyne.Widget {
// 	Predict.Leaders_display = []string{}
// 	Predict.Leaders_list = widget.NewList(
// 		func() int {
// 			return len(Predict.Leaders_display)
// 		},
// 		func() fyne.CanvasObject {
// 			return widget.NewLabel("")
// 		},
// 		func(i widget.ListItemID, o fyne.CanvasObject) {
// 			o.(*widget.Label).SetText(Predict.Leaders_display[i])
// 		})
//
// 	return Predict.Leaders_list
// }

func ShowPredictionControls() {
	DisablePreditions(false)
	Predict.Higher.Show()
	Predict.Lower.Show()
	// prediction leaderboard
	// holdero.Actions.NameEntry.Show()
	// holdero.Actions.NameEntry.Text = menu.CheckPredictionName(Predict.Contract)
	// holdero.Actions.NameEntry.Refresh()
}

func setPredictionControls(str string) (item string) {
	split := strings.Split(str, "   ")
	if len(split) >= 3 {
		trimmed := strings.Trim(split[2], " ")
		if len(trimmed) == 64 {
			item = str
			menu.Control.P_contract.SetText(trimmed)
			go SetPredictionInfo(trimmed)
			if menu.CheckActivePrediction(trimmed) {
				ShowPredictionControls()
			} else {
				DisablePreditions(true)
			}
		}
	}

	return
}

func SetPredictionInfo(scid string) {
	info := GetPrediction(rpc.Daemon.Connect, scid)
	if info != "" {
		Predict.Info.SetText(info)
		Predict.Info.Refresh()
	}
}

func SetPredictionPrices(d bool) {
	if d {
		_, btc := holdero.GetPrice("BTC-USDT")
		_, dero := holdero.GetPrice("DERO-USDT")
		_, xmr := holdero.GetPrice("XMR-USDT")
		/// custom feed with rpc.Display.P_feed
		prices := "Current Price feed from dReams Client\nBTC: " + btc + "\nDERO: " + dero + "\nXMR: " + xmr

		Predict.Prices.SetText(prices)
	}
}

// List object for populating public dPrediction contracts
func PredictionListings(tab *container.AppTabs) fyne.CanvasObject {
	Predict.Predict_list = widget.NewList(
		func() int {
			return len(menu.Control.Predict_contracts)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(canvas.NewImageFromImage(nil), widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(menu.Control.Predict_contracts[i])
			if menu.Control.Predict_contracts[i][0:2] != "  " {
				var key string
				split := strings.Split(menu.Control.Predict_contracts[i], "   ")
				if len(split) >= 3 {
					trimmed := strings.Trim(split[2], " ")
					if len(trimmed) == 64 {
						key = trimmed
					}
				}

				badge := canvas.NewImageFromResource(menu.DisplayRating(menu.Control.Contract_rating[key]))
				badge.SetMinSize(fyne.NewSize(35, 35))
				o.(*fyne.Container).Objects[0] = badge
			}
		})

	var item string

	Predict.Predict_list.OnSelected = func(id widget.ListItemID) {
		if id != 0 && menu.Connected() {
			item = setPredictionControls(menu.Control.Predict_contracts[id])
			Predict.Favorite_list.UnselectAll()
			Predict.Owned_list.UnselectAll()
		} else {
			DisablePreditions(true)
		}
	}

	save := widget.NewButton("Favorite", func() {
		menu.Control.Predict_favorites = append(menu.Control.Predict_favorites, item)
		sort.Strings(menu.Control.Predict_favorites)
	})

	rate := widget.NewButton("Rate", func() {
		if len(Predict.Contract) == 64 {
			if !menu.CheckOwner(Predict.Contract) {
				reset := tab.Selected().Content
				tab.Selected().Content = menu.RateConfirm(Predict.Contract, tab, reset)
				tab.Selected().Content.Refresh()
			} else {
				log.Println("[dReams] You own this contract")
			}
		}
	})

	cont := container.NewBorder(
		nil,
		container.NewBorder(nil, nil, save, rate, layout.NewSpacer()),
		nil,
		nil,
		Predict.Predict_list)

	return cont
}

// List object for populating favorite dPrediction contracts
func PredicitionFavorites() fyne.CanvasObject {
	Predict.Favorite_list = widget.NewList(
		func() int {
			return len(menu.Control.Predict_favorites)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(menu.Control.Predict_favorites[i])
		})

	var item string

	Predict.Favorite_list.OnSelected = func(id widget.ListItemID) {
		if menu.Connected() {
			item = setPredictionControls(menu.Control.Predict_favorites[id])
			Predict.Predict_list.UnselectAll()
			Predict.Owned_list.UnselectAll()
		} else {
			DisablePreditions(true)
		}
	}

	remove := widget.NewButton("Remove", func() {
		if len(menu.Control.Predict_favorites) > 0 {
			Predict.Favorite_list.UnselectAll()
			new := menu.Control.Predict_favorites
			for i := range new {
				if new[i] == item {
					copy(new[i:], new[i+1:])
					new[len(new)-1] = ""
					new = new[:len(new)-1]
					menu.Control.Predict_favorites = new
					break
				}
			}
		}
		Predict.Favorite_list.Refresh()
		sort.Strings(menu.Control.Predict_favorites)
	})

	cont := container.NewBorder(
		nil,
		container.NewBorder(nil, nil, nil, remove, layout.NewSpacer()),
		nil,
		nil,
		Predict.Favorite_list)

	return cont
}

// List object for populating owned dPrediction contracts
func PredictionOwned() fyne.CanvasObject {
	Predict.Owned_list = widget.NewList(
		func() int {
			return len(menu.Control.Predict_owned)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(menu.Control.Predict_owned[i])
		})

	Predict.Owned_list.OnSelected = func(id widget.ListItemID) {
		if menu.Connected() {
			setPredictionControls(menu.Control.Predict_owned[id])
			Predict.Predict_list.UnselectAll()
			Predict.Favorite_list.UnselectAll()
		} else {
			DisablePreditions(true)
		}

	}

	return Predict.Owned_list
}

// prediction leaderboard
// func Remove() fyne.Widget {
// 	Predict.Remove_button = widget.NewButton("Remove", func() {
// 		namePopUp(2)
// 	})
//
// 	Predict.Remove_button.Hide()
//
// 	return Predict.Remove_button
// }

func P_initResults(p, amt, eA, c, to, u, d, r, f, m string, ta, tb, tc int) (info string) { /// prediction info, initialized
	end_time, _ := rpc.MsToTime(eA)
	utc := end_time.String()
	add := rpc.StringToInt(eA)
	end := strconv.Itoa(add + (tc * 1000))
	end_pay, _ := rpc.MsToTime(end)
	rf := strconv.Itoa(tb / 60)

	result, err := strconv.ParseFloat(to, 32)

	if err != nil {
		log.Println("[Predictions]", err)
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
			info = "SCID: \n" + Predict.Contract + "\n\n" + p + wfp + "\n\nNode: " + f + "\n\nMark: " + mark + "\nRound Pot: " + s +
				"\n\nPredictions: " + c + "\nHigher Predictions: " + u + "\nLower Predictions: " + d +
				"\n\nPayout After: " + end_pay.String() + "\nRefund if not paid within " + rf + " minutes\n\nRounds Completed: " + r
		} else {
			info = "SCID: \n" + Predict.Contract + "\n\n" + p + wfp + "\n\nMark: " + mark + "\nRound Pot: " + s +
				"\n\nPredictions: " + c + "\nHigher Predictions: " + u + "\nLower Predictions: " + d +
				"\n\nPayout After: " + end_pay.String() + "\nRefund if not paid within " + rf + " minutes\n\nRounds Completed: " + r
		}

	} else {
		var live string
		if now > rpc.Predict.Buffer {
			live = "\n\nAccepting " + p + " Predictions "
		} else {
			left := rpc.Predict.Buffer - now
			live = "\n\n" + p + "\nBuffer ends in " + strconv.Itoa(int(left)) + " seconds"
		}

		node := ""
		if isOnChainPrediction(p) {
			node = "\n\nNode: " + f
		}

		if m == "0" {
			pw := strconv.Itoa(ta / 60)
			info = "SCID: \n" + Predict.Contract + live + node +
				"\n\nCloses at: " + utc + "\nMark posted with in " + pw + " minutes of close\n\nPrediction Amount: " + amt + " Dero\nRound Pot: " + s + " \n\nPredictions: " + c +
				"\nHigher Predictions: " + u + "\nLower Predictions: " + d + "\n\nPayout After: " + end_pay.String() + "\nRefund if not paid within " + rf + " minutes\n\nRounds Completed: " + r
		} else {
			info = "SCID: \n" + Predict.Contract + live + node +
				"\n\nCloses at: " + utc + "\nMark: " + m + "\n\nPrediction Amount: " + amt + " Dero\nRound Pot: " + s + "\n\nPredictions: " + c +
				"\nHigher Predictions: " + u + "\nLower Predictions: " + d + "\n\nPayout After: " + end_pay.String() + "\nRefund if not paid within " + rf + " minutes\n\nRounds Completed: " + r
		}
	}

	return
}

func roundResults(fr, m string) string { /// prediction results text
	if len(Predict.Contract) == 64 && fr != "" {
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
			if holdero.CoinDecimal(split[0]) == 8 {
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
				if holdero.CoinDecimal(split[0]) == 8 {
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

func P_no_initResults(fr, tx, r, m string) (info string) { /// prediction info, not initialized
	info = "SCID: \n" + Predict.Contract + "\n" + "\nRound Completed\n\nRound Mark: " + m +
		"\nRound Results: " + roundResults(fr, m) + "\n\nPayout TXID: " + tx + "\n\nRounds Completed: " + r

	rpc.Display.Prediction = ""

	return
}

func GetPrediction(d bool, scid string) (info string) {
	if d && menu.Gnomes.Init && !menu.GnomonClosing() && menu.Gnomes.Sync {
		predicting, _ := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "predicting", menu.Gnomes.Indexer.ChainHeight, true)
		url, _ := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_url", menu.Gnomes.Indexer.ChainHeight, true)
		final, _ := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_final", menu.Gnomes.Indexer.ChainHeight, true)
		//final_tx, _ := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_final_txid", menu.Gnomes.Indexer.ChainHeight, true)
		_, amt := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_amount", menu.Gnomes.Indexer.ChainHeight, true)
		_, init := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_init", menu.Gnomes.Indexer.ChainHeight, true)
		_, up := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_up", menu.Gnomes.Indexer.ChainHeight, true)
		_, down := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_down", menu.Gnomes.Indexer.ChainHeight, true)
		_, count := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_#", menu.Gnomes.Indexer.ChainHeight, true)
		_, end := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_end_at", menu.Gnomes.Indexer.ChainHeight, true)
		_, buffer := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "buffer", menu.Gnomes.Indexer.ChainHeight, true)
		_, pot := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_total", menu.Gnomes.Indexer.ChainHeight, true)
		_, rounds := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "p_played", menu.Gnomes.Indexer.ChainHeight, true)
		_, mark := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "mark", menu.Gnomes.Indexer.ChainHeight, true)
		_, time_a := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "time_a", menu.Gnomes.Indexer.ChainHeight, true)
		_, time_b := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "time_b", menu.Gnomes.Indexer.ChainHeight, true)
		_, time_c := menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, "time_c", menu.Gnomes.Indexer.ChainHeight, true)

		var pre, p_played, p_final, p_mark string
		if init != nil {
			if init[0] == 1 {
				rpc.Predict.Init = true

				if buffer != nil {
					rpc.Predict.Buffer = int64(buffer[0])
				}

				rpc.Predict.Amount = amt[0]
				if predicting != nil {
					pre = predicting[0]
				}
				rpc.Display.Prediction = pre

				p_amt := fmt.Sprint(float64(rpc.Predict.Amount) / 100000)

				p_down := fmt.Sprint(down[0])
				p_up := fmt.Sprint(up[0])
				p_count := fmt.Sprint(count[0])

				p_pot := fmt.Sprint(pot[0])
				p_played := fmt.Sprint(rounds[0])

				var p_feed string
				if url != nil {
					rpc.Display.P_feed = url[0]
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
							if holdero.CoinDecimal(pre) == 8 {
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

				rpc.Predict.Init = false
				txid := rpc.FetchPredictionFinal(scid)

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
						if holdero.CoinDecimal(split[0]) == 8 {
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
