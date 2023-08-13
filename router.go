package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route -- HTTP route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// NewRouter -- the HTTP router
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}

// Routes -- List of active routes
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	{
		"Members",
		"GET",
		"/members",
		Members,
	},
	{
		"MembersS",
		"GET",
		"/memberss",
		MembersS,
	},
	Route{
		"EditMember",
		"GET",
		"/editmember/{memberId}",
		EditMember,
	},
	Route{
		"SaveMember",
		"POST",
		"/savemember/{memberId}",
		SaveMember,
	},
	Route{
		"Login",
		"POST",
		"/login",
		Login,
	},
	Route{
		"NotAdminUser",
		"GET",
		"/notadminuser",
		NotAdminUser,
	},
}
