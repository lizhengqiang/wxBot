package bot

import (
	"net/url"
	"strings"
)

func init() {

}

type Contact struct {
	VerifyFlag  int64
	UserName    string
	RemarkName  string
	NickName    string
	DisplayName string
}
type GetContactResponse struct {
	MemberList []*Contact
}

// 获取联系人列表
func (bot *WeixinBot) GetContact() error {
	mps := []*Contact{}
	groups := []*Contact{}
	contacts := []*Contact{}

	// SpecialUsers := []string{"newsapp", "fmessage", "filehelper", "weibo", "qqmail", "fmessage", "tmessage", "qmessage", "qqsync", "floatbottle", "lbsapp", "shakeapp", "medianote", "qqfriend", "readerapp", "blogapp", "facebookapp", "masssendapp", "meishiapp", "feedsapp", "voip", "blogappweixin", "weixin", "brandsessionholder", "weixinreminder", "wxid_novlwrv3lqwv11", "gh_22b87fa7cb3c", "officialaccounts", "notification_messages", "wxid_novlwrv3lqwv11", "gh_22b87fa7cb3c", "wxitil", "userexperience_alarm", "notification_messages"}
	response := &GetContactResponse{}
	queryValues := &url.Values{}
	queryValues.Add(passTicket, bot.getProperty(passTicket))
	queryValues.Add(skey, bot.getProperty(skey))
	queryValues.Add("r", bot.timestamp())
	// fmt.Sprintf("/webwxgetcontact?pass_ticket=%s&skey=%s&r=%s", bot.getProperty(passTicket), bot.getProperty(skey), bot.timestamp())
	u := "/webwxgetcontact?" + queryValues.Encode()
	bot.PostJson(u, &EmptyRequest{}, response)
	// 处理好友信息
	for _, contact := range response.MemberList {

		if contact.VerifyFlag != 0 {
			mps = append(mps, contact)
			continue
		}
		if strings.Contains(contact.UserName, "@@") {
			groups = append(groups, contact)
			continue
		}
		if strings.Contains(contact.UserName, "@") {
			contacts = append(contacts, contact)
			continue
		}
	}
	bot.marshal(memberList, response.MemberList)
	bot.marshal(groupList, groups)
	bot.marshal(mpList, mps)
	bot.Println(mps)
	bot.Println(bot.Get(mpList))
	bot.marshal(contactList, contacts)
	bot.Println("联系人信息获取完毕")
	return nil
}

func (bot *WeixinBot) GetRemarkName(id string) (name string) {
	contacts := []*Contact{}
	bot.unmarshal(memberList, &contacts)
	for _, contact := range contacts {
		if contact.UserName == id {
			if contact.RemarkName == "" {
				return contact.NickName
			} else {
				return contact.RemarkName
			}
		}
	}
	return "未知"
}

