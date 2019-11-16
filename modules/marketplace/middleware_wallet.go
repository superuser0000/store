package marketplace

import (
	"github.com/gocraft/web"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
)

func (c *Context) BitcoinWalletMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.UserBitcoinWallets = c.ViewUser.FindUserBitcoinWallets()
	for _, w := range c.UserBitcoinWallets {
		w.UpdateBalance(false)
	}
	if len(c.UserBitcoinWallets) > 0 {
		c.UserBitcoinWallet = &c.UserBitcoinWallets[0]
	}
	balance := c.UserBitcoinWallets.Balance()
	c.UserBitcoinBalance = &balance
	next(w, r)
}

func (c *Context) EthereumWalletMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.UserEthereumWallets = c.ViewUser.FindUserEthereumWallets()
	for _, w := range c.UserEthereumWallets {
		w.UpdateBalance(false)
	}
	if len(c.UserEthereumWallets) > 0 {
		c.UserEthereumWallet = &c.UserEthereumWallets[0]
	}
	c.UserEthereumBalance = &apis.ETHWalletBalance{
		Balance: c.UserEthereumWallets.Balance().Balance,
	}
	next(w, r)
}

func (c *Context) WalletsMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	if c.ViewUser != nil {
		c.UserBitcoinWallets = c.ViewUser.FindUserBitcoinWallets()
		c.UserEthereumWallets = c.ViewUser.FindUserEthereumWallets()
		btcBalance := c.UserBitcoinWallets.Balance()
		c.UserBitcoinBalance = &btcBalance
		c.UserEthereumBalance = &apis.ETHWalletBalance{
			Balance: c.UserEthereumWallets.Balance().Balance,
		}
	}

	next(w, r)
}
