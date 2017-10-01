package controllers

import (
	"github.com/jcherianucla/bloggo/models"
	"net/http"
)

var GetPosts = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Content-Type", "application/json")
	},
)
