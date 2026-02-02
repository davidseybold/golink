package golink

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	configFile = flag.String("config", "", "path of config file")
)

func Run(logger *slog.Logger) error {
	flag.Parse()

	if *configFile == "" {
		return errors.New("config file is required")
	}

	cfg, err := loadConfig(*configFile)
	if err != nil {
		return err
	}

	linksByShort, err := buildLinksIndex(cfg.Links)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			short := strings.TrimPrefix(r.URL.Path, "/")
			if short == "" {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("Nothing to see here—try /<short>\n"))
				return
			}

			long, ok := linksByShort[short]
			if !ok {
				http.NotFound(w, r)
				logger.Info("golink not found", "short", short)
				return
			}

			target := normalizeTarget(cfg.DefaultScheme, long)
			http.Redirect(w, r, target, http.StatusFound)
			logger.Info("redirected", "short", short, "long", long, "target", target)
		}),
	}

	logger.Info("golink listening", "addr", server.Addr)

	return server.ListenAndServe()
}

type Config struct {
	Port          int    `yaml:"port"`
	DefaultScheme string `yaml:"default_scheme"`
	Links         []Link `yaml:"links"`
}

type Link struct {
	Short string `yaml:"short"`
	Long  string `yaml:"long"`
}

func loadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}

	if cfg.Port == 0 {
		return Config{}, errors.New("config port is required")
	}

	if cfg.DefaultScheme == "" {
		return Config{}, errors.New("config default_scheme is required")
	}

	return cfg, nil
}

func buildLinksIndex(links []Link) (map[string]string, error) {
	linksByShort := make(map[string]string, len(links))
	for _, l := range links {
		short := strings.TrimSpace(l.Short)
		long := strings.TrimSpace(l.Long)
		if short == "" {
			return nil, errors.New("link short is required")
		}
		if long == "" {
			return nil, fmt.Errorf("link long is required for short %q", short)
		}
		if _, exists := linksByShort[short]; exists {
			return nil, fmt.Errorf("duplicate link short %q", short)
		}
		linksByShort[short] = long
	}
	return linksByShort, nil
}

func normalizeTarget(defaultScheme, long string) string {
	if strings.HasPrefix(long, "http://") || strings.HasPrefix(long, "https://") {
		return long
	}
	return defaultScheme + "://" + long
}
