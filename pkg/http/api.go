package http

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/bigkevmcd/peanut/pkg/config"
)

// TODO: Add logr.Logger

// APIRouter is an HTTP API for accessing app configurations.
type APIRouter struct {
	*http.ServeMux
	cfg *config.Config
}

// ListApps returns the list of configured apps.
func (a *APIRouter) ListApps(w http.ResponseWriter, r *http.Request) {
	result := listAppsResponse{Apps: []appResponse{}}
	for _, v := range a.cfg.Apps {
		result.Apps = append(result.Apps, appResponse{Name: v.Name})
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("failed to encode resource as JSON: %s", err)
	}
}

// GetApp returns a specific app.
func (a *APIRouter) GetApp(w http.ResponseWriter, r *http.Request) {
	app := a.cfg.App(r.PathValue("name"))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(app); err != nil {
		log.Printf("failed to encode resource as JSON: %s", err)
	}
}

// GetAppConfig returns a specific app's desired state.
func (a *APIRouter) GetAppConfig(w http.ResponseWriter, r *http.Request) {
	app := a.cfg.App(r.PathValue("name"))
	w.Header().Set("Content-Type", "application/json")

	desired, err := config.ParseManifests(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := createConfigResponse(app, desired[app.Name])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode resource as JSON: %s", err)
	}
}

// GetEnvironment returns a specific app.
func (a *APIRouter) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	env := a.cfg.App(r.PathValue("name")).Environment(r.PathValue("env"))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(envResponse{Environment: env}); err != nil {
		log.Printf("failed to encode resource as JSON: %s", err)
	}
}

// NewRouter creates and returns a new APIRouter.
func NewRouter(cfg *config.Config) *APIRouter {
	api := &APIRouter{ServeMux: http.NewServeMux(), cfg: cfg}
	api.HandleFunc("GET /", api.ListApps)
	api.HandleFunc("GET /apps/{name}", api.GetApp)
	api.HandleFunc("GET /apps/{name}/desired", api.GetAppConfig)
	api.HandleFunc("GET /apps/{name}/envs/{env}", api.GetEnvironment)
	return api
}

type listAppsResponse struct {
	Apps []appResponse `json:"apps"`
}

type appResponse struct {
	Name string `json:"name"`
}

type envResponse struct {
	Environment *config.Environment `json:"environment"`
}

type configSvcResponse struct {
	Name   string   `json:"name"`
	Images []string `json:"images"`
}

type configEnvResponse struct {
	Name     string               `json:"name"`
	RelPath  string               `json:"rel_path"`
	Services []*configSvcResponse `json:"services"`
}

type configResponse struct {
	Name         string               `json:"name"`
	RepoURL      string               `json:"repo_url"`
	Path         string               `json:"path"`
	Environments []*configEnvResponse `json:"environments"`
}

func createConfigResponse(app *config.App, state map[string]map[string][]string) (*configResponse, error) {
	r := &configResponse{
		Name:         app.Name,
		RepoURL:      app.RepoURL,
		Path:         app.Path,
		Environments: []*configEnvResponse{},
	}
	err := app.EachEnvironment(func(env *config.Environment) error {
		respEnv := &configEnvResponse{Name: env.Name, RelPath: env.RelPath, Services: []*configSvcResponse{}}
		for svc, imgs := range state[env.Name] {
			sort.Strings(imgs)
			respSvc := &configSvcResponse{Name: svc, Images: []string{}}
			respSvc.Images = append(respSvc.Images, imgs...)
			respEnv.Services = append(respEnv.Services, respSvc)
		}
		r.Environments = append(r.Environments, respEnv)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}
