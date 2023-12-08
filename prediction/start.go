package prediction

import (
	"fmt"
	"image/color"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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

const app_tag = "dPrediction"

var version = semver.MustParse("0.3.0-dev")
var gnomon = gnomes.NewGnomes()

// Check prediction package version
func Version() semver.Version {
	return version
}

// Start dPrediction dApp
func StartApp() {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
	gnomes.InitLogrusLog(logrus.InfoLevel)
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
	menu.Theme.Img = *canvas.NewImageFromResource(nil)
	d := dreams.AppObject{
		App:        a,
		Window:     w,
		Background: container.NewStack(&menu.Theme.Img),
	}
	d.SetChannels(1)

	closeFunc := func() {
		save := dreams.SaveData{
			Skin:   config.Skin,
			DBtype: gnomon.DBStorageType(),
		}

		if rpc.Daemon.Rpc == "" {
			save.Daemon = config.Daemon
		} else {
			save.Daemon = []string{rpc.Daemon.Rpc}
		}

		menu.WriteDreamsConfig(save)
		menu.SetClose(true)
		gnomon.Stop(app_tag)
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
	menu.Control.Ratings = make(map[string]uint64)
	gnomon.SetDBStorageType("boltdb")
	gnomon.SetFastsync(true)

	// Initialize asset widgets
	//asset_selects := []fyne.Widget{}

	// Create dwidget connection box with controls
	connect_box := dwidget.NewHorizontalEntries(app_tag, 1)
	connect_box.Button.OnTapped = func() {
		rpc.GetAddress(app_tag)
		rpc.Ping()
		if rpc.Daemon.IsConnected() && !gnomon.IsInitialized() && !gnomon.IsStarting() {
			filter := []string{
				GetPredictCode(0),
				GetSportsCode(0),
				rpc.GetSCCode(rpc.GnomonSCID),
				rpc.GetSCCode(rpc.RatingSCID)}

			go gnomes.StartGnomon(app_tag, gnomon.DBStorageType(), filter, 0, 0, nil)
		}
	}

	connect_box.Disconnect.OnChanged = func(b bool) {
		if !b {
			gnomon.Stop(app_tag)
		}
	}

	connect_box.AddDaemonOptions(config.Daemon)

	connect_box.Container.Objects[0].(*fyne.Container).Add(menu.StartIndicators())

	// Initialize profile widgets
	line := canvas.NewLine(bundle.TextColor)
	form := []*widget.FormItem{}
	form = append(form, widget.NewFormItem("Name", menu.NameEntry()))
	form = append(form, widget.NewFormItem("", layout.NewSpacer()))
	form = append(form, widget.NewFormItem("", container.NewVBox(line)))
	form = append(form, widget.NewFormItem("Avatar", holdero.AvatarSelect(menu.Assets.SCIDs)))
	form = append(form, widget.NewFormItem("Theme", menu.ThemeSelect()))
	form = append(form, widget.NewFormItem("", layout.NewSpacer()))
	form = append(form, widget.NewFormItem("", container.NewVBox(line)))

	profile_spacer := canvas.NewRectangle(color.Transparent)
	profile_spacer.SetMinSize(fyne.NewSize(450, 0))

	profile := container.NewCenter(container.NewBorder(profile_spacer, nil, nil, nil, widget.NewForm(form...)))

	// Layout tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Predict", LayoutPredictItems(&d)),
		container.NewTabItem("Sports", LayoutSportsItems(&d)),
		// TODO profile
		container.NewTabItem("Assets", menu.PlaceAssets(app_tag, profile, nil, resourceDServiceCirclePng, &d)),
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
					gnomes.GnomonEndPoint()
				}

				if rpc.Daemon.IsConnected() && gnomon.IsInitialized() {
					connect_box.Disconnect.SetChecked(true)
					if gnomon.IsRunning() {
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
					connect_box.Disconnect.SetChecked(false)
					Disconnected()
				}

				if rpc.Daemon.IsConnected() {
					rpc.Startup = false
				}

				d.SignalChannel()

			case <-d.Closing(): // exit
				logger.Printf("[%s] Closing...", app_tag)
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

	go func() {
		time.Sleep(450 * time.Millisecond)
		w.SetContent(container.NewStack(d.Background, container.NewStack(bundle.NewAlpha180(), tabs), container.NewVBox(layout.NewSpacer(), connect_box.Container)))
	}()
	w.ShowAndRun()
	<-done
	logger.Printf("[%s] Closed", app_tag)
}
