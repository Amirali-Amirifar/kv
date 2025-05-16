package service

import (
	"github.com/Amirali-Amirifar/kv/pkg/kvController/api"
)

type KvController struct {
}

func NewKvController() KvController {
	handler := api.KvRouteHandler{}
	router := api.SetupRouter(handler)
	router.Run(":8080")
	return KvController{}
}
