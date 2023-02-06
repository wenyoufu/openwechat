package openwechat

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

//func TestLogin(t *testing.T) {
//	bot := DefaultBot(Desktop)
//	bot.LoginCallBack = func(body []byte) {
//		t.Log("login success")
//	}
//	if err := bot.Login(); err != nil {
//		t.Error(err)
//	}
//}

//func TestLogout(t *testing.T) {
//	bot := DefaultBot(Desktop)
//	bot.LoginCallBack = func(body []byte) {
//		t.Log("login success")
//	}
//	bot.LogoutCallBack = func(bot *Bot) {
//		t.Log("logout")
//	}
//	bot.MessageHandler = func(msg *Message) {
//		if msg.IsText() && msg.Content == "logout" {
//			bot.Logout()
//		}
//	}
//	if err := bot.Login(); err != nil {
//		t.Error(err)
//		return
//	}
//	bot.Block()
//}

func TestMessageHandle(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.MessageHandler = func(msg *Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	bot.Block()
}

func TestFriends(t *testing.T) {
	bot := DefaultBot(Desktop)
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	friends, err := user.Friends()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(friends)
}

func TestGroups(t *testing.T) {
	bot := DefaultBot(Desktop)
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	groups, err := user.Groups()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(groups)
}

func TestPinUser(t *testing.T) {
	bot := DefaultBot(Desktop)
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	friends, err := user.Friends()
	if err != nil {
		t.Error(err)
		return
	}
	if friends.Count() > 0 {
		f := friends.First()
		f.Pin()
		time.Sleep(time.Second * 5)
		f.UnPin()
	}
}

func TestSender(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.MessageHandler = func(msg *Message) {
		if msg.IsSendByGroup() {
			fmt.Println(msg.SenderInGroup())
		} else {
			fmt.Println(msg.Sender())
		}
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	bot.Block()
}

// TestGetUUID
// @description: 获取登录二维码(UUID)
// @param t
func TestGetUUID(t *testing.T) {
	bot := DefaultBot(Desktop)

	uuid, err := bot.Caller.GetLoginUUID()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(uuid)
}

// TestLoginWithUUID
// @description: 使用UUID登录
// @param t
func TestLoginWithUUID(t *testing.T) {
	uuid := "oZZsO0Qv8Q=="
	bot := DefaultBot(Desktop)
	bot.SetUUID(uuid)
	err := bot.Login()
	if err != nil {
		t.Errorf("登录失败: %v", err.Error())
		return
	}
}

func TestBBB(t *testing.T) {
	//bot := DefaultBot()
	bot := DefaultBot(Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
		if msg.IsSendByGroup() {
			// 首先尝试从缓存里面查找, 如果没有找到则从服务器获取
			members, err := msg.bot.self.Members()
			if err != nil {
				return
			}
			user, exist := members.GetByUserName(msg.FromUserName)
			if strings.Contains(user.NickName, "幸福合作") {
				if !exist {
					// 找不到, 从服务器获取
					user = &User{self: msg.bot.self, UserName: msg.FromUserName}
					err = user.Detail()
				}
				user2, exist := members.GetByUserName(msg.ToUserName)
				if !exist {
					// 找不到, 从服务器获取
					user2 = &User{self: msg.bot.self, UserName: msg.ToUserName}
					err = user2.Detail()
				}
				user3, _ := msg.SenderInGroup()
				fmt.Println(user3)
				//if err == nil /* && strings.Contains(sender.NickName, "四季平安")*/ {
				msg.ReplyText("???")
				//fmt.Println("test!!!", sender.NickName)
			}
			//}
		}
	}

	// 注册登陆二维码回调
	bot.UUIDCallback = PrintlnQrcodeUrl
	//bot.ScanCallBack = PrintlnScanCode
	//bot.LoginCallBack = PrintlnLogin
	// 登陆
	//if err := bot.Login(); err != nil {
	//	fmt.Println(err)
	//	return
	//}
	// 创建热存储容器对象
	reloadStorage := NewFileHotReloadStorage("storage.json")

	defer reloadStorage.Close()
	if err := bot.HotLogin(reloadStorage, NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return
	}
	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)
	fmt.Println("好友数量为：", friends.Count())
	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)
	fmt.Println("群主数量为：", groups.Count())
	mps, err := self.Mps()
	fmt.Println(mps, err)
	fmt.Println("公众号数量为：", mps.Count())
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
