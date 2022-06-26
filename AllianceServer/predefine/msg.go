package predefine

var MsgFlags = map[int]string{
	SUCCESS:                            "成功",
	NOT_LOGIN:                          "当前未登录",
	REGISTER_FAILED:                    "注册失败",
	LOGIN_FAILED:                       "登录失败",
	GET_ALLIANCE_LIST_FALIED:           "获取公会列表失败",
	CREATE_ALLIANCE_FAILED:             "创建公会失败",
	JOIN_ALLIANCE_FALIED:               "加入公会失败",
	LEAVE_ALLIANCE_FALILED:             "解散公会失败",
	DESTROY_ALLIANCE_ITEM_FALIED:       "销毁公会物品失败",
	TIDY_ALLIANCE_ITEM_FALIED:          "整理公会物品失败",
	COMMIT_ALLIANCE_ITEM_FALIED:        "提交公会物品失败",
	INCREASE_ALLIANCE_CAPACITY_FALILED: "扩容公会仓库失败",
}

func GetMsg(code int) string {
	retString, ok := MsgFlags[code]
	if !ok {
		return ""
	}

	return retString
}
