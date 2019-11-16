package marketplace

import (
	"github.com/jasonlvhit/gocron"
)

func StartMessageboardCron() {
	gocron.Every(5).Minutes().Do(RefreshViewThreadsMaterializedView)
}
