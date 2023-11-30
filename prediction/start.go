package prediction

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/blang/semver/v4"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/sirupsen/logrus"
)

const app_tag = "dPrediction"

var version = semver.MustParse("0.3.0-dev")

// Check prediction package version
func Version() semver.Version {
	return version
}

// Start dPrediction dApp
func StartApp() {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
	menu.InitLogrusLog(logrus.InfoLevel)
	config := menu.ReadDreamsConfig(app_tag)

	// Initialize Fyne app and window
	a := app.NewWithID("dPrediction Desktop Client")
	a.Settings().SetTheme(bundle.DeroTheme(config.Skin))
	w := a.NewWindow(app_tag)
	w.SetIcon(resourceDServiceIconPng)
	w.Resize(fyne.NewSize(1400, 800))
	w.SetMaster()
	done := make(chan struct{})

	// Initialize dReams AppObject and close func
	dreams.Theme.Img = *canvas.NewImageFromResource(nil)
	d := dreams.AppObject{
		App:        a,
		Window:     w,
		Background: container.NewStack(&dreams.Theme.Img),
	}
	d.SetChannels(1)

	closeFunc := func() {
		save := dreams.SaveData{
			Skin:   config.Skin,
			DBtype: menu.Gnomes.DBType,
		}

		if rpc.Daemon.Rpc == "" {
			save.Daemon = config.Daemon
		} else {
			save.Daemon = []string{rpc.Daemon.Rpc}
		}

		menu.WriteDreamsConfig(save)
		menu.CloseAppSignal(true)
		menu.Gnomes.Stop(app_tag)
		d.StopProcess()
		w.Close()
	}

	w.SetCloseIntercept(closeFunc)

	// Handle ctrl-c close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		closeFunc()
	}()

	// Initialize vars
	rpc.InitBalances()
	menu.Control.Contract_rating = make(map[string]uint64)
	menu.Gnomes.DBType = "boltdb"
	menu.Gnomes.Fast = true

	// Initialize asset widgets
	asset_selects := []fyne.Widget{}

	// Create dwidget connection box with controls
	connect_box := dwidget.NewHorizontalEntries(app_tag, 1)
	connect_box.Button.OnTapped = func() {
		rpc.GetAddress(app_tag)
		rpc.Ping()
		if rpc.Daemon.IsConnected() && !menu.Gnomes.IsInitialized() && !menu.Gnomes.Start {
			filter := []string{
				GetPredictCode(0),
				GetSportsCode(0),
				rpc.GetSCCode(rpc.GnomonSCID),
				rpc.GetSCCode(rpc.RatingSCID)}

			go menu.StartGnomon(app_tag, menu.Gnomes.DBType, filter, 0, 0, nil)
		}
	}

	connect_box.Disconnect.OnChanged = func(b bool) {
		if !b {
			menu.Gnomes.Stop(app_tag)
		}
	}

	connect_box.AddDaemonOptions(config.Daemon)

	connect_box.Container.Objects[0].(*fyne.Container).Add(menu.StartIndicators())

	// Layout tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Predict", LayoutPredictItems(&d)),
		container.NewTabItem("Sports", LayoutSportsItems(&d)),
		container.NewTabItem("Assets", menu.PlaceAssets(app_tag, asset_selects, resourceDServiceCirclePng, d.Window)),
		container.NewTabItem("Log", rpc.SessionLog(app_tag, version)))

	tabs.SetTabLocation(container.TabLocationBottom)
	d.SetTab("Predict")
	tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Predict":
			d.SetTab("Predict")
		case "Sports":
			d.SetTab("Sports")
		default:
			d.SetTab("")
			// nothing
		}
	}

	// Stand alone process
	go func() {
		time.Sleep(3 * time.Second)
		ticker := time.NewTicker(3 * time.Second)

		for {
			select {
			case <-ticker.C: // do on interval
				rpc.Ping()
				rpc.EchoWallet(app_tag)
				rpc.GetDreamsBalances(rpc.SCIDs)
				rpc.GetWalletHeight(app_tag)

				connect_box.RefreshBalance()
				if !rpc.Startup {
					menu.GnomonEndPoint()
				}

				if rpc.Daemon.IsConnected() && menu.Gnomes.IsInitialized() {
					connect_box.Disconnect.SetChecked(true)
					if menu.Gnomes.IsRunning() {
						menu.DisableIndexControls(false)
						menu.Gnomes.IndexContains()
						scids := " Indexed SCIDs: " + strconv.Itoa(int(menu.Gnomes.SCIDS))
						menu.Assets.Gnomes_index.Text = scids
						menu.Assets.Gnomes_index.Refresh()
						if menu.Gnomes.HasIndex(2) {
							menu.Gnomes.Checked(true)
						}
					}

					if menu.Gnomes.Indexer.LastIndexedHeight >= menu.Gnomes.Indexer.ChainHeight-3 {
						Predict.Public.List.Refresh()
						Sports.Public.List.Refresh()
						menu.Gnomes.Synced(true)
					} else {
						menu.Gnomes.Synced(false)
						menu.Gnomes.Checked(false)
					}
				} else {
					menu.DisableIndexControls(true)
					connect_box.Disconnect.SetChecked(false)
				}

				if rpc.Daemon.IsConnected() {
					rpc.Startup = false
				}

				d.SignalChannel()

			case <-d.Closing(): // exit
				logger.Printf("[%s] Closing...", app_tag)
				if menu.Gnomes.Icon_ind != nil {
					menu.Gnomes.Icon_ind.Stop()
				}
				ticker.Stop()
				d.CloseAllDapps()
				time.Sleep(time.Second)
				done <- struct{}{}
				return
			}
		}
	}()

	go func() {
		time.Sleep(450 * time.Millisecond)
		w.SetContent(container.NewStack(d.Background, container.NewStack(bundle.NewAlpha180(), tabs), container.NewVBox(layout.NewSpacer(), connect_box.Container)))
	}()
	w.ShowAndRun()
	<-done
	logger.Printf("[%s] Closed", app_tag)
}
