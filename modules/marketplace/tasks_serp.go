package marketplace

import (
	"github.com/jasonlvhit/gocron"
)

func StartSERPCron() {
	gocron.Every(30).Minutes().Do(RefreshSerpItemsMaterializedView)
}
