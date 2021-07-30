package login_vo


type TwoAuthForm struct {
	TwoAuthCode string `form:"two_auth_code" binding:"required"`

}
