package utils

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/4paradigm/openaios-platform/src/billing/conf"
	"github.com/4paradigm/openaios-platform/src/internal/mongodb"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"time"
)

type PodInfo struct {
	UserId          string    `bson:"userId,omitempty"`
	PodName         string    `bson:"podName,omitempty"`
	PodUID          types.UID `bson:"podUID,omitempty"`
	InstanceId      string    `bson:"instanceId,omitempty"`
	ComputeunitList []string  `bson:"computeunitList,omitempty"`
	StartTime       time.Time `bson:"startTime,omitempty"`
	UpdateTime      time.Time `bson:"updateTime,omitempty"`
	Count           int64     `bson:"count,omitempty"`
}

type AccountInfo struct {
	UserId           string   `bson:"userId,omitempty"`
	Balance          *float64 `bson:"balance,omitempty"`
	CallbackUrl      string   `bson:"callbackUrl,omitempty"`
	ComputeunitGroup []string `bson:"computeunitGroup,omitempty"`
}

type InstanceInfo struct {
	// using _id as instance id
	InstanceName string `bson:"instanceName,omitempty"`
	UserId       string `bson:"userId,omitempty"`
}

var database = conf.GetMongodbDatabase()

var podColl = "pod"
var userColl = "user"
var instanceColl = "instance"
var computeunitGroupColl = "computeunitGroup"
var computeunitColl = "computeunit"

func InitColl(client *mongo.Client) error {
	if client == nil {
		return errors.New("mongodb client is nil.")
	}
	if err := mongodb.CreateUniqueIndex(client, database, podColl, "podUID"); err != nil {
		return err
	}
	if err := mongodb.CreateUniqueIndex(client, database, userColl, "userId"); err != nil {
		return err
	}
	if err := mongodb.CreateUniqueIndex(client, database, computeunitGroupColl, "groupName"); err != nil {
		return err
	}
	if err := mongodb.CreateUniqueIndex(client, database, computeunitColl, "id"); err != nil {
		return err
	}

	return nil
}

func CreateUserWithBalance(client *mongo.Client, userId string,
	balance float64, callbackUrl string) error {
	_, err := mongodb.InsertOneDocument(client, database, userColl,
		AccountInfo{UserId: userId,
			Balance:          &balance,
			CallbackUrl:      callbackUrl,
			ComputeunitGroup: []string{"default"}})
	return errors.Wrap(err, response.GetRuntimeLocation())
}

func GetUserBalance(client *mongo.Client, userId string) (float64, error) {
	uniqueKey := AccountInfo{UserId: userId}
	document := mongodb.FindOneDocument(client, database, userColl, uniqueKey)
	if document == nil {
		return 0, errors.New("Cannot find user " + response.GetRuntimeLocation())
	} else {
		var accountInfo AccountInfo
		err := document.Decode(&accountInfo)
		if err != nil {
			return 0, errors.Wrap(err, response.GetRuntimeLocation())
		}
		return *accountInfo.Balance, nil
	}
}

func DeleteUser(client *mongo.Client, userId string) error {
	uniqueKey := AccountInfo{UserId: userId}
	return mongodb.DeleteOneDocument(client, database, userColl, uniqueKey)
}

func ModifyUserBalance(client *mongo.Client, userId string, balance float64) error {
	uniqueKey := AccountInfo{UserId: userId}
	minBalance := 0.0
	_, err := mongodb.UpdateOneDocument(client, database, userColl, uniqueKey,
		mongodb.MongodbOperation{Operator: "$inc", Document: AccountInfo{Balance: &balance}})
	if err != nil {
		return errors.Wrap(err, response.GetRuntimeLocation())
	}
	_, err = mongodb.UpdateOneDocument(client, database, userColl, uniqueKey,
		mongodb.MongodbOperation{Operator: "$max", Document: AccountInfo{Balance: &minBalance}})
	if err != nil {
		return errors.Wrap(err, response.GetRuntimeLocation())
	}
	//else if modifyCount == 0 {
	//	log.Warnf("cannot find such user %s.", userId)
	//}
	return nil
}

func UpdateUserAccount(client *mongo.Client, userID string, balance *float64,
	callbackURL *string, computeunitGroup *[]string) error {
	accountInfo := AccountInfo{}
	uniqueKey := AccountInfo{UserId: userID}

	if balance != nil {
		accountInfo.Balance = balance
	}
	if callbackURL != nil {
		accountInfo.CallbackUrl = *callbackURL
	}
	if computeunitGroup != nil {
		accountInfo.ComputeunitGroup = *computeunitGroup
	}

	_, err := mongodb.UpdateOneDocument(client, database, userColl, uniqueKey,
		mongodb.MongodbOperation{Operator: "$set", Document: accountInfo})
	if err != nil {
		return errors.Wrap(err, "update user account failed "+response.GetRuntimeLocation())
	}
	//else if modifyCount == 0 {
	//	return errors.New("cannot find such user " + GetRuntimeLocation())
	//}
	return nil
}

func GetAccountList(client *mongo.Client) ([]AccountInfo, error) {
	cursor, err := mongodb.FindDocuments(client, database, userColl, "")
	if err != nil {
		return nil, errors.Wrap(err, "find documents failed. "+response.GetRuntimeLocation())
	}
	var accountList = []AccountInfo{}
	for cursor.Next(context.Background()) {
		var accountInfo = AccountInfo{}
		if err = cursor.Decode(&accountInfo); err != nil {
			log.Warn(err)
			continue
		}
		accountList = append(accountList, accountInfo)
	}
	return accountList, nil
}

func CheckUserNoBalance(client *mongo.Client, billMap map[string]float64) error {
	operator := mongodb.ComparisonQueryOperator{Operation: "$lte", Value: 0}
	cursor, err := mongodb.FindDocuments(client, database, userColl, "balance", operator)
	if err != nil {
		return err
	}
	var accountInfo AccountInfo
	for cursor.Next(context.Background()) {
		err = cursor.Decode(&accountInfo)
		if err != nil {
			log.Error(err.Error())
		} else if _, ok := billMap[accountInfo.UserId]; ok {
			log.Warn(accountInfo.UserId + " has no balance.")
			request, err := http.NewRequest(http.MethodDelete, accountInfo.CallbackUrl, nil)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			if resp.StatusCode != http.StatusOK {
				log.Warn(resp.StatusCode)
			}
		}
	}
	return nil
}

func UpdateOrInsertPod(client *mongo.Client, pod PodInfo) error {
	uniqueKey := PodInfo{PodUID: pod.PodUID}
	podInfo := PodInfo{ComputeunitList: pod.ComputeunitList, StartTime: pod.StartTime, PodName: pod.PodName,
		UpdateTime: pod.UpdateTime, PodUID: pod.PodUID, UserId: pod.UserId, InstanceId: pod.InstanceId}
	setOperation := mongodb.MongodbOperation{Operator: "$set", Document: podInfo}
	incOperation := mongodb.MongodbOperation{Operator: "$inc", Document: PodInfo{Count: 1}}
	return mongodb.UpdateOrInsertOneDocument(client, database, podColl,
		uniqueKey, setOperation, incOperation)
}

func CreateInstance(client *mongo.Client, userId string, instanceName string) (string, error) {
	balance, err := GetUserBalance(client, userId)
	if err != nil {
		log.Error(err.Error())
		return "", errors.New("Cannot find user.")
	} else if balance <= 0 {
		return "", errors.New("User does not have enough balance.")
	}

	// create instance
	instanceInfo := InstanceInfo{InstanceName: instanceName, UserId: userId}
	result, err := mongodb.InsertOneDocument(client, database, instanceColl, instanceInfo)
	if err != nil {
		return "", err
	}
	return mongodb.InsertedIdToObjectId(result.InsertedID)
}
