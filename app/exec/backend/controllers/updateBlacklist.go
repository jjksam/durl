package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
	"strconv"
)

type updateBlacklistReq struct {
	Ip string `valid:"IP"`
}

// UpdateBlacklist
// 函数名称: UpdateBlacklist
// 功能: 根据短链修改短链接信息
// 输入参数:
//	   fullUrl: 原始url
//	   isFrozen: 是否冻结
//	   expirationTime: 过期时间
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021/12/07 5:44 下午 #
func (c *BackendController) UpdateBlacklist() {

	req := updateBlacklistReq{}
	// 效验请求参数格式
	c.BaseCheckParams(&req)

	id := c.Ctx.Input.Param(":id")
	intId, _ := strconv.Atoi(id)

	// 查询此短链
	fields := map[string]interface{}{"id": intId}
	engine := db.NewDbService()
	urlInfo := engine.GetBlacklistInfo(fields)
	if urlInfo.Id == 0 {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	// 初始化需要更新的内容
	updateData := make(map[string]interface{})
	updateData["ip"] = req.Ip

	// 统计结果总条数
	total := engine.GetBlacklistListTotal(updateData)
	if total > 0 {
		c.ErrorMessage(comm.ErrRepeatCommit, comm.MsgRepeatCommitErr)
		return
	}

	// 修改此短链信息
	_, err := engine.UpdateBlacklistById(intId, updateData)
	if err != nil {
		c.ErrorMessage(comm.ErrSysDb, comm.MsgNotOk)
		return
	}

	c.FormatResp(comm.OK, comm.OK, comm.MsgOk)
	return
}
