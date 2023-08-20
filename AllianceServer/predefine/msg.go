package predefine

var MsgFlags = map[int]string{
	SUCCESS:                         "成功",
	NotLogin:                        "当前未登录",
	RegisterFailed:                  "注册失败",
	LoginFailed:                     "登录失败",
	GetAllianceListFalied:           "获取公会列表失败",
	CreateAllianceFailed:            "创建公会失败",
	JoinAllianceFalied:              "加入公会失败",
	LeaveAllianceFaliled:            "解散公会失败",
	DestroyAllianceItemFalied:       "销毁公会物品失败",
	TidyAllianceItemFalied:          "整理公会物品失败",
	CommitAllianceItemFalied:        "提交公会物品失败",
	IncreaseAllianceCapacityFaliled: "扩容公会仓库失败",
	CommonFailed:                    "通用失败",
}

func GetMsg(code int) string {
	retString, ok := MsgFlags[code]
	if !ok {
		return ""
	}

	return retString
}
