package marketplace

import (
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

type Context struct {
	*util.Context
	// Localization
	Localization Localization `json:"-"`

	// General
	CaptchaId string `json:"captcha_id,omitempty"`
	Error     string `json:"error,omitempty"`

	// User Permissions
	CanEdit            bool `json:"can_edit,omitempty"`
	EnableStoreStaffUI bool `json:"-,omitempty"`
	CanPostWarnings    bool `json:"-"`

	// Misc
	Pgp                 string                `json:"pgp,omitempty"`
	UserSettingsHistory []UserSettingsHistory `json:"user_settings_history,omitempty"`
	Language            string                `json:"language,omitempty"`

	// Paging & sorting
	SelectedPage  int    `json:"selected_page,omitempty"`
	Pages         []int  `json:"-,omitempty"`
	Page          int    `json:"page,omitempty"`
	NumberOfPages int    `json:"number_of_pages,omitempty"`
	Query         string `json:"query,omitempty"`
	SortBy        string `json:"sort_by,omitempty"`
	Domestic      bool   `json:"domestic,omitempty"`

	// Static Pages
	StaticPage  StaticPage   `json:"-,omitempty"`
	StaticPages []StaticPage `json:"-,omitempty"`

	// Menu
	Categories             []Category         `json:"-,omitempty"`
	Cities                 map[string]int     `json:"-,omitempty"`
	City                   string             `json:"city,omitempty"`
	GeoCities              []City             `json:"geo_cities,omitempty"`
	CityMetroStations      []CityMetroStation `json:"metro_stations,omitempty"`
	CityID                 int                `json:"city_id,omitempty"`
	Countries              []Country          `json:"countries,omitempty"`
	Quantity               int                `json:"quantity,omitempty"`
	SelectedPackageType    string             `json:"selected_package_type,omitempty"`
	SelectedSection        string             `json:"-,omitempty"`
	SelectedSectionID      int                `json:"-,omitempty"`
	SelectedStatus         string             `json:"selected_status,omitempty"`
	SelectedShippingStatus string             `json:"selected_shipping_status,omitempty"`
	ShippingFrom           string             `json:"shipping_from,omitempty"`
	ShippingFromList       []string           `json:"shipping_from_list,omitempty"`
	ShippingTo             string             `json:"shipping_to,omitempty"`
	ShippingToList         []string           `json:"shipping_to_list,omitempty"`
	Account                string             `json:"account,omitempty"`
	Currency               string             `json:"currency,omitempty"`

	// Categories
	Category       string `json:"category,omitempty"`
	SubCategory    string `json:"sub_category,omitempty"`
	SubSubCategory string `json:"sub_sub_category,omitempty"`
	CategoryID     int    `json:"category_id,omitempty"`

	// Items page
	GroupPackages                        []GroupPackage                       `json:"-,omitempty"`
	GroupPackagesByTypeOriginDestination map[GroupedPackageKey][]GroupPackage `json:"-,omitempty"`
	GroupAvailability                    *GroupPackage                        `json:"group_package,omitempty"`
	NumberOfItems                        int                                  `json:"number_of_items,omitempty"`
	PackageCurrency                      string                               `json:"package_currency,omitempty"`
	PackagePrice                         string                               `json:"package_price,omitempty"`

	// Models
	ExtendedUsers []ExtendedUser `json:"-,omitempty"`
	Thread        Thread         `json:"-,omitempty"`
	Threads       []Thread       `json:"-,omitempty"`
	RatingReview  *RatingReview  `json:"rating-review,omitempty"`

	// View Models
	ViewCurrentTransactionStatuses []ViewCurrentTransactionStatus `json:"transaction_statuses,omitempty"`
	ViewExtendedUsers              []ViewExtendedUser             `json:"users,omitempty"`
	ViewItem                       *ViewItem                      `json:"item,omitempty,omitempty"`
	ViewItemCategories             []ViewItemCategory             `json:"item_categories,omitempty"`
	ViewItemCategory               *ViewItemCategory              `json:"item_category,omitempty"`
	ViewItems                      []ViewItem                     `json:"items,omitempty"`
	ViewMarketplaceUser            *ViewUser                      `json:"marketplace_user,omitempty"`
	ViewMessage                    ViewMessage                    `json:"-"`
	ViewMessages                   []ViewMessage                  `json:"-"`
	ViewPackage                    ViewPackage                    `json:"-"`
	ViewPackages                   []ViewPackage                  `json:"packages,omitempty"`
	ViewPrivateThreads             []ViewPrivateThread            `json:"threads,omitempty"`
	ViewStore                      *ViewStore                     `json:"store,omitempty"`
	ViewStores                     []ViewStore                    `json:"stores,omitempty"`
	ViewStoreWarnings              []ViewStoreWarning             `json:"warnings"`
	ViewThread                     *ViewThread                    `json:"thread,omitempty"`
	ViewThreads                    []ViewThread                   `json:"-"`
	ViewTransaction                *ViewTransaction               `json:"transaction,omitempty"`
	ViewTransactions               []ViewTransaction              `json:"transactions,omitempty"`
	ViewUser                       *ViewUser                      `json:"user,omitempty"`
	ViewUsers                      []ViewUser                     `json:"-"`
	ViewUserStore                  *ViewStore                     `json:"user_store,omitempty"`

	// Stats
	NumberOfDailyTransactions     int `json:"-"`
	NumberOfMonthlyTransactions   int `json:"-"`
	NumberOfPrivateMessages       int `json:"-"`
	NumberOfSupportMessages       int `json:"-"`
	NumberOfTransactions          int `json:"-"`
	NumberOfUnreadPrivateMessages int `json:"-"`
	NumberOfUnreadSupportMessages int `json:"-"`
	NumberOfWeeklyTransactions    int `json:"-"`
	NumberOfDisputes              int `json:"-"`

	// Admin Stats
	NumberOfUsers int `json:"-"`

	// --- Vendor Statistics ---
	NumberOfVendors       int `json:"-"`
	NumberOfFreeVendors   int `json:"-"`
	NumberOfGoldVendors   int `json:"-"`
	NumberOfSilverVendors int `json:"-"`
	NumberOfBronzeVendors int `json:"-"`

	// --- User Statistics ---
	NumberOfNewUsers           int         `json:"-"`
	NumberOfActiveUsers        int         `json:"-"`
	NumberOfWeeklyActiveUsers  int         `json:"-"`
	NumberOfOnlineUsers        int         `json:"-"`
	NumberOfMonthlyActiveUsers int         `json:"-"`
	NumberOfInvitedUsers       int         `json:"-"`
	StatsItems                 []StatsItem `json:"-"`

	// --- Staff Statistics ---
	StaffStats StaffStats `json:"-"`

	// Auth
	SecretText string `json:"secret_text,omitempty"`
	InviteCode string `json:"invite_code,omitempty"`

	// Bitcoin Wallets
	UserBitcoinBalance       *apis.BTCWalletBalance    `json:"btc_balance,omitempty"`
	UserBitcoinWallets       UserBitcoinWallets        `json:"-"`
	UserBitcoinWallet        *UserBitcoinWallet        `json:"btc_wallet,omitempty"`
	UserBitcoinWalletActions []UserBitcoinWalletAction `json:"btc_wallet_actions,omitempty"`

	// Ethereum Wallets
	UserEthereumBalance       *apis.ETHWalletBalance     `json:"eth_balance,omitempty"`
	UserEthereumWallets       UserEthereumWallets        `json:"-"`
	UserEthereumWallet        *UserEthereumWallet        `json:"eth_wallet,omitempty"`
	UserEthereumWalletActions []UserEthereumWalletAction `json:"eth_wallet_actions,omitempty"`

	// Referrals
	ReferralPayments []ReferralPayment `json:"-"`

	//Dispute
	Dispute      *Dispute      `json:"dispute,omitempty"`
	Disputes     []Dispute     `json:"disputes,omitempty"`
	DisputeClaim *DisputeClaim `json:"dispute_claim,omitempty"`

	// Deposit
	Deposit         *Deposit        `json:"-"`
	Deposits        Deposits        `json:"-"`
	DepositsSummary DepositsSummary `json:"-"`

	// Support
	ViewSupportTicket  *ViewSupportTicket  `json:"ticket,omitempty"`
	ViewSupportTickets []ViewSupportTicket `json:"tickets,omitempty"`

	// New Items List page
	ViewSerpItems  []ViewSerpItem  `json:"serp_items,omitempty"`
	ViewSerpStores []ViewSerpStore `json:"serp_stores,omitempty"`

	// Front Page Items Lists
	CountVendors int `json:"-,omitempty"`
	CountItems   int `json:"-,omitempty"`

	// Currency Rates
	CurrencyRates map[string]map[string]float64 `json:"currency_rates"`
	USDBTCRate    float64                       `json:"-"`

	// Wallet page
	BTCFee           float64                `json:"btc_fee,omitempty"`
	Amount           float64                `json:"amount,omitempty"`
	Address          string                 `json:"address,omitempty"`
	Description      string                 `json:"description,omitempty"`
	BTCPaymentResult *apis.BTCPaymentResult `json:"btc_payment_result,omitempty"`
	ETHPaymentResult *apis.ETHPaymentResult `json:"eth_payment_results,omitempty""`

	// Membership Plans
	PriceBTC float64 `json:"-"`
	PriceETH float64 `json:"-"`
	PriceUSD float64 `json:"-"`

	// Advertising
	Advertisings     []Advertising `json:"-"`
	AdvertisingCost  float64       `json:"-"`
	HideAdvertisings bool          `json:"-"`

	// ApiSession
	APISession   *APISession `json:"api_session,omitempty"`
	IsAPIRequest bool        `json:"-"`

	// CSRF Token
	CSRFToken string `json:"-"`
	SiteName  string `json:"-"`
	SiteURL   string `json:"-"`

	// Page Statistics
	PageRenderStart time.Time `json:"-"`
	NumberOfQueries int       `json:"-"`
	PageRenderTime  uint64    `json:"-"`

	// Notifications
	Notification Notification `json:"notifications"`
}
