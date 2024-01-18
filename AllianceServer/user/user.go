package user

import (
	"AllianceServer/api"
	"AllianceServer/predefine"
	"AllianceServer/util"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"strconv"
)

type Conn struct {
	conn        *net.TCPConn // 连接handler
	isLogin     bool         // 是否登录标识
	userAccount string       // 玩家账号
}

// NewConnection 新建连接
func NewConnection(conn *net.TCPConn) *Conn {
	return &Conn{
		conn:        conn,
		isLogin:     false,
		userAccount: "",
	}
}

// DealConnection 处理连接
func (c *Conn) DealConnection() {
	defer func() {
		_ = c.conn.Close()
	}()
	for {
		length := make([]byte, 4) // 长度的字节数固定为4
		if _, err := io.ReadFull(c.conn, length); err != nil {
			return
		} // 先拿出表示长度的数据
		realLength := binary.LittleEndian.Uint32(length)
		data := make([]byte, realLength)
		if _, err := io.ReadFull(c.conn, data); err != nil {
			return
		} // 读取真正的数据
		content := string(data)
		log.Println("receive data: ", content)
		retContent, _ := c.Dispatch(content) // 指令分发
		retLength := len(retContent)
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, int32(retLength)) // 执行结果返回
		if err != nil {
			log.Println("binary.Write failed:", err)
		}
		_, _ = c.conn.Write(append(buf.Bytes(), []byte(retContent)...))
	}
}

// Dispatch 消息分发
func (c *Conn) Dispatch(content string) (string, int) {
	stringSlice := util.DealString(content)
	code := predefine.SUCCESS
	var retString string
	length := len(stringSlice)
	if !c.isLogin && !(stringSlice[0] == "login" || stringSlice[0] == "register") {
		code = predefine.NotLogin
	} else {
		switch stringSlice[0] {
		case "register": // 注册并登录
			retString, code = register(content)
			c.isLogin = true
			c.userAccount = stringSlice[1]
		case "login": // 登录
			retString, code = login(content)
			c.isLogin = true
			c.userAccount = stringSlice[1]
		case "allianceList": // 查看已有公会列表
			retString, code = allianceList()
		case "createAlliance": // 创建公会
			if length != 2 {
				code = predefine.CreateAllianceFailed
			} else {
				err := api.CreateAlliance(c.userAccount, stringSlice[1])
				if err != nil {
					log.Println("create alliance failed, err:", err)
					code = predefine.GetAllianceListFalied
				}
				retString = "创建公会成功"
			}
		case "joinAlliance": // 加入公会
			if length != 2 {
				code = predefine.JoinAllianceFalied
			} else {
				err := api.JoinAlliance(c.userAccount, stringSlice[1])
				if err != nil {
					log.Println("join alliance failed, err:", err)
					code = predefine.JoinAllianceFalied
				}

				retString = "加入公会成功"
			}
		case "dismissAlliance": // 解散公会
			err := api.DismissAlliance(c.userAccount)
			if err != nil {
				log.Println("leave alliance failed, err:", err)
				code = predefine.LeaveAllianceFaliled
			}

			retString = "解散公会成功"
		case "destroyItem": // 销毁公会物品
			if length != 2 {
				code = predefine.DestroyAllianceItemFalied
			} else {
				num, _ := strconv.Atoi(stringSlice[1])
				err := api.DestroyItem(c.userAccount, num)
				if err != nil {
					log.Println("destroy alliance item failed, err:", err)
					code = predefine.DestroyAllianceItemFalied
				}

				retString = "删除物品成功"
			}
		case "clearup": // 整理公会物品
			err := api.TidyItems(c.userAccount)
			if err != nil {
				log.Println("tidy alliance item failed, err:", err)
				code = predefine.TidyAllianceItemFalied
			}

			retString = "clearup alliance ok!"
		case "storeItem": // 提交公会物品
			if length != 3 {
				code = predefine.CommitAllianceItemFalied
			} else {
				id, _ := strconv.Atoi(stringSlice[1])
				num, _ := strconv.Atoi(stringSlice[2])
				err := api.CommitItem(c.userAccount, id, num)
				if err != nil {
					log.Println("commit item failed, err:", err)
					code = predefine.CommitAllianceItemFalied
				}

				retString = "提交物品成功"
			}
		case "increaseCapacity": // 扩容公会仓库
			err := api.IncreaseCapacity(c.userAccount)
			if err != nil {
				log.Println("increase capacity failed, err:", err)
				code = predefine.IncreaseAllianceCapacityFaliled
			}

			retString = "扩容公会仓库成功"
		case "whichAlliance": // 查看自己的公会信息
			ret, err := api.WhichAlliance(c.userAccount)
			if err != nil {
				log.Println("which alliance failed, err:", err)
			}

			retString = ret
		case "allianceItems": // 查看公会物品信息
			ret, err := api.AllianceItems(c.userAccount)
			if err != nil {
				log.Println("alliance items failed, err:", err)
			}

			retString = ret
		default: // 默认情况不做任何处理
			retString = "no command"
		}
	}

	if code != predefine.SUCCESS {
		retString = predefine.GetMsg(code)
	}

	return retString, code
}

func register(content string) (string, int) {
	stringSlice := util.DealString(content)
	length := len(stringSlice)
	var code int
	var retString string

	if length != 3 {
		code = predefine.RegisterFailed
	} else {
		err := api.Register(stringSlice[1], stringSlice[2])
		if err != nil {
			log.Println("register failed, err:", err)
			code = predefine.RegisterFailed
		}
		err = api.Login(stringSlice[1], stringSlice[2])
		if err != nil {
			log.Println("login failed, err:", err)
			code = predefine.LoginFailed
		}
		retString = "注册并登录成功"
	}

	return retString, code
}

// login 登陆
func login(content string) (string, int) {
	stringSlice := util.DealString(content)
	length := len(stringSlice)
	var code int
	var retString string

	if length != 3 {
		code = predefine.LoginFailed
	} else {
		err := api.Login(stringSlice[1], stringSlice[2])
		if err != nil {
			log.Println("login failed, err:", err)
			code = predefine.LoginFailed
		}
		retString = "登录成功"
	}

	return retString, code
}

// allianceList 联盟列表
func allianceList() (string, int) {
	var code int
	var retString string
	ret, err := api.AllianceList()
	if err != nil {
		log.Println("get alliance list failed, err:", err)
		code = predefine.GetAllianceListFalied
	}

	retString = ret

	return retString, code

}

// createAlliance 创建联盟
func createAlliance() {

}

// joinAlliance 加入联盟
func joinAlliance() {

}

// 离开联盟
func leaveAlliance() {

}

// dismissAlliance 解散联盟
func dismissAlliance() {

}

// destroyItem 销毁物品
func destroyItem() {

}

// clearup 清理仓库
func clearup() {

}

// storeItem 存储物品
func storeItem() {

}

// increaseCapacity 仓库扩容
func increaseCapacity() {

}

// whichAlliance 查看联盟信息
func whichAlliance() {

}

// allianceItems 查看联盟物品
func allianceItems() {

}
