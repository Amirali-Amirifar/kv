package api

type NodeRegisterHandlerRequest struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type NodeRegisterHandlerResponse struct {
	Error string `json:"error"`
}
