package request

type MessageSubscribe struct {
	Token  *string `header:"Authorization,omitempty" validate:"required"`
	ChatID *uint64 `param:"chat_id,omitempty" validate:"number,required"`
	UserID *uint64 `param:"id,omitempty" validate:"number,required"`
}

func (ms *MessageSubscribe) GetToken() string {
	return *ms.Token
}

func (ms *MessageSubscribe) SetToken(token string) {
	*ms.Token = token
}
