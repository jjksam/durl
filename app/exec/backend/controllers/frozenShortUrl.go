package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"strconv"
)

// 函数名称: FrozenShortUrl
// 功能: 冻结ShortUrl
// 输入参数:
//     id: 数据id
// 输出参数:
// 返回: 冻结操作结果
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/26 1:56 下午 #

func (c *BackendController) FrozenShortUrl() {

	id := c.Ctx.Input.Param(":id")
	intId, _ := strconv.ParseUint(id, 10, 32)
	uint32Id := uint32(intId)

	// 查询此短链
	fields := map[string]interface{}{"id": uint32Id}
	engine := db.NewDbService()
	urlInfo := engine.GetShortUrlInfo(fields)
	if urlInfo.ShortNum == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 冻结/解冻ShortUrl
	if urlInfo.IsFrozen == 1 {
		c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
		return
	}

	updateData := make(map[string]interface{})
	updateData["is_frozen"] = 0
	_, err := engine.UpdateUrlById(uint32Id, urlInfo.ShortNum, updateData)
	if err != nil {
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}
	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
