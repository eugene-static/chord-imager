package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"chord-drawer/app/internal/domain/service/chord"
	"chord-drawer/app/internal/middleware"
	"chord-drawer/app/internal/session"
	"chord-drawer/app/internal/transport/http/dto"
	"chord-drawer/app/pkg/config"
)

const (
	indexPath  = "app/internal/transport/http/ui/html/home.html"
	staticPath = "app/internal/transport/http/ui/static"
)
const devYear = 2023

type Service interface {
	GetNames(*chord.NameInfo) (*chord.NameResponse, error)
	GetTab(*chord.TabInfo) (*chord.TabResponse, error)
	GetPNG(*chord.ImageInfo) (*chord.ImageResponse, error)
}

type ChordHandler struct {
	log     *slog.Logger
	cfg     *config.Config
	service Service
}

func NewHandler(log *slog.Logger, cfg *config.Config, service *chord.Service) *ChordHandler {
	return &ChordHandler{
		log:     log,
		cfg:     cfg,
		service: service,
	}
}

func (h *ChordHandler) Register(router *http.ServeMux, manager *session.Manager) {
	fileServer := http.FileServer(http.Dir(staticPath))
	tempServer := http.FileServer(http.Dir(h.cfg.Server.TempStorage))
	staticStorage := fmt.Sprintf("/%s/", h.cfg.Server.StaticStorage)
	tempStorage := fmt.Sprintf("/%s/", h.cfg.Server.TempStorage)
	router.Handle(staticStorage, http.StripPrefix(staticStorage, fileServer))
	router.Handle(tempStorage, http.StripPrefix(tempStorage, middleware.CacheControl(tempServer)))
	router.HandleFunc("/", middleware.Logging(middleware.SessionHandler(h.indexHandler, manager), h.log))

}
func (h *ChordHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		h.clientError(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.log.Debug("parsing html template file", slog.String("file", indexPath))
		tmpl, err := template.ParseFiles(indexPath)
		if err != nil {
			h.serverError(w, "failed to parse html template", err)
			return
		}
		t := time.Now().Year()
		year := ""
		if t != devYear {
			year = fmt.Sprintf("-%d", t)
		}
		h.log.Debug("applying html template")
		w.WriteHeader(http.StatusOK)
		err = tmpl.Execute(w, year)
		if err != nil {
			h.serverError(w, "failed to apply html template", err)
			return
		}
	case http.MethodPost:
		request := &dto.NameRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			h.clientError(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.log.Info("decoded data",
			slog.String("event", request.Event),
			slog.String("name", request.Info.Name),
			slog.String("pattern", request.Info.Pattern),
			slog.Int("fret", request.Info.Fret),
			slog.Bool("capo", request.Info.Capo))
		if request.Event == "GetName" {
			info := request.Info
			nameInfo := &chord.NameInfo{
				Pattern: info.Pattern,
				Fret:    info.Fret,
				Capo:    info.Capo,
			}
			names, err := h.service.GetNames(nameInfo)
			if err != nil {
				h.serverError(w, "failed to analyze data", err)
				return
			}
			baseChord := &dto.Name{
				Root:     names.BaseName.Root,
				Quality:  names.BaseName.Quality,
				Extended: names.BaseName.Extended,
				Altered:  names.BaseName.Altered,
				Omitted:  names.BaseName.Omitted,
			}
			h.log.Info("base name",
				slog.String("root", names.BaseName.Root),
				slog.String("quality", names.BaseName.Quality),
				slog.String("extended", names.BaseName.Extended),
				slog.String("altered", names.BaseName.Altered),
				slog.String("omitted", names.BaseName.Omitted))
			var variations []dto.Name
			for i, v := range names.Variations {
				variations = append(variations, dto.Name{
					Root:     v.Root,
					Quality:  v.Quality,
					Extended: v.Extended,
					Altered:  v.Altered,
					Omitted:  v.Omitted,
				})
				h.log.Info("variations",
					slog.Int("var", i),
					slog.String("root", v.Root),
					slog.String("quality", v.Quality),
					slog.String("extended", v.Extended),
					slog.String("altered", v.Altered),
					slog.String("omitted", v.Omitted))
			}
			nameResponse := &dto.NameResponse{
				BaseName:   baseChord,
				Variations: variations,
			}
			err = json.NewEncoder(w).Encode(&nameResponse)
			if err != nil {
				h.serverError(w, "failed to write response data", err)
				return
			}
		}
		if request.Event == "GetTab" {
			info := request.Info
			tabInfo := &chord.TabInfo{
				Name: info.Name,
				Info: chord.NameInfo{
					Pattern: info.Pattern,
					Fret:    info.Fret,
					Capo:    info.Capo,
				},
			}
			chordTab, err := h.service.GetTab(tabInfo)
			if err != nil {
				h.serverError(w, "failed to get tab", err)
				return
			}
			tabResponse := &dto.TabResponse{Tab: chordTab.Tab}
			h.log.Info("tab", slog.String("name", request.Event))
			err = json.NewEncoder(w).Encode(&tabResponse)
			if err != nil {
				h.serverError(w, "failed to write response data", err)
				return
			}
		}
		if request.Event == "GetPNG" {
			filePath, ok := r.Context().Value(h.cfg.Session.IDKey).(string)
			if !ok {
				h.serverError(w, "failed to get session", errors.New("can't get value from context"))
			}
			info := request.Info
			imgInfo := &chord.ImageInfo{
				Name: info.Name,
				Path: filePath,
				Info: chord.NameInfo{
					Pattern: info.Pattern,
					Fret:    info.Fret,
					Capo:    info.Capo,
				},
			}
			img, err := h.service.GetPNG(imgInfo)
			if err != nil {
				h.serverError(w, "failed to analyze data", err)
				return
			}
			imgResponse := &dto.ImageResponse{URL: img.URL}
			h.log.Info("image", slog.String("url", img.URL))
			err = json.NewEncoder(w).Encode(&imgResponse)
			if err != nil {
				h.serverError(w, "failed to write response data", err)
				return
			}
		}
	}
}
func (h *ChordHandler) serverError(w http.ResponseWriter, msg string, err error) {
	h.log.Error(msg, slog.Any("details", err))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (h *ChordHandler) clientError(w http.ResponseWriter, msg string, code int) {
	h.log.Error(msg, slog.Any("details", http.StatusText(code)))
	http.Error(w, http.StatusText(code), code)
}
