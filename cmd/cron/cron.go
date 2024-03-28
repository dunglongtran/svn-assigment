package cron

import (
	"SVN-interview/cmd/cron/jobs"
	"SVN-interview/internal/common"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func StartCronJob(ctx context.Context, appCtx *common.AppContext) {
	c := cron.New()
	/** Add Job Handler*/
	// Show status job running
	c.AddFunc("@every 1m", func() {
		fmt.Println("Cron job running:", time.Now())
	})

	_, err := c.AddFunc("@every 2m", func() {
		jobs.FillHistoriesPriceJobHandler(appCtx)
	})
	if err != nil {
		fmt.Println("Cron job FillHistoriesPriceJobHandler error :", time.Now())
	}

	// start cron job
	c.Start()

	go func() {
		<-ctx.Done() // wait until context is destroyed
		fmt.Println("Cron job stopping.")
		c.Stop() // stop all cron jobs
	}()
}
