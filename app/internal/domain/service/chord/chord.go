package chord

import (
	"os"
	"path/filepath"
	"sync"

	"chord-drawer/app/pkg/config"
	"github.com/eugene-static/chord-analyzer/analyzer"
)

type Service struct {
	cfg *config.Config
	mu  *sync.Mutex
}

func NewService(cfg *config.Config, mu *sync.Mutex) *Service {
	return &Service{
		cfg: cfg,
		mu:  mu,
	}
}

func (s *Service) GetNames(request *NameInfo) (*NameResponse, error) {
	info := analyzer.NewChordInfo(request.Pattern, request.Fret, request.Capo)
	names, err := info.GetNames()
	if err != nil {
		return nil, err
	}
	var baseName *Name
	var variations []*Name
	baseName = &Name{
		Root:     names.Base.Root,
		Quality:  names.Base.Quality,
		Extended: names.Base.Extended,
		Altered:  names.Base.Altered,
		Omitted:  names.Base.Omitted,
	}
	for _, v := range names.Variations {
		variations = append(variations, &Name{
			Root:     v.Root,
			Quality:  v.Quality,
			Extended: v.Extended,
			Altered:  v.Altered,
			Omitted:  v.Omitted,
		})
	}
	return &NameResponse{
		BaseName:   baseName,
		Variations: variations,
	}, nil
}

func (s *Service) GetTab(request *TabInfo) (*TabResponse, error) {
	info := analyzer.NewChordInfo(
		request.Info.Pattern,
		request.Info.Fret,
		request.Info.Capo,
	)
	tab, err := info.BuildTab(request.Name)
	if err != nil {
		return nil, err
	}
	return &TabResponse{Tab: tab}, nil
}

func (s *Service) GetPNG(request *ImageInfo) (*ImageResponse, error) {
	info := analyzer.NewChordInfo(request.Info.Pattern, request.Info.Fret, request.Info.Capo)
	data, err := info.BuildPNG(request.Name)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := os.CreateTemp(request.Path, "*.png")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return nil, err
	}
	url, err := filepath.Rel("", file.Name())
	if err != nil {
		return nil, err
	}
	url = filepath.ToSlash(file.Name())
	return &ImageResponse{URL: url}, nil
}
