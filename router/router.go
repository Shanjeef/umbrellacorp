package router

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter returns a configured gorilla mux router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, routes := range routesRegistry {
		for _, route := range routes {
			router.Name(route.Name).Methods(route.Methods...).Path(route.Path).Handler(handle(route.HandlerFunc))
		}
	}
	return router
}

// HandlerFunc is umbrellaCorp's handler fn signature that is decorated with http.HandlerFunc
type HandlerFunc func(Request) (Response, error)

// Route defines details to associate a handler with an http router. See router.handle(HandlerFunc) fn for more details
type Route struct {
	Name        string
	Methods     []string
	Path        string
	HandlerFunc HandlerFunc
}

// Routes is a list of Route objects
type Routes []Route

// Contains returns a matching Route in receiver list based on the parameterized route's Name
func (routes Routes) Contains(r Route) *Route {
	for i, route := range routes {
		if route.Name == r.Name {
			return &routes[i]
		}
	}
	return nil
}

var routesRegistry map[string]Routes

// RegisterRoutes is a utility function to associate a list of routes with an entity. If an existing method and path has already been registered
// an error is returned
func RegisterRoutes(entity string, routes Routes) error {
	if routesRegistry == nil {
		routesRegistry = map[string]Routes{}
	}
	existing := routesRegistry[entity]
	for _, route := range routes {
		if existingRoute := existing.Contains(route); existingRoute != nil {
			return fmt.Errorf("Route %s already registered with path: %s, method: %v", route.Name, existingRoute.Path, existingRoute.Methods)
		}
		existing = append(existing, route)
	}

	routesRegistry[entity] = existing
	return nil
}

func handle(handlerFn HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Read up to 1 MB of data from the client
		body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1000000))
		if err != nil {
			err = fmt.Errorf("Failed to read request body. Err: %v", err.Error())
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		defer req.Body.Close()

		request := Request{}
		if len(body) > 0 {
			err = json.Unmarshal(body, &request.Info)
			if err != nil {
				err = fmt.Errorf("Failed to unmarshal request body. Err: %v", err.Error())
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
		}

		resp, err := handlerFn(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(resp.Info); err != nil {
			err = fmt.Errorf("Failed to marshal response details. Err: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
