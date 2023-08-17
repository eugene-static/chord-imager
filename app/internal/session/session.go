package session

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"chord-drawer/app/pkg/config"
	"chord-drawer/app/pkg/random"
)

const (
	idLength = 16
)

type Session struct {
	Registered time.Time
	Expired    *time.Timer
	FilePath   string
}

type Manager struct {
	storage map[string]*Session
	log     *slog.Logger
	cfg     *config.Config
	mu      *sync.Mutex
}

func Init(log *slog.Logger, cfg *config.Config, mu *sync.Mutex) (*Manager, error) {
	if err := os.Mkdir(cfg.Server.TempStorage, os.ModeDir); err != nil {
		if errors.Is(err, os.ErrExist) {
			mu.Lock()
			if err = os.RemoveAll(filepath.Join(cfg.Server.TempStorage, "*")); err != nil {
				return nil, err
			}
			mu.Unlock()
		} else {
			return nil, err
		}
	}
	return &Manager{
		storage: make(map[string]*Session),
		cfg:     cfg,
		log:     log,
		mu:      mu,
	}, nil
}

func (m *Manager) Create(id string) (string, error) {
	if id == "" {
		id = random.String(idLength)
	}
	filePath := filepath.Join(m.cfg.Server.TempStorage, id)
	if err := os.Mkdir(filePath, os.ModeDir); err != nil && !errors.Is(err, os.ErrExist) {
		m.log.Error("failed to create session", slog.String("session", id), slog.Any("details", err))
		return "", err
	}
	session := &Session{
		Registered: time.Now(),
		Expired:    time.AfterFunc(time.Duration(m.cfg.Session.Lifetime)*time.Minute, m.clean(id)),
		FilePath:   filePath,
	}
	m.log.Info("recording new session", slog.String("session ID", id))
	m.storage[id] = session
	return id, nil
}

func (m *Manager) Check(id string) bool {
	_, ok := m.storage[id]
	return ok
}

func (m *Manager) Get(id string) *Session {
	return m.storage[id]
}

func (m *Manager) UpdateTime(id string) {
	timer := m.storage[id].Expired
	if !timer.Stop() {
		<-timer.C
	}
	m.storage[id].Expired = time.AfterFunc(time.Duration(m.cfg.Session.Lifetime)*time.Minute, m.clean(id))
}

func (m *Manager) clean(id string) func() {
	return func() {
		m.log.Info("session timed out",
			slog.String("id", id),
			slog.String("duration", time.Since(m.storage[id].Registered).String()))
		m.mu.Lock()
		defer m.mu.Unlock()
		if err := os.RemoveAll(m.storage[id].FilePath); err != nil {
			m.log.Warn("failed to remove folder",
				slog.Any("details", err),
				slog.String("dir", m.storage[id].FilePath))
		}
		delete(m.storage, id)

	}
}

func (m *Manager) Reset() {
	m.log.Info("clear and close session manager")
	for _, s := range m.storage {
		if !s.Expired.Stop() {
			<-s.Expired.C
		}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := os.RemoveAll(m.cfg.Server.TempStorage); err != nil {
		m.log.Error("failed to remove temp storage",
			slog.Any("details", err))
	}
	m.storage = nil
}

//func (s *Manager) RunCleaner(mu *sync.Mutex) {
//	go func() {
//		for {
//			for _, session := range s.storage {
//				select {
//				case <-session.Registered.C:
//					s.log.Debug("session time has expired", slog.String("session", session.SessionID))
//					if _, err := os.Stat(session.DirPath); os.IsNotExist(err) {
//						continue
//					}
//					mu.Lock()
//					err := os.RemoveAll(session.DirPath)
//					if err != nil {
//						s.log.Error("failed to remove folder", slog.Any("details", err), slog.String("dir", session.DirPath))
//					}
//					delete(s.storage, session.SessionID)
//					mu.Unlock()
//				}
//			}
//			time.Sleep(cleanerLoop * time.Second)
//		}
//	}()
//}
