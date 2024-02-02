package request

type ContactCreate struct {
	Token       string `header:"Authorization,omitempty" validate:"required"`
	UserID      uint64 `param:"id,omitempty" validate:"number,required"`
	ContactID   uint64 `json:"contact_id,omitempty" validate:"number,required"`
	ContactName string `json:"contact_name,omitempty" validate:"required"`
}

func (cc *ContactCreate) GetToken() string {
	return cc.Token
}

func (cc *ContactCreate) SetToken(token string) {
	cc.Token = token
}

type ContactDelete struct {
	Token     string `header:"Authorization,omitempty" validate:"required"`
	UserID    uint64 `param:"id,omitempty" validate:"number,required"`
	ContactID uint64 `param:"contact_id,omitempty" validate:"number,required"`
}

func (cd *ContactDelete) GetToken() string {
	return cd.Token
}

func (cd *ContactDelete) SetToken(token string) {
	cd.Token = token
}
