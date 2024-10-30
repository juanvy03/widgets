package api

import (
	"fmt"
	"net/http"
)

func LinkWidget(w http.ResponseWriter, req *http.Request) {
	fmt.Println(w, req)
}
