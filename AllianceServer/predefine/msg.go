package predefine

var MsgFlags = map[int]string{
	SUCCESS:                            "aa",
	NOT_LOGIN:                          "bb",
	REGISTER_FAILED:                    "cc",
	LOGIN_FAILED:                       "dd",
	GET_ALLIANCE_LIST_FALIED:           "ee",
	CREATE_ALLIANCE_FAILED:             "ff",
	JOIN_ALLIANCE_FALIED:               "gg",
	LEAVE_ALLIANCE_FALILED:             "hh",
	DESTROY_ALLIANCE_ITEM_FALIED:       "ww",
	TIDY_ALLIANCE_ITEM_FALIED:          "xx",
	COMMIT_ALLIANCE_ITEM_FALIED:        "zz",
	INCREASE_ALLIANCE_CAPACITY_FALILED: "qq",
}

func GetMsg(code int) string {
	retString, ok := MsgFlags[code]
	if !ok {
		return ""
	}

	return retString
}
