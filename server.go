package main

import (
	main_router "Redbus_backend/Router"
	"fmt"
	"net/http"
)

func main() {
	router := main_router.GetRouter()
	fmt.Println("Server is running : 3000")
	http.ListenAndServe(":3000", router)
}
