package services

import (
	"kafka/domain"
	"kafka/helpers"
	"kafka/kafkas"
	"kafka/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(user domain.User) (*mongo.InsertOneResult, *utils.Resterr) {
	if err := user.Vaildate(); err != nil {
		return nil, err
	}
	guid := xid.New()
	user.ID = guid.String()
	pass := utils.HashPasswordMD5(user.Password)
	user.Password = ""
	kakfaproderr := kafkas.Produce(user)
	if kakfaproderr != nil {
		return nil, kakfaproderr
	}
	user.Password = pass
	insertNo, restErr := user.Create()
	if restErr != nil {
		return nil, restErr
	}
	return insertNo, nil
}

func LoginUser(user domain.User) (string, *utils.Resterr) {

	pass := user.Password
	result := &domain.User{Email: user.Email}
	user.Password = ""
	kakfaproderr := kafkas.Produce(user)
	if kakfaproderr != nil {
		return "", kakfaproderr
	}

	if err := result.FindUser(); err != nil {
		return "", err
	}

	hpass := result.Password

	passcheck := utils.CheckHash(hpass, pass)

	if !passcheck {
		return "", utils.BadRequest("Invalid Password")
	}

	userSignedStruct := &domain.UserJWTsigneDetails{
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			Subject:   user.Email,
		},
	}

	tokenResp, tokenErr := helpers.CreateToken(userSignedStruct)

	if tokenErr != nil {
		return "", utils.BadRequest("Token Not Create")
	}

	return tokenResp, nil

}

func GetByID(email string) (*domain.User, *utils.Resterr) {

	kakfaproderr := kafkas.Produce(domain.User{Email: email})
	if kakfaproderr != nil {
		return nil, kakfaproderr
	}

	result := &domain.User{Email: email}
	if err := result.FindUser(); err != nil {
		return nil, err
	}
	result.Password = ""
	return result, nil
}

func DeleteUser(email string) (*mongo.DeleteResult, *utils.Resterr) {

	kakfaproderr := kafkas.Produce(domain.User{Email: email})
	if kakfaproderr != nil {
		return nil, kakfaproderr
	}
	result := &domain.User{Email: email}

	deleteID, err := result.Delete()

	if err != nil {
		return nil, err
	}

	return deleteID, nil
}

func UpdateUser(isParital bool, user domain.User) (*mongo.UpdateResult, *utils.Resterr) {

	kakfaproderr := kafkas.Produce(user)
	if kakfaproderr != nil {
		return nil, kakfaproderr
	}
	result := &domain.User{Email: user.Email}

	if isParital {
		if user.Name != "" {
			result.Name = user.Name
		}

		if user.Email != "" {
			result.Email = user.Email
		}
	} else {
		result.Name = user.Name
		result.Email = user.Email
	}

	updateID, err := result.Update()
	if err != nil {
		return nil, err
	}

	return updateID, nil

}
