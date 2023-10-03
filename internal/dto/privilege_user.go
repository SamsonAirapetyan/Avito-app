package dto

type PrivilegedUserCreateDTO struct {
	UserID        int      `json:"user_id"`
	PrivilegeList []string `json:"add_privilege"`
}

type PrivilegedUserDTO struct {
	UserID    int    `json:"user_id"`
	Privilege string `json:"privilege_title" validate:"required,uppercase,max=20"`
}

type PrivilegedUserDeleteDTO struct {
	UserID        int      `json:"user_id"`
	PrivilegeList []string `json:"del_privilege"`
}
