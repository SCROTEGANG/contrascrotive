package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/holedaemon/discord"
	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	defaultAddr = ":8080"
)

var (
	ErrInvalidOption = errors.New("server: invalid option")

	defaultScopes = []string{"identify", "guilds"}
)

type Server struct {
	Addr      string
	GuildID   string
	SecretKey string
	Domain    string
	Debug     bool

	ClientID     string
	ClientSecret string
	RedirectURI  string

	Discord *discord.Client
	Logger  *zap.Logger
	State   map[string]string
	Key     jwk.Key
}

func New(opts ...Option) (*Server, error) {
	s := &Server{
		State: make(map[string]string),
	}

	for _, o := range opts {
		o(s)
	}

	if err := s.defaults(); err != nil {
		return nil, err
	}

	var logger *zap.Logger
	if s.Debug {
		dl, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

		logger = dl
	} else {
		pl, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}

		logger = pl
	}
	s.Logger = logger

	jk, err := jwk.New([]byte(s.SecretKey))
	if err != nil {
		return nil, err
	}
	s.Key = jk

	oa := &oauth2.Config{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		RedirectURL:  s.RedirectURI,
		Endpoint:     discord.Endpoint,
		Scopes:       defaultScopes,
	}

	disc, err := discord.New(
		discord.WithOAuth2Config(oa),
	)
	if err != nil {
		return nil, err
	}
	s.Discord = disc

	return s, nil
}

func (s *Server) defaults() error {
	if s.Addr == "" {
		s.Addr = defaultAddr
	}

	if s.GuildID == "" {
		return invalidOption("guild id cannot be blank")
	}

	if s.ClientID == "" {
		return invalidOption("client id cannot be blank")
	}

	if s.ClientSecret == "" {
		return invalidOption("client secret cannot be blank")
	}

	if s.RedirectURI == "" {
		return invalidOption("redirect uri cannot be blank")
	}

	if s.SecretKey == "" {
		return invalidOption("jwt secret cannot be blank")
	}

	if s.Domain == "" {
		return invalidOption("domain cannot be blank")
	}

	return nil
}

func (s *Server) Start(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Get("/", s.handleAuthDiscord)
	r.Get("/callback", s.handleAuthDiscordCallback)

	srv := &http.Server{
		Addr:    s.Addr,
		Handler: r,
	}

	go func() {
		<-ctx.Done()

		if err := srv.Shutdown(context.Background()); err != nil {
			s.Logger.Error("error shutting down server", zap.Error(err))
		}
	}()

	return srv.ListenAndServe()
}
