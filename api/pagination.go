package api

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func parsePagination(c echo.Context) (uint64, uint64, uint64) {
	limit, _ := strconv.ParseUint(c.QueryParam("limit"), 10, 64)
	if limit == 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	page, _ := strconv.ParseUint(c.QueryParam("page"), 10, 64)
	var offset uint64 = 0
	if page == 0 {
		page = 1
	}
	offset = (page - 1) * limit

	return limit, page, offset
}
