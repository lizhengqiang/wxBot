package bot

func (this *WeixinBot) FindMp(username string) *Contact {
	mps := []*Contact{}
	this.unmarshal(mpList, &mps)

	for _, contact := range mps {
		if contact.UserName == username {
			return contact
		}
	}
	return nil
}
