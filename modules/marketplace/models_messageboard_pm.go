package marketplace

import (
	"fmt"
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type PrivateThread struct {
	Uuid             string    `json:"uuid"`
	LastUpdated      time.Time `json:"last_updated"`
	LastReadDate     time.Time `json:"last_read"`
	NumberOfMessages int       `json:"number_of_messages"`
	SenderUsername   string    `json:"sender_username"`
	ReceiverUsername string    `json:"receiver_username"`
	SenderUserUuid   string    `json:"-"`
	ReceiverUserUuid string    `json:"-"`
}

type ViewPrivateThread struct {
	*PrivateThread
	LastUpdatedStr  string `json:"last_updated_str"`
	DisplayUsername string `json:"diplay_username"`
	IsRead          bool   `json:"is_read"`
}

func (pt PrivateThread) ViewPrivateThread(lang, username string) ViewPrivateThread {
	vpt := ViewPrivateThread{
		PrivateThread:  &pt,
		LastUpdatedStr: util.HumanizeTime(pt.LastUpdated, lang),
		IsRead:         pt.LastReadDate.After(pt.LastUpdated),
	}
	if username == pt.ReceiverUsername {
		vpt.DisplayUsername = pt.SenderUsername
	} else {
		vpt.DisplayUsername = pt.ReceiverUsername
	}
	return vpt
}

type PrivateThreads []PrivateThread

func (msbts PrivateThreads) ViewPrivateThreads(lang, username string) []ViewPrivateThread {
	var vpts []ViewPrivateThread
	for _, msbt := range msbts {
		vpt := msbt.ViewPrivateThread(lang, username)
		vpts = append(vpts, vpt)
	}

	return vpts
}

/*
	Queries
*/

func FindPrivateThreads(userUuid string) PrivateThreads {
	threads := PrivateThreads{}

	database.
		Table("v_private_threads").
		Joins(fmt.Sprintf("left outer join thread_perusal_statuses on v_private_threads.uuid=thread_perusal_statuses.thread_uuid AND thread_perusal_statuses.user_uuid = '%s'", userUuid)).
		Where("(sender_user_uuid=? OR reciever_user_uuid=?)", userUuid, userUuid).
		Select("v_private_threads.*, thread_perusal_statuses.last_read_date").
		Find(&threads)

	return threads
}

func FindThreadPerusalStatus(threadUuid, userUuid string) (*ThreadPerusalStatus, error) {
	var ts ThreadPerusalStatus

	err := database.
		Where(&ThreadPerusalStatus{
			UserUuid:   userUuid,
			ThreadUuid: threadUuid,
		}).
		First(&ts).Error

	if err != nil {
		return nil, err
	}
	return &ts, err
}

func FindPrivateThread(u1, u2 User) Thread {
	var (
		thread Thread
	)

	database.
		Table("v_threads").
		Where(`
			type=? AND (
				(reciever_user_uuid=? AND sender_user_uuid=?) OR
				(reciever_user_uuid=? AND sender_user_uuid=?)
			)`,
			"private",
			u1.Uuid,
			u2.Uuid,
			u2.Uuid,
			u1.Uuid,
		).
		Preload("SenderUser").
		Preload("RecieverUser").
		First(&thread)

	return thread
}

/*
	Collection Methods
*/

func (ths PrivateThreads) CountUnread() int {
	counter := 0
	for _, th := range ths {
		if th.LastReadDate.After(th.LastUpdated) {
			counter += 1
		}
	}
	return counter
}

/*
	Cache
*/

func CacheFindPrivateThreads(userUuid string) PrivateThreads {
	key := "private-threads-" + userUuid
	cThreads, _ := cache15m.Get(key)
	if cThreads == nil {
		threads := FindPrivateThreads(userUuid)
		cache15m.Set(key, threads)

		return threads
	}

	return cThreads.(PrivateThreads)
}

func CacheClearPrivateThreads(userUuid string) {
	key := "private-threads-" + userUuid
	cache15m.Remove(key)
}

/*
	DB Views
*/

func setupPrivateThreadsFunctions() {
	database.Exec("DROP VIEW IF EXISTS v_private_threads CASCADE;")
	database.Exec(`
	CREATE VIEW v_private_threads AS (
		select 
			vm_threads.uuid,
			vm_threads.is_read_by_reciever,
			vm_threads.last_updated,
			vm_threads.number_of_messages,
			u1.username as sender_username,
			u2.username as receiver_username,
			vm_threads.sender_user_uuid,
			vm_threads.reciever_user_uuid

		from
			vm_threads
		join users u1 on u1.uuid=vm_threads.sender_user_uuid
		join users u2 on u2.uuid=vm_threads.reciever_user_uuid
		ORDER BY vm_threads.last_updated DESC
	);`)
}
