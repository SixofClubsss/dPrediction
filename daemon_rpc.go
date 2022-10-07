package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/ybbus/jsonrpc/v3"
)

const (
	DAEMON_REMOTE_DEFAULT_A  = "dero-node.mysrv.cloud:10102"
	DAEMON_REMOTE_DEFAULT_B  = "derostats.io:10102"
	DAEMON_MAINNET_DEFAULT   = "127.0.0.1:10102"
	DAEMON_SIMULATOR_DEFAULT = "127.0.0.1:20000"
	P_SC_ID                  = "9af71c91a2569a03c91a6a20ba2d4970a06ab83b4caddbf9e4e33b47731169bf"
	S_SC_ID                  = "ab4e5ec853e22bf5b4eec44988c6e7fa6285cd109b50c222ed911f337e483fdb"
)

var daemonConnectBool bool
var amount int
var p_initialized bool
var leaders map[string]uint64

func ping() error { /// ping blockchain for connection
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientD := jsonrpc.NewClient(daemonAddress)
	var result string
	err := rpcClientD.CallFor(ctx, &result, "DERO.Ping")
	if err != nil {
		daemonConnectBool = false
		return nil
	}

	if result == "Pong " {
		daemonConnectBool = true
	} else {
		daemonConnectBool = false
	}

	return err
}

func getHeight() error { /// get current height and displays
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientD := jsonrpc.NewClient(daemonAddress)
	var result *rpc.Daemon_GetHeight_Result
	err := rpcClientD.CallFor(ctx, &result, "DERO.GetHeight")

	if err != nil {
		return nil
	}
	h := result.Height
	log.Printf("Daemon Height: %d \n", h)
	str := strconv.FormatUint(result.Height, 10)
	currentHeight.SetText("Height: " + str)

	return err

}

func p_getSC() error { /// search sc using getsc method
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	rpcClientD := jsonrpc.NewClient(daemonAddress)
	params := &rpc.GetSC_Params{
		SCID:      P_SC_ID,
		Code:      false,
		Variables: true,
	}

	var result *rpc.GetSC_Result
	err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", params)

	if err != nil {
		log.Println(err)
		return nil
	}

	p_predict := result.VariableStringKeys["predicting"]
	p_amt := result.VariableStringKeys["p_amount"]
	amount = int(p_amt.(float64))
	p_init := result.VariableStringKeys["p_init"]
	p_up := result.VariableStringKeys["p_up"]
	p_down := result.VariableStringKeys["p_down"]
	p_count := result.VariableStringKeys["p_#"]
	p_endAt := result.VariableStringKeys["p_end_at"]
	p_pot := result.VariableStringKeys["p_total"]
	p_rounds := result.VariableStringKeys["p_played"]
	p_feed_url := result.VariableStringKeys["p_url"]
	p_final := result.VariableStringKeys["p_final"]
	p_txid := result.VariableStringKeys["p_final_txid"]
	p_mark := result.VariableStringKeys["mark"]

	pre := fmt.Sprint(p_predict)
	amt := fmt.Sprint(float64(amount) / 100000)
	in := fmt.Sprint(p_init)
	d := fmt.Sprint(p_down)
	u := fmt.Sprint(p_up)
	c := fmt.Sprint(p_count)

	t := fmt.Sprint(p_pot)
	r := fmt.Sprint(p_rounds)
	feed := fmt.Sprint(p_feed_url)
	final := fmt.Sprint(p_final)
	txid := fmt.Sprint(p_txid)
	mark := fmt.Sprint(p_mark)

	if in == "1" {
		p_initialized = true
		p, err := hex.DecodeString(pre)
		if err != nil {
			log.Println(err, "prediction value")
			return nil
		}

		f, err2 := hex.DecodeString(feed)
		if err2 != nil {
			log.Println(err2, "feed url")
			return nil
		}

		end_at := uint(p_endAt.(float64))
		eA := fmt.Sprint(end_at * 1000)
		if p_mark != nil {
			p_initResults(string(p), amt, eA, c, t, u, d, r, string(f), mark, true) /// predict display
		} else {
			p_initResults(string(p), amt, eA, c, t, u, d, r, string(f), "", false)
		}

	} else {
		p_initialized = false
		f, err := hex.DecodeString(final)
		if err != nil {
			log.Println(err, "final amount")
			return nil
		}
		p_no_initResults(string(f), txid, r, mark)

	}

	return err
}

func checkNames(name string) error { /// confirms if name is on leaderboard
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	rpcClientD := jsonrpc.NewClient(daemonAddress)
	params := &rpc.GetSC_Params{
		SCID:      P_SC_ID,
		Code:      false,
		Variables: true,
	}

	var result *rpc.GetSC_Result
	err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", params)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	if result.VariableStringKeys["u_"+name] != nil {
		n := result.VariableStringKeys["u_"+name]
		leaders[name] = uint64(n.(float64))
	}

	return err
}
