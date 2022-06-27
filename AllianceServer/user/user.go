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

// 新建连接
func NewConnection(conn *net.TCPConn) *Conn {
	return &Conn{
		conn:        conn,
		isLogin:     false,
		userAccount: "",
	}
}

// 处理连接
func (c *Conn) DealConnection() {
	defer func() {
		_ = c.conn.Close()
	}()
	for true {
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
		retContent, _ := c.dispatch(content) // 指令分发
		retLength := len(retContent)
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, int32(retLength)) // 执行结果返回
		if err != nil {
			log.Println("binary.Write failed:", err)
		}
		_, _ = c.conn.Write(append(buf.Bytes(), []byte(retContent)...))
	}
}

// 消息分发
func (c *Conn) dispatch(content string) (string, int) {
	stringSlice := util.DealString(content)
	code := 0
	var retString string
	length := len(stringSlice)
	if !c.isLogin && !(stringSlice[0] == "login" || stringSlice[0] == "register") {
		code = predefine.NOT_LOGIN
	} else {
		switch stringSlice[0] {
		case "register": // 注册并登录
			if length != 3 {
				code = 2
			} else {
				err := api.Register(stringSlice[1], stringSlice[2])
				if err != nil {
					log.Println("register failed, err:", err)
					code = 2
				}
				err = api.Login(stringSlice[1], stringSlice[2])
				if err != nil {
					log.Println("login failed, err:", err)
					code = 3
				}
				c.isLogin = true
				c.userAccount = stringSlice[1]
				retString = "注册并登录成功"
			}
		case "login": // 登录
			if length != 3 {
				code = 3
			} else {
				err := api.Login(stringSlice[1], stringSlice[2])
				if err != nil {
					log.Println("login failed, err:", err)
					code = 3
				}
				c.isLogin = true
				c.userAccount = stringSlice[1]
				retString = "登录成功"
			}
		case "allianceList": // 查看已有公会列表
			ret, err := api.AllianceList()
			if err != nil {
				log.Println("get alliance list failed, err:", err)
				code = 4
			}
			retString = ret
		case "createAlliance": // 创建公会
			if length != 2 {
				code = 5
			} else {
				err := api.CreateAlliance(c.userAccount, stringSlice[1])
				if err != nil {
					log.Println("create alliance failed, err:", err)
					code = 5
				}
				retString = "创建公会成功"
			}
		case "joinAlliance": // 加入公会
			if length != 2 {
				code = 6
			} else {
				err := api.JoinAlliance(c.userAccount, stringSlice[1])
				if err != nil {
					log.Println("join alliance failed, err:", err)
					code = 6
				}
				retString = "加入公会成功"
			}
		case "dismissAlliance": // 解散公会
			err := api.DismissAlliance(c.userAccount)
			if err != nil {
				log.Println("leave alliance failed, err:", err)
				code = 7
			}
			retString = "解散公会成功"
		case "destroyItem": // 销毁公会物品
			if length != 2 {
				code = 8
			} else {
				num, _ := strconv.Atoi(stringSlice[1])
				err := api.DestroyItem(c.userAccount, num)
				if err != nil {
					log.Println("destroy alliance item failed, err:", err)
					code = 8
				}
				retString = "删除物品成功"
			}
		case "clearup": // 整理公会物品
			err := api.TidyItems(c.userAccount)
			if err != nil {
				log.Println("tidy alliance item failed, err:", err)
				code = 9
			}
			retString = "clearup alliance ok!"
		case "storeItem": // 提交公会物品
			if length != 3 {
				code = 10
			} else {
				id, _ := strconv.Atoi(stringSlice[1])
				num, _ := strconv.Atoi(stringSlice[2])
				err := api.CommitItem(c.userAccount, id, num)
				if err != nil {
					log.Println("commit item failed, err:", err)
					code = 10
				}
				retString = "提交物品成功"
			}
		case "increaseCapacity": // 扩容公会仓库
			err := api.IncreaseCapacity(c.userAccount)
			if err != nil {
				log.Println("increase capacity failed, err:", err)
				code = 11
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
		default:
			retString = "no command"
		}
	}

	if code != 0 {
		retString = predefine.GetMsg(code)
	}
	return retString, code
}
