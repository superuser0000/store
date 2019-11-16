package marketplace

import (
	"math/rand"
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type Notification struct {
	Uuid string `json:"-" gorm:"primary_key"`

	NewPrivateMessages  int `json:"new_private_messages"`
	NewPurchases        int `json:"new_transactions"`
	NewPurchaseMessages int `json:"new_transaction_messages"`

	// TimeStamps
	CreatedAt *time.Time `json:"-" gorm:"index"`
	UpdatedAt *time.Time `json:"-" gorm:"index"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

func NewRandomNotification() Notification {
	return Notification{
		Uuid:                util.GenerateUuid(),
		NewPrivateMessages:  rand.Intn(1000),
		NewPurchases:        rand.Intn(1000),
		NewPurchaseMessages: rand.Intn(1000),
	}
}

func GetCurrentNotification(user User) Notification {

	store := user.Store()
	storeUuid := ""
	if store != nil {
		storeUuid = store.Uuid
	}
	numberOfTx := CountNumberOfTransactions(user.Uuid, storeUuid)

	return Notification{
		Uuid:               util.GenerateUuid(),
		NewPrivateMessages: CountPrivateMessages(user),
		NewPurchases:       numberOfTx,
	}
}
