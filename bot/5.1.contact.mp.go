package bot

func (bot *WeixinBot) FindMp(username string) *Contact {
	mps := []*Contact{}
	bot.unmarshal(mpList, &mps)

	for _, contact := range mps {
		if contact.UserName == username {
			return contact
		}
	}
	return nil
}
