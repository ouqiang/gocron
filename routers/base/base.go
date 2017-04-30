package base

import (
    "gopkg.in/macaron.v1"
    "github.com/ouqiang/gocron/models"
)

func ParsePageAndPageSize(ctx *macaron.Context, params models.CommonMap)  {
    page := ctx.QueryInt("page")
    pageSize := ctx.QueryInt("page_size")
    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = models.PageSize
    }

    params["Page"] = page
    params["PageSize"] = pageSize
}