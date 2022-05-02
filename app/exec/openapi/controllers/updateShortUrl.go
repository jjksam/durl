package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"strconv"
)

type updateShortUrlReq struct {
	FullUrl        string `form:"fullUrl" valid:"Required"`
	IsFrozen       int    `form:"isFrozen"`
	ExpirationTime int    `form:"expirationTime"`
}

// UpdateShortUrl
// 函数名称: UpdateShortUrl
// 功能: 根据短链修改短链接信息
// 输入参数:
//	   fullUrl: 原始url
//	   isFrozen: 是否冻结
//	   expirationTime: 过期时间
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/18 5:46 下午 #
func (c *OpenApiController) UpdateShortUrl() {

	req := updateShortUrlReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	id := c.Ctx.Input.Param(":id")
	intId, _ := strconv.Atoi(id)

	// 查询此短链
	fields := map[string]interface{}{"id": intId}
	engine := db.NewDbService()
	urlInfo := engine.GetShortUrlInfo(fields)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 初始化需要更新的内容
	updateData := make(map[string]interface{})
	updateData["expiration_time"] = req.ExpirationTime
	updateData["full_url"] = req.FullUrl
	updateData["is_frozen"] = req.IsFrozen

	// 修改此短链信息
	_, err := engine.UpdateUrlById(intId, urlInfo.ShortNum, updateData)
	if err != nil {
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}

	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
