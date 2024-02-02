package request

type ChatCreate struct {
	Token     string `header:"Authorization,omitempty" validate:"required"`
	ContactID uint64 `json:"contact_id,omitempty" validate:"number,required"`
}

func (cc *ChatCreate) GetToken() string {
	return cc.Token
}

func (cc *ChatCreate) SetToken(token string) {
	cc.Token = token
}

type TokenAndChatID struct {
	ID    uint64 `param:"chat_id,omitempty" validate:"number,required"`
	Token string `header:"Authorization,omitempty" validate:"required"`
}

func (tc *TokenAndChatID) GetToken() string {
	return tc.Token
}

func (tc *TokenAndChatID) SetToken(token string) {
	tc.Token = token
}

type ChatDeleteMsg struct {
	ChatID uint64 `param:"chat_id,omitempty" validate:"number,required"`
	MsgID  uint64 `param:"message_id,omitempty" validate:"number,required"`
	Token  string `header:"Authorization,omitempty" validate:"required"`
}

func (cdm *ChatDeleteMsg) GetToken() string {
	return cdm.Token
}

func (cdm *ChatDeleteMsg) SetToken(token string) {
	cdm.Token = token
}
