package utilities

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/balaganapathyparthiban/quickshare-go/db"
	"github.com/robfig/cron/v3"
)

func InitCron() {
	c := cron.New()

	c.AddFunc("@midnight", func() {
		fmt.Println("Cron Job Started")
		iter := db.Store.NewIterator(nil, nil)
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()

			var fdUnmarshal FileData
			json.Unmarshal(value, &fdUnmarshal)

			isExpired := time.Now().After(fdUnmarshal.Expired)

			if isExpired {
				err := os.RemoveAll(fmt.Sprintf("files/%s", key))
				if err != nil {
					fmt.Println(err)
					continue
				}

				db.Store.Delete([]byte(key), nil)
			}

		}
		iter.Release()
		err := iter.Error()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Cron Job Ended")
	})

	c.Start()
}
