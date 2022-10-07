package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/andelf/go-curl"
	"github.com/deroproject/derohe/rpc"
)

type priceFeed struct {
	Success      bool   `json:"success"`
	Initialprice string `json:"initialprice"`
	Price        string `json:"price"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Volume       string `json:"volume"`
	Bid          string `json:"bid"`
	Ask          string `json:"ask"`
}

const (
	pre  = "http://"
	suff = "/json_rpc"
)

func whichDaemon(s string) { /// select menu changes dameon address
	if gnomonInit {
		defaultIndexer.Close()
		gnomonInit = false
	}
	switch s {
	case "REMOTE":
		daemonAddress = pre + DAEMON_REMOTE_DEFAULT_A + suff
		gnomon(DAEMON_REMOTE_DEFAULT_A)
	case "MAINNET":
		daemonAddress = pre + DAEMON_MAINNET_DEFAULT + suff
		gnomon(DAEMON_MAINNET_DEFAULT)
	case "SIMULATOR":
		daemonAddress = pre + DAEMON_SIMULATOR_DEFAULT + suff
		gnomon(DAEMON_SIMULATOR_DEFAULT)
	case "CUSTOM":
		customPopUp() /// enter custom address
	default:
		daemonAddress = pre + DAEMON_REMOTE_DEFAULT_A + suff
		gnomon(DAEMON_REMOTE_DEFAULT_A)
	}
}

func isDaemonConnected() { /// check if daemon is connected
	if daemonConnectBool {
		if !daemonCheckBox.Checked {
			log.Println("Daemon Connected")
		}
		daemonCheckBox.SetChecked(true)
	} else {
		log.Println("Daemon Not Connected")
		currentHeight.SetText("Height:")
		defaultIndexer.Close()
		if daemonCheckBox.Checked {
			daemonCheckBox.SetChecked(false)
		}
	}
}

func isWalletConnected() { /// check if wallet is connected
	if walletConnectBool {
		if !walletCheckBox.Checked {
			log.Println("Wallet Connected")
			walletCheckBox.SetChecked(true)
		}
		GetBalance()
		enablePrediction(p_initialized)
		removeButton.Enable()
	} else {
		log.Println("Wallet Not Connected")
		enablePrediction(false)
		enableSports(false)
		if walletCheckBox.Checked {
			walletCheckBox.SetChecked(false)
			walletConnectBool = false
			removeButton.Disable()
		}
	}

	if walletCheckBox.Checked { /// if wallet is connected and any changes to inputs, show disconnected
		checkPass()
		if pre+rpcWalletInput.Text+suff != walletAddress {
			walletBalance.SetText("Balance: ")
			walletAddress = ""
			walletCheckBox.SetChecked(false)
			walletConnectBool = false
		}
	}
}

func checkPass() { /// check if user:pass has changed
	data := []byte(rpcLoginInput.Text)
	hash := sha256.Sum256(data)

	if hash != passHash {
		walletBalance.SetText("Balance: ")
		walletCheckBox.SetChecked(false)
		walletConnectBool = false
	}
}

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func convertString_Int(s string) int {
	intValue, int_err := strconv.Atoi(s)

	if int_err != nil {
		log.Println("Could Not Convert String to Int", int_err)
		return 0
	} else {
		return intValue
	}

}

func trimTeamA(s string) string {
	for i := range s {
		if i > 3 {
			return s[0:3]
		}
	}
	return s[:0]
}

func trimTeamB(s string) string {
	for i := range s {
		if i > 3 {
			return s[i:]
		}
	}
	return s[:0]
}

func teamReturn(t string) string {
	var team string
	switch t {
	case "0":
		team = "team_a"
	case "1":
		team = "team_b"
	default:
		team = "none"

	}

	return team

}

func enablePrediction(init bool) {
	if init {
		upButton.Enable()
		downButton.Enable()

	} else {
		upButton.Disable()
		downButton.Disable()
	}
}

func enableSports(init bool) {
	if init {
		a_button.Enable()
		b_button.Enable()
		game_select.Enable()
		s_multi.Enable()

	} else {
		a_button.Disable()
		b_button.Disable()
		game_select.Disable()
		s_multi.Disable()
	}
}

func configExists() bool {
	if _, err := os.Stat("config.json"); err == nil {
		return true
	} else {
		return false
	}
}

func writeConfig(u rpc.Argument) {
	file, err := os.Create("config.json")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	jsonByte, _ := json.MarshalIndent(u, "", " ")

	_, err2 := file.Write(jsonByte)

	if err2 != nil {
		log.Fatal("Error writing config file: ", err2)
	}
}

func readConfig() interface{} {
	file, err := os.ReadFile("config.json")

	if err != nil {
		log.Fatal("Error reading config file: ", err)
		return nil
	}

	var payload rpc.Argument
	err2 := json.Unmarshal(file, &payload)
	if err2 != nil {
		log.Fatal("Error during unmarshal: ", err)
		return nil
	}

	return payload.Value
}

func getPrice(coin string) string {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	var url string
	var found priceFeed
	switch coin {
	case "BTC":
		url = "https://tradeogre.com/api/v1/ticker/usdt-btc"
	case "DERO":
		url = "https://tradeogre.com/api/v1/ticker/usdt-dero"
	case "XMR":
		url = "https://tradeogre.com/api/v1/ticker/usdt-xmr"
	default:
		return ""
	}
	easy.Setopt(curl.OPT_URL, url)

	curlBtc := func(buf []byte, userdata interface{}) bool {
		err := json.Unmarshal([]byte(buf), &found)

		if err != nil {
			fmt.Println(err)
		}
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, curlBtc)

	if err := easy.Perform(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	return found.Price
}
