package main

import (
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	users    = []string{"bob", "peter", "john", "alex", "tom"}
	paths    = []string{"/", "/login", "/api/v1"}
	errs     = []string{"out of memory", "cpu throttled", "circuit break"}
	services = []string{"frontend", "signup", "accounting", "api"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func loginError(l *zap.SugaredLogger, user string, attempt int) {
	l.Infow(
		"user failed to log in",
		"user", user,
		"action", "LOGIN",
		"result", "FAILED",
		"attempt", attempt,
	)
}

func loginSuccess(l *zap.SugaredLogger, user string, logins int) {
	l.Infow(
		"user successfully logged in",
		"user", user,
		"action", "LOGIN",
		"result", "SUCCESS",
		"login", logins,
	)
}

func serviceFailed(l *zap.SugaredLogger, service, err string) {
	l.Errorw(
		"service failed",
		"service", service,
		"err", errors.New(err),
	)
}

func initLogger(path string) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{path, "stdout"}
	cfg.EncoderConfig = encoderConfig

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func main() {
	logger := initLogger("./demo_log.log")
	defer logger.Sync()
	sugar := logger.Sugar()

	// Log failed logins
	go func() {
		attempts := make([]int, len(users))
		for {
			i := rand.Intn(500) % len(users)
			attempts[i]++
			loginError(sugar, users[i], attempts[i])
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		}
	}()

	// Log normal logins
	go func() {
		logins := make([]int, len(users))
		for {
			if time.Now().Second() >= 40 && time.Now().Second() <= 59 {
				i := rand.Intn(500) % len(users)
				logins[i]++
				loginSuccess(sugar, users[i], logins[i])
			}
			time.Sleep(time.Second * time.Duration(rand.Intn(4)+1))
		}
	}()

	// Log a random error
	go func() {
		for {
			iError := (1 + rand.Intn(500)) % len(errs)
			iService := (1 + rand.Intn(500)) % len(services)

			serviceFailed(sugar, services[iService], errs[iError])
			time.Sleep(time.Second * time.Duration(rand.Intn(10)+1))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	sugar.Info("Shutting down")
}
