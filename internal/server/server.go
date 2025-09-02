package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"booking/internal/domain"
	"booking/pkg/config"
	"booking/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type FiberApp struct {
	config *config.GatewayConfig
	Log    logger.Logger
	App    *fiber.App
}

func NewFiber(
	config *config.GatewayConfig,
	log logger.Logger,
) *FiberApp {
	fiberCfg := fiber.Config{
		AppName: "gateway",
		StructValidator: &structValidator{
			validate: validator.New(),
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// cek apakah error validasi
			if ve, ok := err.(*ValidationError); ok {
				return c.Status(fiber.StatusBadRequest).JSON(domain.HttpResponse{
					Success: false,
					Message: domain.ErrInvalidRequest.Error(),
					Data:    ve.Errors,
				})
			}

			// fallback ke default fiber error
			if fe, ok := err.(*fiber.Error); ok {
				return c.Status(fe.Code).JSON(domain.HttpResponse{
					Success: false,
					Message: fe.Message,
					Data:    nil,
				})
			}

			// internal server error default
			log.Error(err, "internal server error")
			return c.Status(fiber.StatusInternalServerError).JSON(domain.HttpResponse{
				Success: false,
				Message: "internal server error",
				Data:    nil,
			})
		},
	}
	app := fiber.New(fiberCfg)

	// global middleware
	app.Use(recover.New())
	app.Use(fiberLogger.New())

	return &FiberApp{
		config: config,
		Log:    log,
		App:    app,
	}
}

func (f *FiberApp) Run() {
	port := fmt.Sprintf(":%s", f.config.Port)
	go func() {
		if err := f.App.Listen(port); err != nil {
			f.Log.Fatal(err, "error starting server")
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	f.Log.Info("Gracefully shutting down gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // timeout 10 detik
	defer cancel()

	if err := f.App.ShutdownWithContext(ctx); err != nil {
		f.Log.Error(err, "error shutting down Fiber app")
	}

	f.Log.Info("Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	f.Log.Info("Fiber was successful shutdown.")
}

// validator
type ValidationError struct { // buat type validation error, agar bias di compare di error handler fiber (harus ada method error() return string)
	Errors []map[string]string `json:"errors"`
}

func (e *ValidationError) Error() string { // supaya suatu struct bisa dianggap sebagai error, dia harus implement interface built-in error yang namanya Error() dan return string
	return "validation failed"
}

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(i any) error { // ini akan di panggil dari fiber otomatis (c.bind)
	if err := v.validate.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok { // cek apakah errornya dari golang validatorv10
			t := reflect.TypeOf(i)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}

			var errs []map[string]string
			for _, e := range validationErrors {
				field, _ := t.FieldByName(e.StructField())
				customMsg := field.Tag.Get("message")
				if customMsg == "" {
					customMsg = e.Error()
				}
				errs = append(errs, map[string]string{
					"field": e.Field(),
					"error": customMsg,
				})
			}
			// supaya *ValidationError bisa dilempar (return err) dari handler seolah-olah itu error biasa
			// mangkanya dia harus punya method error()
			// nanti dia akan di lempar otomatis ama fiber ke error handler yang udah di buat
			return &ValidationError{Errors: errs}
		}
		return err // jika errornya bukan dari golang validator return error asli
	}
	return nil // kalo tidak ada error berarti tidak otomatis di lempar ke error handler fiber
}
