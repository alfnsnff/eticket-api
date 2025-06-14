package response

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Params struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`   // e.g. "name:asc"
	Search string `json:"search"` // e.g. "john"
	Path   string `json:"path"`   // current request path
}

func GetParams(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "")
	search := c.DefaultQuery("search", "")
	path := c.Request.URL.Path

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	return Params{
		Page:   page,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Search: search,
		Path:   path,
	}
}
