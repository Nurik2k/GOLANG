package main

import (
	"context"
	"time"
)

func main() {
	lg := logrus.WithField("_start_type", "http")

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Use(middleware.RequestID())
	e.Logger.SetLevel(log.OFF)

	deliveryHttp.Register(e, deliveryHttp.New(service, timeoutSeconds))

	lg.WithFields(logrus.Fields{"address": httpAddr}).Info("listening")
	go func() {
		errs <- e.Start(httpAddr)
	}()

	return grace.NewService("http", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			lg.Errorf("failed to shutdown http server: %v", err)
		}
	})
}
