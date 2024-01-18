package vo

type RegisterAndLoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Username     string `form:"username" json:"username" validate:"required,min=2,max=20"`
	Password     string `form:"password" json:"password"`
	Mobile       string `form:"mobile" json:"mobile" validate:"required,checkMobile"`
	Avatar       string `form:"avatar" json:"avatar"`
	Nickname     string `form:"nickname" json:"nickname" validate:"min=0,max=20"`
	Introduction string `form:"introduction" json:"introduction" validate:"min=0,max=255"`
	Status       uint   `form:"status" json:"status" validate:"oneof=1 2"`
	RoleIds      []uint `form:"roleIds" json:"roleIds" validate:"required"`
}

type UserListRequest struct {
	Username string `json:"username" form:"username" `
	Mobile   string `json:"mobile" form:"mobile" `
	Nickname string `json:"nickname" form:"nickname" `
	Status   uint   `json:"status" form:"status" `
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

type DeleteUserRequest struct {
	UserIds []uint `json:"userIds" form:"userIds"`
}

type ChangePwdRequest struct {
	OldPassword string `json:"oldPassword" form:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" form:"newPassword" validate:"required"`
}
