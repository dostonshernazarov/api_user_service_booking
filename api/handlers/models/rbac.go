package models

type Policy struct {
	Role     string `json:"role"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

type ListPolicyResponse struct {
	Policies []*Policy `json:"policies"`
}

type CreateNewRoleRequest struct {
	UserId string `json:"userId"`
	Path   string `json:"path"`
	Role   string `json:"role"`
}
