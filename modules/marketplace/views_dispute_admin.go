package marketplace

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) AdminDisputeList(w web.ResponseWriter, r *web.Request) {
	// transaction type
	if len(r.URL.Query()["status"]) > 0 {
		c.SelectedStatus = r.URL.Query()["status"][0]
	}
	// pages
	pageSize := 20
	if len(r.URL.Query()["page"]) > 0 {
		strPage := r.URL.Query()["page"][0]
		page, err := strconv.ParseInt(strPage, 10, 32)
		if err != nil || page < 0 {
			http.NotFound(w, r.Request)
			return
		}
		c.Page = int(page)
	} else {
		c.Page = 1
	}
	c.NumberOfTransactions = CountDisputedTransactions("", c.SelectedStatus)
	c.NumberOfPages = int(math.Ceil(float64(c.NumberOfTransactions) / float64(pageSize)))
	for i := 0; i < c.NumberOfPages; i++ { // paging
		c.Pages = append(c.Pages, i+1)
	}

	transactions := Transactions(GetDisputedTransactionsPaged(pageSize, c.Page-1, "", c.SelectedStatus))
	c.ViewTransactions = transactions.ViewTransactions()
	util.RenderTemplate(w, "dispute/admin/list", c)
}

func (c *Context) ViewDisputePartialRefund(w web.ResponseWriter, r *web.Request) {

	if !(c.ViewUser.IsAdmin || c.ViewUser.IsStaff) {
		http.Error(w, "User is not admin or staff", 502)
		return
	}

	refundPercent, err := strconv.ParseInt(r.FormValue("percent"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	err = c.Dispute.PartialRefund(c.ViewUser.Uuid, float64(refundPercent)/100.)
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	redirectUrl := fmt.Sprintf("/dispute/%s", c.Dispute.Uuid)
	http.Redirect(w, r.Request, redirectUrl, 302)
	return
}
