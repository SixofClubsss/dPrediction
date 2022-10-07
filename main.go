package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var myApp = app.New()
var myWindow fyne.Window
var quit chan struct{}
var daemonAddress string
var walletAddress string

const (
	MIN_WIDTH  = 380
	MIN_HEIGHT = 800
)

func main() {
	myApp = app.New()
	myWindow = myApp.NewWindow("dPrediction")
	myWindow.SetMaster()
	myWindow.Resize(fyne.NewSize(MIN_WIDTH, MIN_HEIGHT))
	myWindow.SetFixedSize(true)
	myWindow.SetIcon(resourceDReamTablesIconPng)
	myWindow.SetCloseIntercept(func() { /// do when app close
		stopLoop()
		log.Println("Putting gnomes to sleep. This will take ~5sec.")
		defaultIndexer.Close()
		myWindow.Close()
	})

	leaders = make(map[string]uint64)
	/// organize content
	wallet_content_a := container.NewVBox(balanceDisplay(), rpcLoginEdit(), rpcWalletEdit(), rpcConnectButton())
	wallet_content_b := container.NewVBox(layout.NewSpacer(), daemonSelectOption(), daemonConnectBox(), walletConnectBox(), heightDisplay())
	wallet_content_c := container.NewVBox(ownerButton())
	wallet_content := container.NewVBox(wallet_content_c, layout.NewSpacer(), wallet_content_b, wallet_content_a)

	sports_content := container.NewVBox(s_scDisplayTop())
	sports_scroll := container.NewVScroll(sports_content)
	sports_scroll.SetMinSize(fyne.NewSize(160, 660))
	sports_buttons := container.NewVBox(gameOptions(), teamA(), teamB())
	s_sc := container.NewVBox(sports_scroll, layout.NewSpacer(), sports_buttons)

	predict_content_a := container.NewVBox(higher(), lower())
	predict_content_b := container.NewVBox(nameEdit(), change())
	predict_content_c := container.NewVBox(p_scDisplayTop(), p_scDisplayBottom())
	predict_scroll := container.NewScroll(predict_content_c)
	predict_scroll.SetMinSize(fyne.NewSize(180, 600))
	predict_content := container.NewVBox(predict_scroll, layout.NewSpacer(), predict_content_b, predict_content_a)

	leaders_scroll := container.NewScroll(leadersDisplay())
	leaders_scroll.SetMinSize(fyne.NewSize(180, 780))
	leaders_contnet := container.NewVBox(leaders_scroll, layout.NewSpacer(), remove())

	tabs := container.NewAppTabs(
		container.NewTabItem("Wallet", wallet_content),
		container.NewTabItem("Sports", s_sc),
		container.NewTabItem("Predictions", predict_content),
		container.NewTabItem("Leaderboard", leaders_contnet),
	)
	tabs.OnSelected = func(ti *container.TabItem) {
		getBook()
		makeLeaderBoard()
	}
	/// start process loop, set content and start main app
	fetchLoop()
	tabs.SetTabLocation(container.TabLocationTop)
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

func fetchLoop() { ///
	var ticker = time.NewTicker(6 * time.Second)
	quit = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C: /// do on interval
				ping()
				isDaemonConnected()
				isWalletConnected()
				getHeight()
				p_getSC()
			case <-quit: /// exit loop
				log.Println("Exiting...")
				ticker.Stop()
				return
			}
		}
	}()
}

func stopLoop() { /// exit loop on close
	quit <- struct{}{}
}
