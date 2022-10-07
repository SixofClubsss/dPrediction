package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/ybbus/jsonrpc/v3"
)

const (
	WALLET_MAINNET_DEFAULT   = "http://127.0.0.1:10103/json_rpc"
	WALLET_TESTNET_DEFAULT   = "http://127.0.0.1:40103/json_rpc"
	WALLET_SIMULATOR_DEFAULT = "http://127.0.0.1:30000/json_rpc"
	CHANGE_FEE               = 10000
	REMOVE_FEE               = 5000
)

var passHash [32]byte
var walletConnectBool bool
var wallet string

func getAddress() error { /// get address with auth
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})
	var result *rpc.GetAddress_Result
	err := rpcClientW.CallFor(ctx, &result, "GetAddress")

	if err != nil {
		walletConnectBool = false
		walletCheckBox.SetChecked(false)
		fmt.Println(err)
		return nil
	}

	addressLen := len(result.Address)
	if addressLen == 66 {
		walletConnectBool = true
		walletCheckBox.SetChecked(true)
		wallet = result.Address
		log.Println("Wallet Connected")
		fmt.Println("Dero Address: " + result.Address)
		data := []byte(rpcLoginInput.Text)
		passHash = sha256.Sum256(data)
		//checkOwners(wallet)
	}

	return err
}

func GetBalance() error { /// get wallet balance
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})
	var result *rpc.GetBalance_Result
	err := rpcClientW.CallFor(ctx, &result, "GetBalance")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	atomic := float64(result.Unlocked_Balance) /// unlocked balance in atomic units
	a := atomic / 100000
	str := strconv.FormatFloat(a, 'f', 5, 64)
	walletBalance.SetText("Balance: " + str)

	return err
}

func predictHigher() error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})

	a := uint64(amount)
	nameVal := nameInput.Text
	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Predict"}
	arg2 := rpc.Argument{Name: "pre", DataType: "U", Value: 1}
	arg3 := rpc.Argument{Name: "name", DataType: "S", Value: nameVal}
	args := rpc.Arguments{arg1, arg2, arg3}
	txid := rpc.Transfer_Result{}
	p := &rpc.SC_Invoke_Params{
		SC_ID:           P_SC_ID,
		SC_RPC:          args,
		SC_DERO_Deposit: a,
		Ringsize:        2,
	}
	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	if configExists() {
		fmt.Printf("Config exists\n")
	} else {
		writeConfig(arg3)
		nameInput.Disable()
		changeNameButton.Enable()
		fmt.Printf("Config does not exist\n")
	}

	log.Println("Prediction TX:", txid)

	return err
}

func predictLower() error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})

	a := uint64(amount)
	nameVal := nameInput.Text
	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Predict"}
	arg2 := rpc.Argument{Name: "pre", DataType: "U", Value: 0}
	arg3 := rpc.Argument{Name: "name", DataType: "S", Value: nameVal}
	args := rpc.Arguments{arg1, arg2, arg3}
	txid := rpc.Transfer_Result{}
	p := &rpc.SC_Invoke_Params{
		SC_ID:           P_SC_ID,
		SC_RPC:          args,
		SC_DERO_Deposit: a,
		Ringsize:        2,
	}
	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	if configExists() {
		fmt.Printf("Config exists\n")
	} else {
		writeConfig(arg3)
		nameInput.Disable()
		changeNameButton.Enable()
		fmt.Printf("Config does not exist\n")
	}
	log.Println("Prediction TX:", txid)

	return err
}

func nameChange() error { /// change leaderboard name
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})

	nameVal := nameInput.Text
	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "NameChange"}
	arg2 := rpc.Argument{Name: "name", DataType: "S", Value: nameVal}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}
	p := &rpc.SC_Invoke_Params{
		SC_ID:           P_SC_ID,
		SC_RPC:          args,
		SC_DERO_Deposit: CHANGE_FEE,
		Ringsize:        2,
	}
	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	writeConfig(arg2)
	nameInput.Disable()
	log.Println("Name Change TX:", txid)

	return err
}

func removeAddress() error { /// change leaderboard name
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})

	nameVal := nameInput.Text
	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Remove"}
	arg2 := rpc.Argument{Name: "name", DataType: "S", Value: nameVal}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}
	p := &rpc.SC_Invoke_Params{
		SC_ID:           P_SC_ID,
		SC_RPC:          args,
		SC_DERO_Deposit: REMOVE_FEE,
		Ringsize:        2,
	}
	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	writeConfig(arg2)
	nameInput.Disable()
	log.Println("Remove TX:", txid)

	return err
}

func pickTeam(pick int) error { /// pick sports team
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})
	n := game_select.Text
	a := uint64(getSportsAmt(n))
	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Book"}
	arg2 := rpc.Argument{Name: "n", DataType: "S", Value: n}
	arg3 := rpc.Argument{Name: "pre", DataType: "U", Value: pick}
	args := rpc.Arguments{arg1, arg2, arg3}
	txid := rpc.Transfer_Result{}
	p := &rpc.SC_Invoke_Params{
		SC_ID:           S_SC_ID,
		SC_RPC:          args,
		SC_DERO_Deposit: a,
		Ringsize:        2,
	}
	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	log.Println("Pick TX:", txid)

	return err
}

func setSports(end, amt int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "S_start"}
	arg2 := rpc.Argument{Name: "end", DataType: "U", Value: end}
	arg3 := rpc.Argument{Name: "amt", DataType: "U", Value: amt}
	arg4 := rpc.Argument{Name: "league", DataType: "S", Value: s_leagueInput.Text}
	arg5 := rpc.Argument{Name: "game", DataType: "S", Value: s_gameInput.Text}
	arg6 := rpc.Argument{Name: "feed", DataType: "S", Value: s_feedInput.Text}
	args := rpc.Arguments{arg1, arg2, arg3, arg4, arg5, arg6}
	txid := rpc.Transfer_Result{}
	params := &rpc.SC_Invoke_Params{
		SC_ID:           S_SC_ID,
		SC_RPC:          args,
		SC_DERO_Deposit: 0, /// TODO add inital deposit
		Ringsize:        2,
	}
	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	log.Println("Set Sports TX:", txid)

	return err
}

func setPrediction(end, amt int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "P_start"}
	arg2 := rpc.Argument{Name: "end", DataType: "U", Value: end}
	arg3 := rpc.Argument{Name: "amt", DataType: "U", Value: amt}
	arg4 := rpc.Argument{Name: "predict", DataType: "S", Value: p_nameInput.Text}
	arg5 := rpc.Argument{Name: "feed", DataType: "S", Value: p_feedInput.Text}
	args := rpc.Arguments{arg1, arg2, arg3, arg4, arg5}
	txid := rpc.Transfer_Result{}
	p := &rpc.SC_Invoke_Params{
		SC_ID:           P_SC_ID,
		SC_RPC:          args,
		SC_DERO_Deposit: 0, /// TODO add inital deposit
		Ringsize:        2,
	}
	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	log.Println("Set Prediction TX:", txid)

	return err
}

func postPrediction(price int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Post"}
	arg2 := rpc.Argument{Name: "price", DataType: "U", Value: price}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}
	p := &rpc.SC_Invoke_Params{
		SC_ID:           P_SC_ID,
		SC_RPC:          args,
		SC_DERO_Deposit: 0,
		Ringsize:        2,
	}
	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	log.Println("Post TX:", txid)

	return err
}

func endSports(num, team string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "S_end"}
	arg2 := rpc.Argument{Name: "n", DataType: "S", Value: num}
	arg3 := rpc.Argument{Name: "team", DataType: "S", Value: team}
	args := rpc.Arguments{arg1, arg2, arg3}
	txid := rpc.Transfer_Result{}
	params := &rpc.Transfer_Params{
		Transfers: []rpc.Transfer{},
		SC_Value:  0,
		SC_ID:     S_SC_ID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      400,
		Signer:    "",
	}
	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	log.Println("End Sports TX:", txid)

	return err
}

func endPredition(price int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	rpcClientW := jsonrpc.NewClientWithOpts(walletAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcLoginInput.Text)),
		},
	})

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "P_end"}
	arg2 := rpc.Argument{Name: "price", DataType: "U", Value: price}
	args := rpc.Arguments{arg1, arg2, arg2}
	txid := rpc.Transfer_Result{}

	params := &rpc.Transfer_Params{
		Transfers: []rpc.Transfer{},
		//SC_Code:   "",
		SC_Value: 0,
		SC_ID:    P_SC_ID,
		SC_RPC:   args,
		Ringsize: 2,
		Fees:     400,
		//Signer:   "",
	}
	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	log.Println("End Predition TX:", txid)

	return err
}
