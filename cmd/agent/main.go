package main

import (
	"context"
	"github.com/LeslyDev/metrics_and_alerts_service/internal"
	"os"
	"os/signal"
	_ "runtime"
	"sync"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup
	mclient := internal.NewMetricClient()
	//mclient.Kek()
	wg.Add(1)
	go mclient.UpdateMetrics(ctx, &wg)
	wg.Add(1)
	go mclient.SendMetrics(ctx, &wg)
	wg.Wait()
}
