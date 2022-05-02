package controllers

import (
	"durl/app/share/comm"
	"durl/app/share/dao/db"
)

type BatchFrozenShortUrlReq struct {
	Ids      []int `from:"ids" valid:"Required"`
	IsFrozen int   `from:"isFrozen"`
}

type BatchFrozenShortUrlRes struct {
	RequestCount int   `json:"requestCount"`
	UpdateCount  int   `json:"updateCount"`
	ErrIds       []int `json:"errIds"`
}

// BatchFrozenShortUrl
// 函数名称: BatchFrozenShortUrl
// 功能: 批量冻结/解冻Url
// 输入参数:
//		BatchFrozenShortUrlReq{}
// 输出参数:
// 返回: BatchFrozenShortUrlRes{}
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/11/26 2:15 下午 #
func (c *BackendController) BatchFrozenShortUrl() {

	req := BatchFrozenShortUrlReq{}

	c.BaseCheckParams(&req)

	// 查询待操作Url信息
	fields := map[string]interface{}{"id": req.Ids}
	engine := db.NewDbService()
	data := engine.GetAllShortUrl(fields)
	if data == nil {
		c.ErrorMessage(comm.ErrNotFound, comm.MsgParseFormErr)
		return
	}

	var updateIds []int
	errIds := make([]int, 0)
	var insertShortNum []int
	// 提交id数量与查询出的数据量不一致
	// 需要以数据库数据为准筛选出差集，准备进行错误返回
	requestCount := len(req.Ids)
	updateCount := len(data)
	if updateCount != requestCount {

		// 将请求操作的id 提为key
		mapData := make(map[int]interface{})
		for _, v := range data {
			mapData[v.Id] = v.ShortNum
		}

		for _, v := range req.Ids {
			if mapData[v] != nil {
				updateIds = append(updateIds, v)
				insertShortNum = append(insertShortNum, mapData[v].(int))
			} else {
				errIds = append(errIds, v)
			}
		}

	} else {
		updateIds = req.Ids
		for _, vv := range data {
			insertShortNum = append(insertShortNum, vv.ShortNum)
		}
	}

	// 正确数据进行批量操作
	// 批量冻结/解冻Url
	updateData := map[string]interface{}{"is_frozen": req.IsFrozen}
	updateWhere := map[string]interface{}{"id": updateIds}

	_, err := engine.BatchUpdateUrlByIds(updateWhere, insertShortNum, updateData)
	if err != nil {
		c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, &BatchFrozenShortUrlRes{
			RequestCount: requestCount,
			UpdateCount:  0,
			ErrIds:       req.Ids,
		})
		return
	}
	res := BatchFrozenShortUrlRes{
		RequestCount: requestCount,
		UpdateCount:  updateCount,
		ErrIds:       errIds,
	}
	c.FormatInterfaceResp(comm.OK, comm.OK, comm.MsgOk, &res)
	return
}
