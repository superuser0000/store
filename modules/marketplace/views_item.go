package marketplace

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gocraft/web"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) viewShowItem(w web.ResponseWriter, r *web.Request) {
	packages := c.ViewItem.Item.PackagesWithoutReservation()
	c.ViewPackages = packages.ViewPackages()
	table := packages.GroupsTable()
	c.GroupPackagesByTypeOriginDestination = table.GetGroupPackagesByTypeOriginDestination()
}

func (c *Context) ViewShowItem(w web.ResponseWriter, r *web.Request) {
	c.viewShowItem(w, r)
	util.RenderTemplate(w, "item/show", c)
}

func (c *Context) EditItem(w web.ResponseWriter, r *web.Request) {
	c.ViewItemCategories = FindAllCategories().ViewCategories(c.ViewUser.Language)
	c.CategoryID = c.ViewItem.DBModel().ItemCategoryID
	c.ViewPackages = c.ViewItem.ViewPackages
	util.RenderTemplateOrAPIResponse(w, r, "item/edit", c, c.IsAPIRequest)
}

func (c *Context) viewItemSave(w web.ResponseWriter, r *web.Request) {

	item := c.ViewItem.DBModel()

	if r.PathParams["item"] == "new" {
		item.Uuid = util.GenerateUuid()
	}

	categoryId, err := strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		c.Error = err.Error()
		c.EditItem(w, r)
		return
	}

	category, err := FindCategoryByID(int(categoryId))
	if err != nil {
		c.Error = err.Error()
		c.EditItem(w, r)
		return
	}

	item.ItemCategory = *category
	item.Name = r.FormValue("name")
	item.Description = r.FormValue("description")
	if c.ViewUserStore != nil && item.StoreUuid == "" {
		item.StoreUuid = c.ViewUserStore.Uuid
		item.Store = *c.ViewUserStore.Store
	}

	validationError := item.Validate()
	if validationError != nil {
		c.Error = validationError.Error()
		c.EditItem(w, r)
		return
	}
	err = util.SaveImage(r, "image", 500, item.Uuid)
	if err != nil && r.PathParams["item"] == "new" {
		c.Error = "Image: " + err.Error()
		c.EditItem(w, r)
		return
	}
	err = item.Save()
	if err != nil {
		c.Error = err.Error()
		c.EditItem(w, r)
		return
	}

	if c.ViewUser.IsStaff || c.ViewUser.IsAdmin {
		now := time.Now()
		item.ReviewedByUserUuid = c.ViewUser.Uuid
		item.ReviewedAt = &now
		item.Save()
	}

	EventNewItem(item)
	viewItem := item.ViewItem(c.ViewUser.Language)
	c.ViewItem = &viewItem
}

func (c *Context) SaveItem(w web.ResponseWriter, r *web.Request) {
	c.viewItemSave(w, r)
	redirectUrl := "/store/" + c.ViewItem.Store.Storename + "/item/" + c.ViewItem.Uuid
	http.Redirect(w, r.Request, redirectUrl, 302)
}

func (c *Context) DeleteItem(w web.ResponseWriter, r *web.Request) {
	item := c.ViewItem.DBModel()
	item.Remove()
	util.RedirectOrAPIResponse(w, r, "/user/"+c.ViewUser.Username, c, c.IsAPIRequest)
}

func (c *Context) ViewItemImage(w web.ResponseWriter, r *web.Request) {
	itemUuid := r.PathParams["item"]
	size := "normal"
	if len(r.URL.Query()["size"]) > 0 {
		size = r.URL.Query()["size"][0]
	}
	util.ServeImage(itemUuid, size, w, r)
}

func (c *Context) ViewItemCategoryImage(w web.ResponseWriter, r *web.Request) {
	categoryID := r.PathParams["id"]
	size := "normal"
	if len(r.URL.Query()["size"]) > 0 {
		size = r.URL.Query()["size"][0]
	}
	// temp
	categoryID = "robocop"
	util.ServeImage(categoryID, size, w, r)
}
