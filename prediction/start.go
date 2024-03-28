package prediction

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/SixofClubsss/Holdero/holdero"
	"github.com/blang/semver/v4"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/sirupsen/logrus"
)

const (
	appName = "dPrediction"
	appID   = "dreamdapps.io.dprediction"
)

var version = semver.MustParse("0.3.1-dev.3")
var gnomon = gnomes.NewGnomes()

// Check prediction package version
func Version() semver.Version {
	return version
}

// Start dPrediction dApp
func StartApp() {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	// Initialize logrus logger to stdout
	gnomes.InitLogrusLog(logrus.InfoLevel)

	// Read config.json file
	config := menu.GetSettings(appName)

	// Initialize Fyne app and window as dreams.AppObject
	d := dreams.NewFyneApp(
		appID,
		appName,
		"P2P Predictions and Sports Bet",
		bundle.DeroTheme(config.Skin),
		resourceDServiceIconPng,
		menu.DefaultBackgroundResource(),
		true)

	// Set one channel for dPrediction routine
	d.SetChannels(1)

	// Initialize closing channels and func
	done := make(chan struct{})

	closeFunc := func() {
		save := dreams.SaveData{
			Skin:   config.Skin,
			DBtype: gnomon.DBStorageType(),
			Theme:  dreams.Theme.Name,
		}

		if rpc.Daemon.Rpc == "" {
			save.Daemon = config.Daemon
		} else {
			save.Daemon = []string{rpc.Daemon.Rpc}
		}

		menu.StoreSettings(save)
		menu.SetClose(true)
		gnomon.Stop(appName)
		d.StopProcess()
		d.Window.Close()
	}

	d.Window.SetCloseIntercept(closeFunc)

	// Handle ctrl-c close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		closeFunc()
	}()

	// Initialize Gnomon vars
	gnomon.SetDBStorageType("boltdb")
	gnomon.SetFastsync(true, true, 10000)

	// Create dwidget connection box, using default OnTapped for connections
	connection := dwidget.NewHorizontalEntries(appName, 1, &d)

	connection.Connected.OnChanged = func(b bool) {
		if b {
			if rpc.Daemon.IsConnected() && !gnomon.IsInitialized() && !gnomon.IsStarting() {
				filter := []string{
					GetPredictCode(0),
					GetSportsCode(0),
					rpc.GetSCCode(rpc.GnomonSCID),
					rpc.GetSCCode(rpc.RatingSCID)}

				go gnomes.StartGnomon(appName, gnomon.DBStorageType(), filter, 0, 0, nil)
			}
		} else {
			gnomon.Stop(appName)
		}
	}

	// Set any saved daemon configs
	connection.AddDaemonOptions(config.Daemon)

	// Adding dReams indicator panel for wallet, daemon and Gnomon
	connection.AddIndicator(menu.StartIndicators(nil))

	// Initialize profile widgets
	line := canvas.NewLine(bundle.TextColor)
	form := []*widget.FormItem{}
	form = append(form, widget.NewFormItem("Name", menu.NameEntry()))
	form = append(form, widget.NewFormItem("", layout.NewSpacer()))
	form = append(form, widget.NewFormItem("", container.NewVBox(line)))
	form = append(form, widget.NewFormItem("Avatar", holdero.AvatarSelect(menu.Assets.SCIDs)))
	form = append(form, widget.NewFormItem("Theme", menu.ThemeSelect(&d)))
	form = append(form, widget.NewFormItem("", layout.NewSpacer()))
	form = append(form, widget.NewFormItem("", container.NewVBox(line)))

	profile := container.NewCenter(container.NewBorder(dwidget.NewSpacer(450, 0), nil, nil, nil, widget.NewForm(form...)))

	// Layout tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Predict", LayoutPredictions(&d)),
		container.NewTabItem("Sports", LayoutSports(&d)),
		container.NewTabItem("Assets", menu.PlaceAssets(appName, profile, nil, resourceDServiceCirclePng, &d)),
		container.NewTabItem("Log", rpc.SessionLog(appName, version)))

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
				rpc.Wallet.Sync()

				connection.RefreshBalance()

				if rpc.Daemon.IsConnected() {
					connection.Connected.SetChecked(true)
					if gnomon.IsRunning() {
						gnomes.EndPoint()
						menu.DisableIndexControls(false)
						gnomon.IndexContains()
						menu.Info.RefreshIndexed()
						if gnomon.HasIndex(2) {
							gnomon.Checked(true)
						}
					}

					if gnomon.GetLastHeight() >= gnomon.GetChainHeight()-3 {
						Predict.Public.List.Refresh()
						Sports.Public.List.Refresh()
						gnomon.Synced(true)
					} else {
						gnomon.Synced(false)
						gnomon.Checked(false)
					}
				} else {
					menu.DisableIndexControls(true)
					connection.Connected.SetChecked(false)
					Disconnected()
				}

				d.SignalChannel()

			case <-d.Closing(): // exit
				logger.Printf("[%s] Closing...", appName)
				if gnomes.Indicator.Icon != nil {
					gnomes.Indicator.Icon.Stop()
				}
				ticker.Stop()
				d.CloseAllDapps()
				time.Sleep(time.Second)
				done <- struct{}{}
				return
			}
		}
	}()

	// Start app and place layout
	go func() {
		time.Sleep(450 * time.Millisecond)
		d.Window.SetContent(container.NewStack(d.Background, container.NewStack(bundle.NewAlpha180(), tabs), container.NewVBox(layout.NewSpacer(), connection.Container)))
	}()
	d.Window.ShowAndRun()
	<-done
	logger.Printf("[%s] Closed", appName)
}
