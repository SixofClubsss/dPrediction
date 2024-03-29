package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/SixofClubsss/dPrediction/prediction"
	"github.com/civilware/Gnomon/indexer"
	"github.com/civilware/Gnomon/structures"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"
)

// Run dReamsService process from dReams prediction package

var gnomon = gnomes.NewGnomes()
var enable_transfers bool
var logger = structures.Logger.WithFields(logrus.Fields{})
var command_line string = `dService
App to run dService as a single process, powered by Gnomon and dReams.

Usage:
  dService [options]
  dService -h | --help

Options:
  -h --help                      Show this screen.
  --daemon=<127.0.0.1:10102>     Set daemon rpc address to connect.
  --wallet=<127.0.0.1:10103>     Set wallet rpc address to connect.
  --login=<user:pass>     	 Wallet rpc user:pass for auth.
  --transfers=<false>            True/false value for enabling processing transfers to integrated address.
  --debug=<true>     		 True/false value for enabling terminal debug.
  --fastsync=<true>	         Gnomon option,  true/false value to define loading at chain height on start up.
  --num-parallel-blocks=<5>      Gnomon option,  defines the number of parallel blocks to index.`

func main() {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	v := prediction.Version().String()

	// Flags when starting dService
	arguments, err := docopt.ParseArgs(command_line, nil, v)
	if err != nil {
		logger.Fatalf("Error while parsing arguments: %s\n", err)
	}

	fastsync := true
	if arguments["--fastsync"] != nil {
		if arguments["--fastsync"].(string) == "false" {
			fastsync = false
		}
	}

	parallel := 1
	if arguments["--num-parallel-blocks"] != nil {
		s := arguments["--num-parallel-blocks"].(string)
		switch s {
		case "2":
			parallel = 2
		case "3":
			parallel = 3
		case "4":
			parallel = 4
		case "5":
			parallel = 5
		default:
			parallel = 1
		}
	}

	// Set default rpc params
	rpc.Daemon.Rpc = "127.0.0.1:10102"
	rpc.Wallet.Rpc = "127.0.0.1:10103"

	if arguments["--daemon"] != nil {
		if arguments["--daemon"].(string) != "" {
			rpc.Daemon.Rpc = arguments["--daemon"].(string)
		}
	}

	if arguments["--wallet"] != nil {
		if arguments["--wallet"].(string) != "" {
			rpc.Wallet.Rpc = arguments["--wallet"].(string)
		}
	}

	if arguments["--login"] != nil {
		if arguments["--login"].(string) != "" {
			rpc.Wallet.UserPass = arguments["--login"].(string)
		}
	}

	// Default false, integrated addresses generated through dReams
	transfers := false
	if arguments["--transfers"] != nil {
		if arguments["--transfers"].(string) == "true" {
			transfers = true
		}
	}

	debug := true
	if arguments["--debug"] != nil {
		if arguments["--debug"].(string) == "false" {
			debug = false
		}
	}

	arguments["--debug"] = false
	indexer.InitLog(arguments, os.Stderr)

	logger.Printf("[dService] %s  OS: %s  ARCH: %s  DREAMS: %s  GNOMON: %s\n", v, runtime.GOOS, runtime.GOARCH, rpc.Version(), structures.Version.String())

	// Check for daemon connection
	rpc.Ping()
	if !rpc.Daemon.IsConnected() {
		logger.Fatalf("[dService] Daemon %s not connected\n", rpc.Daemon.Rpc)
	}

	// Check for wallet connection
	rpc.GetAddress("dService")
	if !rpc.Wallet.IsConnected() {
		logger.Fatalf("[dService] Wallet %s not connected\n", rpc.Wallet.Rpc)
	}

	prediction.Service.Start()
	enable_transfers = transfers
	prediction.Service.Debug = debug
	gnomon.SetFastsync(fastsync, true, 10000)
	gnomon.SetParallel(parallel)
	prediction.Imported = true

	// Start dService from last payload format height at minimum
	height := prediction.PAYLOAD_FORMAT

	// Set up Gnomon search filters
	filter := []string{}
	predict := prediction.GetPredictCode(0)
	if predict != "" {
		filter = append(filter, predict)
	}

	sports := prediction.GetSportsCode(0)
	if sports != "" {
		filter = append(filter, sports)
	}

	// Handle ctrl+c close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		gnomon.Stop("dService")
		rpc.Wallet.Connected(false)
		prediction.Service.Stop()
		menu.SetClose(true)
		for prediction.Service.IsProcessing() {
			logger.Println("[dService] Waiting for service to close")
			time.Sleep(3 * time.Second)
		}
		logger.Println("[dService] Closing")
		os.Exit(0)
	}()

	// Start Gnomon with search filters
	go gnomes.StartGnomon("dService", "boltdb", filter, 0, 0, nil)

	// Routine for checking daemon, wallet connection and Gnomon sync
	go func() {
		for !menu.IsClosing() && !gnomon.IsInitialized() {
			time.Sleep(time.Second)
		}

		logger.Println("[dService] Starting when Gnomon is synced")
		height = uint64(gnomon.GetChainHeight())
		for !menu.IsClosing() && gnomon.IsRunning() && rpc.IsReady() {
			rpc.Ping()
			rpc.EchoWallet("dService")
			gnomon.IndexContains()
			if gnomon.GetLastHeight() >= gnomon.GetChainHeight()-3 && gnomon.HasIndex(9) {
				gnomon.Synced(true)
			} else {
				gnomon.Synced(false)
				if !gnomon.IsStarting() && gnomon.IsInitialized() {
					diff := gnomon.GetChainHeight() - gnomon.GetLastHeight()
					if diff > 3 && prediction.Service.Debug {
						logger.Printf("[dService] Gnomon has %d blocks to go\n", diff)
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()

	// Wait for Gnomon to sync
	for !menu.IsClosing() && (!gnomon.IsSynced() || gnomon.IsStatus("fastsyncing")) {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)

	// Populate SCID of connected wallet
	prediction.PopulatePredictions(nil)
	prediction.PopulateSports(nil)

	// Set added print text and make integrated address for transfers
	add := ""
	if enable_transfers {
		add = "and transactions"
		prediction.MakeIntegratedAddr(prediction.Service.Debug)
	}

	// Start dService
	logger.Printf("[dService] Processing payouts %s\n", add)
	prediction.RunService(height, true, enable_transfers)
}
