package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	cookieJWT = "jwt"
)

func (s *Server) handleAuthDiscord(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	redir := q.Get("redirect")
	state := uuid.Must(uuid.NewV4()).String()

	s.State[state] = redir

	url := s.Discord.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (s *Server) handleAuthDiscordCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	state := r.FormValue("state")
	if state == "" {
		s.Logger.Debug("blank state returned")
		writeError(w, http.StatusBadRequest)
		return
	}

	redir, ok := s.State[state]
	if !ok {
		s.Logger.Debug("state not in cache", zap.String("got", state))
		writeError(w, http.StatusBadRequest)
		return
	}

	tok, err := s.Discord.Exchange(ctx, r.FormValue("code"))
	if err != nil {
		s.Logger.Error("error exchanging code", zap.Error(err))
		writeError(w, http.StatusInternalServerError)
		return
	}

	user, err := s.Discord.User(ctx, tok.AccessToken)
	if err != nil {
		s.Logger.Error("error getting user", zap.Error(err))
		writeError(w, http.StatusInternalServerError)
		return
	}

	guilds, err := s.Discord.UserGuilds(ctx, tok.AccessToken)
	if err != nil {
		s.Logger.Error("error getting user guilds", zap.Error(err))
		writeError(w, http.StatusInternalServerError)
		return
	}

	inGuild := false
	for _, g := range guilds {
		if strings.EqualFold(g.ID, s.GuildID) {
			inGuild = true
		}
	}

	if !inGuild {
		writeError(w, http.StatusForbidden)
		return
	}

	t := jwt.New()
	t.Set("uid", user.ID)
	t.Set(jwt.IssuerKey, "https://scrote.gay")
	t.Set(jwt.AudienceKey, "https://scrote.gay")

	payload, err := jwt.Sign(t, jwa.HS256, s.Key)
	if err != nil {
		s.Logger.Error("unable to sign token", zap.Error(err))
		writeError(w, http.StatusInternalServerError)
		return
	}

	expireDur := (time.Hour * 24) * 30
	expire := time.Now().Add(expireDur)

	cook := &http.Cookie{
		Name:    cookieJWT,
		Value:   string(payload),
		Domain:  s.Domain,
		Path:    "/",
		Expires: expire,
		// SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}

	if s.Debug {
		cook.Secure = false
	} else {
		cook.Secure = true
	}

	http.SetCookie(w, cook)

	// if redir != "" {
	// 	if !strings.HasPrefix(redir, "https://") {
	// 		redir = "https://" + redir
	// 	}

	// 	http.Redirect(w, r, redir, http.StatusSeeOther)
	// 	return
	// }

	// http.Redirect(w, r, "/", http.StatusSeeOther)
	w.Write([]byte(`you have been authenticated!!!!!!!!!!!!!!!!!`))
}
