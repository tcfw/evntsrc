package stsmetrics

import (
	"log"
	"time"

	"github.com/robfig/cron"
)

//Cron triggers the job command every cron iteration
func Cron(natsEndpoint string, cronExpr string) {
	connectNats(natsEndpoint)
	defer natsConn.Close()

	c := cron.New()

	err := c.AddFunc(cronExpr, func() {
		findJobs(0)

		log.Printf("Next run in: %s\n", time.Until(c.Entries()[0].Next))
	})
	if err != nil {
		panic(err)
	}

	log.Println("Starting cron...")

	c.Start()

	select {}
}
