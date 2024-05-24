package handlers

import (
	"akari/handlers/authHandler"
	"akari/models"
)

var (
	Services models.HandlerMap = models.HandlerMap{
		authHandler.Service: authHandler.HandleRequest,
	}
)
