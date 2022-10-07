package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/civilware/Gnomon/indexer"
	"github.com/civilware/Gnomon/storage"
	"github.com/civilware/Gnomon/structures"
	"github.com/deroproject/derohe/rpc"
)

var gnomonInit bool
var defaultIndexer *indexer.Indexer
var Graviton_backend *storage.GravitonStore

func gnomon(ep string) {
	log.Println("Starting gnomes.")
	shasum := fmt.Sprintf("%x", sha1.Sum([]byte("gnomon")))
	db_folder := fmt.Sprintf("gnomondb\\%s_%s", "GNOMON", shasum)
	Graviton_backend = storage.NewGravDB(db_folder, "25ms")

	search_filter := ""
	last_indexedheight := Graviton_backend.GetLastIndexHeight()
	daemon_endpoint := ep
	runmode := "daemon"
	mbl := false
	closeondisconnect := false
	fastsync := false

	defaultIndexer = indexer.NewIndexer(Graviton_backend, search_filter, last_indexedheight, daemon_endpoint, runmode, mbl, closeondisconnect, fastsync)
	go defaultIndexer.StartDaemonMode()
	time.Sleep(3 * time.Second)

	gnomonInit = true

}

// func checkOwners(check string) {
// 	dev, _ := Graviton_backend.GetSCIDValuesByKey(P_SC_ID, "dev", defaultIndexer.ChainHeight, true)
// 	owner1, _ := Graviton_backend.GetSCIDValuesByKey(P_SC_ID, "owner1", defaultIndexer.ChainHeight, true)

// 	if check == dev[0] || check == owner1[0] {
// 		owner_button.Enable()
// 	} else {
// 		owner_button.Disable()
// 	}

// 	fmt.Println(dev[0])
// 	fmt.Println(owner1[0])

// }

func ifPlaying() {
	nameInput.SetText("")
	if _, err := os.Stat("config.json"); err == nil {
		log.Printf("config.json exists\n")
		str := fmt.Sprintf("%v", readConfig())
		nameInput.SetText(str)
	} else {
		log.Printf("config.json does not exist\n")
		var s []*structures.SCIDVariable
		check := Graviton_backend.GetAllSCIDInvokeDetailsByEntrypoint(P_SC_ID, "Predict")
		for i := range check {
			str := "u_" + check[i].Sc_args.Value("name", rpc.DataString).(string)
			key, _ := defaultIndexer.GetSCIDKeysByValue(s, P_SC_ID, str, defaultIndexer.ChainHeight)
			signer := check[i].Sender
			if key != nil && signer == wallet {
				nameInput.SetText(check[i].Sc_args.Value("name", rpc.DataString).(string))
			}

		}
	}

}

func getBook() {

	initValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "s_init", defaultIndexer.ChainHeight, true)
	playedValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "s_played", defaultIndexer.ChainHeight, true)
	//hlValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "hl", defaultIndexer.ChainHeight, true)

	init := convertString_Int(initValue[0])
	played := convertString_Int(playedValue[0])

	game_select.SetText("")
	s_sc_displayT.SetText("SC ID: \n" + S_SC_ID + "\n\nGames Completed: " + playedValue[0] + "\nCurrent Games:\n")
	go func() {
		for {
			iv := strconv.Itoa(played)
			s_initValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "s_init_"+iv, defaultIndexer.ChainHeight, true)
			if s_initValue != nil {
				gameValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "game_"+iv, defaultIndexer.ChainHeight, true)
				leagueValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "league_"+iv, defaultIndexer.ChainHeight, true)
				s_nValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "s_#_"+iv, defaultIndexer.ChainHeight, true)
				s_amtValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "s_amount_"+iv, defaultIndexer.ChainHeight, true)
				s_endValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "s_end_at_"+iv, defaultIndexer.ChainHeight, true)
				s_totalValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "s_total_"+iv, defaultIndexer.ChainHeight, true)
				//s_urlValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "s_url_"+iv, defaultIndexer.ChainHeight, true)
				s_taValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "team_a_"+iv, defaultIndexer.ChainHeight, true)
				s_tbValue, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "team_b_"+iv, defaultIndexer.ChainHeight, true)

				team_a := trimTeamA(gameValue[0])
				team_b := trimTeamB(gameValue[0])

				var end_at float64
				if s, err := strconv.ParseFloat(s_endValue[0], 64); err == nil {
					end_at = s
				} else {
					log.Println("Foat Conversion Error", err)
				}
				end := uint(end_at)
				eA := fmt.Sprint(end * 1000)

				s_Results(gameValue[0], iv, leagueValue[0], s_amtValue[0], eA, s_nValue[0], team_a, team_b, s_taValue[0], s_tbValue[0], s_totalValue[0])

				if walletConnectBool {
					enableSports(true)
				}

			} else {
				enableSports(false)
			}

			if played >= init {
				break
			}

			played++
		}
	}()
}

func getSportsAmt(n string) int {
	amt, _ := Graviton_backend.GetSCIDValuesByKey(S_SC_ID, "s_amount_"+n, defaultIndexer.ChainHeight, true)
	if amt != nil {
		return convertString_Int(amt[0])
	} else {
		return 0
	}

}

func makeLeaderBoard() {
	leadersTotal = nil
	go func() {
		findLeaders := Graviton_backend.GetAllSCIDInvokeDetailsByEntrypoint(P_SC_ID, "Predict")

		for i := range findLeaders {
			str := findLeaders[i].Sc_args.Value("name", rpc.DataString).(string)
			checkNames(str)
		}

		printLeaders()
	}()

}
