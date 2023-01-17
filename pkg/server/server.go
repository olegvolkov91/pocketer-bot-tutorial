package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/olegvolkov91/pocketer-bot/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func NewAuthorizationServer(
	pocketClient *pocket.Client,
	tokenRepository repository.TokenRepository,
	redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{
		pocketClient:    pocketClient,
		tokenRepository: tokenRepository,
		redirectURL:     redirectURL,
	}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	chatIDQuery := r.URL.Query().Get("chat_id")
	if chatIDQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDQuery, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.createAccessToken(r.Context(), chatID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}

func (s *AuthorizationServer) createAccessToken(ctx context.Context, chatID int64) error {
	requestToken, err := s.tokenRepository.Get(chatID, repository.RequestTokens)
	if err != nil {
		return err
	}

	authResp, err := s.pocketClient.Authorize(ctx, requestToken)
	if err != nil {
		return err
	}

	if err := s.tokenRepository.Save(chatID, authResp.AccessToken, repository.AccessTokens); err != nil {
		return err
	}

	return nil
}
