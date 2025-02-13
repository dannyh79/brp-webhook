package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	r "github.com/dannyh79/brp-webhook/internal/repositories"
	rest "github.com/dannyh79/brp-webhook/internal/rest"
	"github.com/dannyh79/brp-webhook/internal/sentry"
	s "github.com/dannyh79/brp-webhook/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	Env                    string `toml:"ENV" env:"ENV" env-default:"production"`
	SentryDsn              string `toml:"SENTRY_DSN" env:"SENTRY_DSN"`
	LineChannelSecret      string `toml:"LINE_CHANNEL_SECRET" env:"LINE_CHANNEL_SECRET"`
	LineChannelAccessToken string `toml:"LINE_CHANNEL_ACCESS_TOKEN" env:"LINE_CHANNEL_ACCESS_TOKEN"`
	D1GroupQueryEndpoint   string `toml:"D1_GROUP_QUERY_ENDPOINT" env:"D1_GROUP_QUERY_ENDPOINT"`
	D1EndpointApiToken     string `toml:"D1_ENDPOINT_API_TOKEN" env:"D1_ENDPOINT_API_TOKEN"`
	Port                   int16  `toml:"PORT" env:"PORT" env-default:"8080"`
}

func main() {
	var cfg config
	err := cleanenv.ReadConfig("config.toml", &cfg)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	errs := validateConfig(cfg, cfg.Env == "development")
	if len(errs) > 0 {
		log.Println(strings.Join(errs, "\n"))
		os.Exit(2)
	}

	monitor := sentry.NewMonitor(cfg.Env, cfg.SentryDsn)
	monitor.Init()
	defer monitor.Flush()

	httpClient := &http.Client{}
	repo := r.NewD1GroupRepository(cfg.D1GroupQueryEndpoint, cfg.D1EndpointApiToken, httpClient)
	sCtx := s.NewServiceContext(
		s.NewUnlistService(repo),
		s.NewRegistrationService(repo),
		s.NewReplyService(cfg.LineChannelAccessToken, httpClient),
		s.NewWelcomeService(cfg.LineChannelAccessToken, httpClient),
	)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(sentry.SentryMiddleware())
	rest.AddRoutes(router, cfg.LineChannelSecret, sCtx)
	err = router.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		log.Fatalf("Error in starting the app: %v", err)
	}
}

type errorMsg = string

func validateConfig(cfg config, printEnv bool) []errorMsg {
	v := reflect.ValueOf(cfg)
	t := reflect.TypeOf(cfg)

	var errs []string
	for i := 0; i < v.NumField(); i++ {
		k := t.Field(i).Name
		v := v.Field(i)

		if v.Kind() == reflect.String && v.String() == "" {
			errs = append(errs, fmt.Sprintf("Missing value for %s", k))
		} else if v.IsZero() {
			errs = append(errs, fmt.Sprintf("Missing value for %s", k))
		} else if printEnv {
			log.Printf("%v: %v", k, v)
		}
	}

	return errs
}
