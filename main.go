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
	users    = []string{"bob", "peter", "john", "alex", "tom"}
	paths    = []string{"/", "/login", "/api/v1"}
	status   = []int{200, 404, 500}
	errs     = []string{"out of memory", "cpu throttled", "circuit break"}
	services = []string{"frontend", "signup", "accounting", "api"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomness() int {
	return 1 + rand.Intn(500)
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

func serviceCall(l *zap.SugaredLogger, status int, service, path string) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	msg := "request received"
	var dur time.Duration

	// Have some outliers.
	rnd := randomness()
	switch {
	case rnd%6 == 0:
		dur = time.Millisecond * time.Duration(1000*rand.Float64())
	case rnd%12 == 0:
		dur = time.Millisecond * time.Duration(5000*rand.Float64())
	default:
		dur = time.Millisecond * time.Duration(500*rand.Float64()+0.01)
	}

	l.Infow(
		msg,
		"service", service,
		"action", "REQUEST",
		"status", status,
		"duration", dur,
		"handler", path,
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
			i := randomness() % len(users)
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
				i := randomness() % len(users)
				logins[i]++
				loginSuccess(sugar, users[i], logins[i])
			}
			time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
		}
	}()

	// Log a random error
	go func() {
		for {
			serviceFailed(sugar,
				services[randomness()%len(services)],
				errs[randomness()%len(errs)],
			)
			time.Sleep(time.Second * time.Duration(rand.Intn(10)+1))
		}
	}()

	// Some service call
	go func() {
		for {
			serviceCall(
				sugar,
				status[randomness()%len(status)],
				services[randomness()%len(services)],
				paths[randomness()%len(paths)],
			)
			time.Sleep(time.Second * time.Duration(3*rand.Float64()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	sugar.Info("Shutting down")
}
