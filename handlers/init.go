package handlers

import (
	customer "umbrellacorp/handlers/customer"
)

// Init initializes all entity handlers
func Init() {
	customer.Init()
}
