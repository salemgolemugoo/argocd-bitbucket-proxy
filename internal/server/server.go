package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/salemgolemugoo/argocd-bitbucket-proxy/internal/config"
	"github.com/salemgolemugoo/argocd-bitbucket-proxy/internal/forwarder"
	"github.com/salemgolemugoo/argocd-bitbucket-proxy/internal/translator"
	"github.com/salemgolemugoo/argocd-bitbucket-proxy/internal/validator"
)

type Server struct {
	cfg *config.Config
	fwd *forwarder.Forwarder
}

func New(cfg *config.Config, argocdURL string) *Server {
	return &Server{
		cfg: cfg,
		fwd: forwarder.New(argocdURL),
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.HandleFunc("GET /readyz", s.handleReadyz)
	mux.HandleFunc("POST /webhook", s.handleWebhook)
	return mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	if s.cfg.BitbucketServerSecret == "" && s.cfg.BitbucketCloudSecret == "" {
		http.Error(w, "no secrets configured", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	eventKey := r.Header.Get("X-Event-Key")
	if eventKey == "" {
		http.Error(w, "missing X-Event-Key header", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	secret := s.resolveSecret(body)
	if secret == "" {
		http.Error(w, "no secret configured for this webhook source", http.StatusUnauthorized)
		return
	}

	signature := r.Header.Get("X-Hub-Signature")
	if err := validator.Validate(body, signature, secret); err != nil {
		slog.Warn("signature validation failed", "error", err)
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	githubEvent, payload, err := translator.Translate(body, eventKey)
	if err != nil {
		slog.Error("translation failed", "error", err, "event", eventKey)
		http.Error(w, "translation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	ghBody, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshal GitHub payload", "error", err)
		http.Error(w, "marshal failed", http.StatusInternalServerError)
		return
	}

	if err := s.fwd.Forward(githubEvent, ghBody); err != nil {
		slog.Error("failed to forward to ArgoCD", "error", err)
		http.Error(w, "forward failed: "+err.Error(), http.StatusBadGateway)
		return
	}

	slog.Info("webhook forwarded", "source_event", eventKey, "github_event", githubEvent)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) resolveSecret(body []byte) string {
	if translator.IsBitbucketServer(body) && s.cfg.BitbucketServerSecret != "" {
		return s.cfg.BitbucketServerSecret
	}
	if translator.IsBitbucketCloud(body) && s.cfg.BitbucketCloudSecret != "" {
		return s.cfg.BitbucketCloudSecret
	}
	return ""
}
