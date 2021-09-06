package web

type Option func(*Server)

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

func WithGuildID(gid string) Option {
	return func(s *Server) {
		s.GuildID = gid
	}
}

func WithClientID(cid string) Option {
	return func(s *Server) {
		s.ClientID = cid
	}
}

func WithClientSecret(cs string) Option {
	return func(s *Server) {
		s.ClientSecret = cs
	}
}

func WithRedirectURI(uri string) Option {
	return func(s *Server) {
		s.RedirectURI = uri
	}
}

func WithDebug(d bool) Option {
	return func(s *Server) {
		s.Debug = d
	}
}

func WithJWTSecret(sec string) Option {
	return func(s *Server) {
		s.SecretKey = sec
	}
}

func WithDomain(d string) Option {
	return func(s *Server) {
		s.Domain = d
	}
}
