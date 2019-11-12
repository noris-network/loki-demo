package main

import (
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	users      = []string{"bob", "peter", "john", "alex", "tom"}
	paths      = []string{"/", "/login", "/api/v1"}
	status     = []int{200, 404, 500}
	errs       = []string{"out of memory", "cpu throttled", "circuit break"}
	services   = []string{"frontend", "signup", "accounting", "api"}
	randomness = 1 + rand.Intn(500)
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

func serviceCall(l *zap.SugaredLogger, service string, status int) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	msg := "request received"
	dur := time.Millisecond * time.Duration(500*rand.Float64()+0.01)

	l.Infow(
		msg,
		"service", service,
		"status", status,
		"duration", dur,
		"traceID", uuid,
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
	logger := initLogger("./logs/demo_log.log")
	defer logger.Sync()
	sugar := logger.Sugar()

	// Log failed logins
	go func() {
		attempts := make([]int, len(users))
		for {
			i := randomness % len(users)
			attempts[i]++
			loginError(sugar, users[i], attempts[i])
			time.Sleep(time.Second * time.Duration(rand.Intn(6)+1))
		}
	}()

	// Log normal logins
	go func() {
		logins := make([]int, len(users))
		for {
			if time.Now().Second() >= 40 && time.Now().Second() <= 59 {
				i := randomness % len(users)
				logins[i]++
				loginSuccess(sugar, users[i], logins[i])
			}
			time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
		}
	}()

	// Log a random error
	go func() {
		for {
			iError := randomness % len(errs)
			iService := randomness % len(services)
			serviceFailed(sugar, services[iService], errs[iError])
			time.Sleep(time.Second * time.Duration(rand.Intn(10)+1))
		}
	}()

	// Some service call
	go func() {
		for {
			iService := randomness % len(services)
			iStatus := randomness % len(status)
			serviceCall(sugar, services[iService], status[iStatus])
			time.Sleep(time.Second * time.Duration(3*rand.Float64()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	sugar.Info("Shutting down")
}
