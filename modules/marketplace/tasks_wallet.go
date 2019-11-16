package marketplace

import (
	"github.com/jasonlvhit/gocron"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func TaskUpdateRecentBitcoinWallets() {
	recentWallets := FindRecentBitcoinWallets()
	util.Log.Debug("[Task] [TaskUpdateRecentBitcoinWallets], number of wallets to update: %d", len(recentWallets))
	for _, w := range recentWallets {
		w.UpdateBalance(false)
	}
}

func TaskUpdateAllBitcoinWallets() {
	util.Log.Debug("[Task] [TaskUpdateAllBitcoinWallets]")
	wallets := FindAllBitcoinWallets()
	for _, w := range wallets {
		w.UpdateBalance(false)
	}
}

func TaskUpdateRecentEthereumWallets() {
	recentWallets := FindRecentEthereumWallets()
	util.Log.Debug("[Task] [TaskUpdateRecentEthereumWallets], number of wallets to update: %d", len(recentWallets))
	for _, w := range recentWallets {
		w.UpdateBalance(false)
	}

}

func TaskUpdateAllEthereumWallets() {
	util.Log.Debug("[Task] [TaskUpdateAllEthereumWallets]")
	wallets := FindAllEthereumWallets()
	for _, w := range wallets {
		w.UpdateBalance(false)
	}
}

func StartWalletsCron() {
	gocron.Every(24).Hours().Do(TaskUpdateAllBitcoinWallets)
	gocron.Every(24).Hours().Do(TaskUpdateAllEthereumWallets)

	gocron.Every(1).Hour().Do(TaskUpdateRecentBitcoinWallets)
	gocron.Every(1).Hour().Do(TaskUpdateRecentEthereumWallets)
}
