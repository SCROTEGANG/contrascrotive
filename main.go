package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/scrotums/contrascrotive/internal/web"
)

func main() {
	useEnv := flag.Bool("u", false, "Load variables from an .env file")
	envFile := flag.String("p", ".env", "Path to .env file")

	flag.Parse()

	if *useEnv {
		if err := godotenv.Load(*envFile); err != nil {
			log.Fatalln("Unable to load .env file:", err)
		}
	}

	addr := os.Getenv("CS_ADDR")
	guildID := os.Getenv("CS_GUILD_ID")
	debug := os.Getenv("CS_DEBUG")
	keySecret := os.Getenv("CS_JWT_SECRET")
	domain := os.Getenv("CS_DOMAIN")

	clientID := os.Getenv("CS_CLIENT_ID")
	clientSecret := os.Getenv("CS_CLIENT_SECRET")
	redirectURI := os.Getenv("CS_REDIRECT_URI")

	db, err := strconv.ParseBool(debug)
	if err != nil {
		db = false
	}

	srv, err := web.New(
		web.WithAddr(addr),
		web.WithGuildID(guildID),
		web.WithClientID(clientID),
		web.WithClientSecret(clientSecret),
		web.WithRedirectURI(redirectURI),
		web.WithDebug(db),
		web.WithJWTSecret(keySecret),
		web.WithDomain(domain),
	)
	if err != nil {
		log.Fatalln("error creating server:", err)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		log.Fatalln("error starting server:", err)
	}
}
