package request

type AccountCreate struct {
	Firstname *string `json:"firstname,omitempty" validate:"required"`
	Phone     *string `json:"phone,omitempty" validate:"e164,required"`
	Username  *string `json:"username,omitempty" validate:"required"`
	Password  *string `json:"password,omitempty" validate:"min=8,required"`

	Lastname  *string `json:"lastname,omitempty"`
	ImagePath *string `json:"img_path,omitempty" validate:"url"`
	Bio       *string `json:"bio,omitempty" validate:"max=100"`
}

type AccountLoginByUsername struct {
	Username *string `json:"username,omitempty" validate:"required"`
	Password *string `json:"password,omitempty" validate:"required"`
}

type AccountLoginByPhone struct {
	Phone    *string `json:"phone,omitempty" validate:"e164,required"`
	Password *string `json:"password,omitempty" validate:"required"`
}

type AccountLoginWholeData struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

type AccountUpdate struct {
	ID    *uint64 `param:"id,omitempty" validate:"number,required"`
	Token *string `header:"Authorization,omitempty" validate:"required"`

	Firstname *string `json:"firstname,omitempty"`
	Phone     *string `json:"phone,omitempty" validate:"e164"`
	Username  *string `json:"username,omitempty"`
	Password  *string `json:"password,omitempty" validate:"min=8"`
	Lastname  *string `json:"lastname,omitempty"`
	ImagePath *string `json:"img_path,omitempty" validate:"url"`
	Bio       *string `json:"bio,omitempty" validate:"max=100"`
}

func (au *AccountUpdate) GetToken() string {
	return *au.Token
}

func (au *AccountUpdate) SetToken(token string) {
	*au.Token = token
}

type AccountSearch struct {
	Token *string `header:"Authorization,omitempty" validate:"required"`

	ID        *uint64 `query:"id,omitempty" validate:"number"`
	Firstname *string `query:"firstname,omitempty"`
	Lastname  *string `query:"lastname,omitempty"`
	Phone     *string `query:"phone,omitempty" validate:"e164"`
	Username  *string `query:"username,omitempty"`
}

func (as *AccountSearch) GetToken() string {
	return *as.Token
}

func (as *AccountSearch) SetToken(token string) {
	*as.Token = token
}
