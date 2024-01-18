package vo

type OperationLogListRequest struct {
	Username string `json:"username" form:"username"`
	Ip       string `json:"ip" form:"ip"`
	Path     string `json:"path" form:"path"`
	Status   int    `json:"status" form:"status"`
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}

type DeleteOperationLogRequest struct {
	OperationLogIds []uint `json:"operationLogIds" form:"operationLogIds"`
}
