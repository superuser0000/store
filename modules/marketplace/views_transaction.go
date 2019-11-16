package marketplace

import (
	"math"
	"net/http"
	"strconv"

	btcqr "github.com/GeertJohan/go.btcqr"
	"github.com/dchest/captcha"
	"github.com/gocraft/web"
	"github.com/mojocn/base64Captcha"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ShowTransaction(w web.ResponseWriter, r *web.Request) {
	c.CaptchaId = captcha.New()
	c.RatingReview, _ = FindRatingReviewByTransactionUuid(c.ViewTransaction.Uuid)
	// this fixes template not rendering issue
	if c.RatingReview == nil {
		rr := RatingReview{}
		c.RatingReview = &rr
	}
	if len(r.URL.Query()["section"]) > 0 {
		section := r.URL.Query()["section"][0]
		c.SelectedSection = section
	} else {
		c.SelectedSection = "payment"
	}
	util.RenderTemplate(w, "transaction/show", c)
}

func (c *Context) ShowTransactionPOST(w web.ResponseWriter, r *web.Request) {
	isCaptchaValid := base64Captcha.VerifyCaptcha(r.FormValue("captcha_id"), r.FormValue("captcha")) || c.ViewUser.IsAdmin || c.ViewUser.IsStaff || c.IsAPIRequest
	if !isCaptchaValid {
		c.Error = "Invalid captcha"
		c.ShowTransaction(w, r)
		return
	}
	message, err := CreateMessage(r.FormValue("text"), c.Thread, *c.ViewUser.User)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ShowTransaction(w, r)
		return
	}
	err = message.AddImage(r)
	if err != nil {
		c.Error = err.Error()
		c.ViewMessage = message.ViewMessage(c.ViewUser.Language)
		c.ShowTransaction(w, r)
		return
	}

	c.TransactionMiddleware(w, r, c.ShowTransaction)
}

func (c *Context) UpdateTransaction(w web.ResponseWriter, r *web.Request) {
	transaction, _ := FindTransactionByUuid(r.PathParams["transaction"])
	if transaction == nil {
		http.NotFound(w, r.Request)
		return
	}
	transaction.UpdateTransactionStatus()
	viewTransaction := transaction.ViewTransaction()
	c.ViewTransaction = &viewTransaction
	c.ShowTransaction(w, r)
}

func (c *Context) ListCurrentTransactionStatuses(w web.ResponseWriter, r *web.Request) {
	pageSize := 20
	if len(r.URL.Query()["status"]) > 0 {
		c.SelectedStatus = r.URL.Query()["status"][0]
		if c.SelectedStatus == "DISPATCH PENDING" {
			c.SelectedShippingStatus = c.SelectedStatus
			c.SelectedStatus = ""
		}
	}
	c.Page = 1
	if len(r.URL.Query()["page"]) > 0 {
		page, err := strconv.ParseInt(r.URL.Query()["page"][0], 10, 32)
		if err != nil || page < 0 {
			http.NotFound(w, r.Request)
			return
		}
		c.Page = int(page)
	}

	c.NumberOfTransactions = CountCurrentTransactionStatusesForBuyer(c.ViewUser.Uuid, c.SelectedStatus, c.SelectedShippingStatus, false)
	c.ViewCurrentTransactionStatuses = FindCurrentTransactionStatusesForBuyer(
		c.ViewUser.Uuid, c.SelectedStatus, false, c.Page-1, pageSize).
		ViewCurrentTransactionStatuses(c.ViewUser.Language)

	if c.ViewUserStore != nil && c.ViewUserStore.Uuid != "" {
		c.NumberOfTransactions += CountCurrentTransactionStatusesForStore(c.ViewUserStore.Uuid, c.SelectedStatus, c.SelectedShippingStatus, false)
		transactions := FindCurrentTransactionStatusesForStore(
			c.ViewUserStore.Uuid, c.SelectedStatus, c.SelectedShippingStatus, false, c.Page-1, pageSize).
			ViewCurrentTransactionStatuses(c.ViewUser.Language)
		c.ViewCurrentTransactionStatuses = append(c.ViewCurrentTransactionStatuses, transactions...)
	}

	c.NumberOfPages = int(math.Ceil(float64(c.NumberOfTransactions) / float64(pageSize)))
	for i := 0; i < c.NumberOfPages; i++ { // paging
		c.Pages = append(c.Pages, i+1)
	}
	util.RenderTemplate(w, "transaction/list", c)
}

func (c *Context) ViewTransactionRecoverGET(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "transaction/recover", c)
}

func (c *Context) ViewTransactionRecoverPOST(w web.ResponseWriter, r *web.Request) {

	if c.ViewUser.IsStore {
		c.Error = "You are store"
		c.ViewTransactionRecoverGET(w, r)
		return
	}

	tx, _ := FindTransactionByUuidRaw(r.FormValue("address"))
	if tx == nil {
		c.Error = "No such transaction"
		c.ViewTransactionRecoverGET(w, r)
		return
	}

	user, _ := FindUserByUuid(tx.BuyerUuid)
	if user != nil {
		c.Error = "Transaction already linked"
		c.ViewTransactionRecoverGET(w, r)
		return
	}

	tx.BuyerUuid = c.ViewUser.Uuid
	http.Redirect(w, r.Request, "/payments/"+tx.Uuid, 302)
}

func (c *Context) TransactionImage(w web.ResponseWriter, r *web.Request) {
	btcqr.DefaultConfig.Scheme = "bitcoin"

	var amount float64
	if c.ViewTransaction.Type == "bitcoin" {
		amount = c.ViewTransaction.BitcoinTransaction.Amount
	}

	req := &btcqr.Request{
		Address: c.ViewTransaction.Uuid,
		Amount:  amount,
		Label:   c.ViewTransaction.Description,
		Message: "",
	}
	code, err := req.GenerateQR()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	png := code.PNG()
	w.Header().Set("Content-type", "image/png")
	w.Write(png)
}

func (c *Context) CompleteTransactionPOST(w web.ResponseWriter, r *web.Request) {

	// Review
	review, _ := FindRatingReviewByTransactionUuid(c.ViewTransaction.Uuid)
	if review == nil {
		review = &RatingReview{
			Uuid: util.GenerateUuid(),
		}
	}

	// quality
	itemQuality, err := strconv.Atoi(r.FormValue("item_quality"))
	if err != nil || itemQuality < 0 || itemQuality > 5 {
		c.Error = "Wrong input for item quality"
		http.NotFound(w, r.Request)
		return
	}
	itemReview := r.FormValue("item_review")
	if len(itemReview) > 255 {
		itemReview = itemReview[0:255]
	}
	// package
	marketplaceQuality, err := strconv.Atoi(r.FormValue("marketplace_quality"))
	if err != nil || marketplaceQuality < 0 || marketplaceQuality > 5 {
		c.Error = "Wrong input for marketplace quality"
		http.NotFound(w, r.Request)
		return
	}
	marketplaceReview := r.FormValue("marketplace_review")
	if len(marketplaceReview) > 255 {
		marketplaceReview = marketplaceReview[0:255]
	}
	// seller
	sellerQuality, err := strconv.Atoi(r.FormValue("seller_quality"))
	if err != nil || sellerQuality < 0 || sellerQuality > 5 {
		c.Error = "Wrong input for seller quality"
		http.NotFound(w, r.Request)
		return
	}
	sellerReview := r.FormValue("seller_review")
	if len(sellerReview) > 255 {
		sellerReview = sellerReview[0:255]
	}

	pkg, _ := FindPackageByUuid(c.ViewTransaction.PackageUuid)
	if pkg != nil {
		review.ItemUuid = pkg.ItemUuid
	}

	review.ItemReview = itemReview
	review.ItemScore = itemQuality
	review.MarketplaceReview = marketplaceReview
	review.MarketplaceScore = marketplaceQuality
	review.SellerReview = sellerReview
	review.SellerScore = sellerQuality
	review.TransactionUuid = c.ViewTransaction.Uuid
	review.StoreUuid = c.ViewTransaction.StoreUuid
	review.UserUuid = c.ViewTransaction.BuyerUuid

	review.Save()

	redirectUrl := "/payments/" + c.ViewTransaction.Uuid + "?section=review"
	util.RedirectOrAPIResponse(w, r, redirectUrl, c, c.IsAPIRequest)
}

func (c *Context) SetTransactionShippingStatus(w web.ResponseWriter, r *web.Request) {
	status := r.FormValue("shipping_status")
	if !(status == "DISPATCHED" || status == "SHIPPED") {
		http.NotFound(w, r.Request)
		return
	}
	c.ViewTransaction.SetShippingStatus(status, "Shipping status changed to "+status, c.ViewUser.Uuid)

	redirectUrl := "/payments/" + c.ViewTransaction.Uuid
	util.RedirectOrAPIResponse(w, r, redirectUrl, c, c.IsAPIRequest)
}

func (c *Context) ReleaseTransaction(w web.ResponseWriter, r *web.Request) {
	if !c.ViewTransaction.Store.UserIsAdministration(c.ViewUser.Uuid) {
		model := c.ViewTransaction.DBModel()

		err := model.Release("User released transaction", c.ViewUser.Uuid)
		if err != nil {
			model.SetTransactionStatus(
				model.CurrentPaymentStatus(),
				model.CurrentAmountPaid(),
				"Failed to release transaction, error: "+err.Error(),
				c.ViewUser.Uuid,
				nil,
			)
		}
	}

	redirectUrl := "/payments/" + c.ViewTransaction.Uuid
	util.RedirectOrAPIResponse(w, r, redirectUrl, c, c.IsAPIRequest)
}

func (c *Context) CancelTransaction(w web.ResponseWriter, r *web.Request) {
	model := c.ViewTransaction.DBModel()
	if model.IsCompleted() && !model.IsDispatched() && !model.IsShipped() {
		err := model.Cancel("User cancelled transaction", c.ViewUser.Uuid)
		if err != nil {
			model.SetTransactionStatus(
				model.CurrentPaymentStatus(),
				model.CurrentAmountPaid(),
				"Failed to cancel transaction",
				c.ViewUser.Uuid,
				nil,
			)
		}
	}

	redirectUrl := "/payments/" + c.ViewTransaction.Uuid
	util.RedirectOrAPIResponse(w, r, redirectUrl, c, c.IsAPIRequest)
}
