package api

import (
	"encoding/json"
	"net/http"

	"github.com/bigkevmcd/peanut/pkg/http/config"
	"github.com/bigkevmcd/peanut/pkg/parser"
	"github.com/go-git/go-git/v5"
	"github.com/julienschmidt/httprouter"
)

type gitParser func(path, opts *git.CloneOptions) *parser.Config

// APIRouter is an HTTP API for accessing app configurations.
type APIRouter struct {
	*httprouter.Router
	cfg    *config.Config
	parser gitParser
}

// ListApps returns the list of configured apps.
func (a *APIRouter) ListApps(w http.ResponseWriter, r *http.Request) {
	result := listAppsResponse{Apps: []appResponse{}}
	for _, v := range a.cfg.Apps {
		result.Apps = append(result.Apps, appResponse{Name: v.Name})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetApp returns a specific app.
func (a *APIRouter) GetApp(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	app := a.cfg.App(params.ByName("name"))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

// GetEnvironment returns a specific app.
func (a *APIRouter) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	env := a.cfg.App(params.ByName("name")).Environment(params.ByName("env"))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(env)
}

// NewRouter creates and returns a new APIRouter.
func NewRouter(cfg *config.Config) *APIRouter {
	api := &APIRouter{Router: httprouter.New(), cfg: cfg}
	api.HandlerFunc(http.MethodGet, "/", api.ListApps)
	api.HandlerFunc(http.MethodGet, "/apps/:name", api.GetApp)
	api.HandlerFunc(http.MethodGet, "/apps/:name/envs/:env", api.GetEnvironment)
	return api
}

type listAppsResponse struct {
	Apps []appResponse `json:"apps"`
}

type appResponse struct {
	Name string `json:"name"`
}
