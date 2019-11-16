package marketplace

import (
	"errors"
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type APISession struct {
	Token                   string    `json:"token" gorm:"primary_key"`
	UserUuid                string    `json:"-"`
	ExpiryDate              time.Time `json:"end_date"`
	IsTwoFactorSession      bool      `json:"is_2fa_session"`
	IsSecondFactorCompleted bool      `json:"is_2fa_completed"`
	SecondFactorSecretText  string    `json:"-"`

	User User `json:"-"`
}

/*
	Model Methods
*/

func (as APISession) IsActive() bool {
	return as.ExpiryDate.After(time.Now())
}

func (as *APISession) Complete2FA() error {
	as.IsSecondFactorCompleted = true
	return as.Save()
}

/*
	Model Interface Implementation
*/

func (as APISession) Validate() error {
	if as.Token == "" {
		return errors.New("Empty Token")
	}
	if as.UserUuid == "" {
		return errors.New("Empty User UUID")
	}
	return nil
}

func (as APISession) Remove() error {
	return database.Delete(&as).Error
}

func (as APISession) Save() error {
	err := as.Validate()
	if err != nil {
		return err
	}

	if existing, _ := FindAPISessionByToken(as.Token); existing != nil {
		return database.Save(&as).Error
	} else {
		return database.Create(&as).Error
	}
}

/*
	Queries
*/

func FindAPISessionByToken(sessionToken string) (*APISession, error) {
	var as APISession

	err := database.
		Preload("User").
		First(&as, "token = ?", sessionToken).
		Error

	if err != nil {
		return nil, err
	}

	return &as, nil
}

func FindAPISessionByUserUuid(userUuid string) (*APISession, error) {
	var as APISession

	err := database.
		Order("expiry_date desc").
		First(&as, "user_uuid = ?", userUuid).
		Error

	if err != nil {
		return nil, err
	}

	return &as, nil
}

func CreateAPISession(user User) (*APISession, error) {

	d, _ := time.ParseDuration("10d")

	as := APISession{
		Token:              util.GenerateUuid(),
		UserUuid:           user.Uuid,
		ExpiryDate:         time.Now().Add(d),
		IsTwoFactorSession: user.TwoFactorAuthentication,
	}

	return &as, as.Save()
}

func FindOrCreateAPISession(user User) (*APISession, error) {
	session, err := FindAPISessionByUserUuid(user.Uuid)
	if err != nil {
		return session, err
	}
	if session.ExpiryDate.After(time.Now()) {
		return session, nil
	}
	return CreateAPISession(user)
}
