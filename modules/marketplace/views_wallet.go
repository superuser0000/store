package marketplace

import (
	"net/http"

	"github.com/gocraft/web"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ViewWalletRecoverGET(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "wallet/recover", c)
}

func (c *Context) ViewWalletRecoverPOST(w web.ResponseWriter, r *web.Request) {

	linkBitcoin := func() string {
		btcws := FindBitcoinWalletsForUser(r.FormValue("address"))
		if len(btcws) == 0 {
			return ""
		}

		user, _ := FindUserByUuid(btcws[0].UserUuid)
		if user != nil {
			return ""
		}

		btcws[0].UserUuid = c.ViewUser.Uuid
		btcws[0].Save()
		return "/wallet/bitcoin/receive"
	}

	linkEthereum := func() string {
		ethws := FindEthereumWalletsForUser(r.FormValue("address"))
		if len(ethws) == 0 {
			return ""
		}

		user, _ := FindUserByUuid(ethws[0].UserUuid)
		if user != nil {
			return ""
		}

		ethws[0].UserUuid = c.ViewUser.Uuid
		ethws[0].Save()
		return "/wallet/bitcoin/receive"
	}

	btcResult := linkBitcoin()
	ethResult := linkEthereum()

	if btcResult != "" {
		http.Redirect(w, r.Request, btcResult, 302)
		return
	}

	if ethResult != "" {
		http.Redirect(w, r.Request, ethResult, 302)
		return
	}

	c.Error = "No such wallet"
	c.ViewWalletRecoverGET(w, r)

}
