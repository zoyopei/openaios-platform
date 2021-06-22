package billingclient

import (
	"context"
	"encoding/json"
	"github.com/4paradigm/openaios-platform/src/internal/billingclient/apigen"
	"github.com/4paradigm/openaios-platform/src/internal/response"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func GetBillingClient(billingServerURL string) (*apigen.Client, error) {
	return apigen.NewClient(billingServerURL)
}

func InitUserBillingAccount(client *apigen.Client, userID string, internalURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	callbackUrl := internalURL + "/internal-api/releases?user=" + userID
	// TODO REMOVE
	initBalance := 1000.0

	params := apigen.PostAccountUseridParams{CallbackUrl: callbackUrl, Balance: &initBalance}
	resp, err := client.PostAccountUserid(ctx, userID, &params)
	if err != nil {
		log.Error(err.Error())
		return errors.New("cannot create user billing account")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		log.Error(resp.Body)
		return errors.New("cannot create user billing account")
	}
	return nil
}

func GetUserBalance(client *apigen.Client, userID string) (*float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetAccountUseridBalance(ctx, userID)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("cannot get user account balance")
	}
	if resp.StatusCode != http.StatusOK {
		log.Error(resp.StatusCode)
		return nil, errors.New("cannot get user account balance")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("cannot get user account balance")
	}
	var balance float64
	err = json.Unmarshal(body, &balance)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("cannot get user account balance")
	}
	return &balance, nil
}

func GetOneComputeUnit(client *apigen.Client, userID string, computeunitID string) (*apigen.ComputeunitInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetComputeunitUseridComputeunitIdComputeunitId(ctx, userID, computeunitID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get single computeunit "+response.GetRuntimeLocation())
	}
	if resp.StatusCode == http.StatusBadRequest {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.FormatInt(int64(resp.StatusCode), 10) + " computeunit not exists. " + response.GetRuntimeLocation())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get single computeunit "+response.GetRuntimeLocation())
	}
	computeunitInfo := apigen.ComputeunitInfo{}
	err = json.Unmarshal(body, &computeunitInfo)
	if err != nil {
		return nil, errors.Wrap(err, "computeunit not exists. "+response.GetRuntimeLocation())
	}
	return &computeunitInfo, nil
}

func GetComputeUnitListByUserID(client *apigen.Client, userID string) ([]apigen.ComputeunitInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetComputeunitUserid(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get computeunit list "+response.GetRuntimeLocation())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.FormatInt(int64(resp.StatusCode), 10) + " cannot get computeunit list " + response.GetRuntimeLocation())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get computeunit list "+response.GetRuntimeLocation())
	}
	computeunitList := []apigen.ComputeunitInfo{}
	err = json.Unmarshal(body, &computeunitList)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get computeunit list "+response.GetRuntimeLocation())
	}
	return computeunitList, nil
}

func GetComputeUnitListByGroupName(client *apigen.Client, groupName string) ([]apigen.ComputeunitInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.GetComputeunitGroupGroupName(ctx, groupName)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get computeunit list "+response.GetRuntimeLocation())
	}
	if resp.StatusCode == http.StatusBadRequest {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.FormatInt(int64(resp.StatusCode), 10) + " cannot get computeunit list " + response.GetRuntimeLocation())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get computeunit list "+response.GetRuntimeLocation())
	}
	computeunitList := []apigen.ComputeunitInfo{}
	err = json.Unmarshal(body, &computeunitList)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get computeunit list "+response.GetRuntimeLocation())
	}
	return computeunitList, nil
}

func CreateInstance(client *apigen.Client, instanceName string, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := apigen.PostInstanceParams{InstanceName: instanceName, UserId: userID}
	resp, err := client.PostInstance(ctx, &params)
	if err != nil {
		return "", errors.Wrap(err, "cannot create instance "+response.GetRuntimeLocation())
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.Wrap(err, "cannot create instance "+response.GetRuntimeLocation())
		}
		if resp.StatusCode == http.StatusBadRequest {
			var respBody apigen.RequestError
			err = json.Unmarshal(body, &respBody)
			if err != nil {
				return "", errors.Wrap(err, "cannot create instance "+response.GetRuntimeLocation())
			}
			if respBody.Message != nil {
				return "", errors.New(*respBody.Message + response.GetRuntimeLocation())
			}
		}
		return "", errors.New("cannot create instance " + response.GetRuntimeLocation())
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func GetComputeUnitPrice(client *apigen.Client, ComputeUnitID string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := apigen.GetComputeunitPriceParams{ComputeunitId: ComputeUnitID}
	resp, err := client.GetComputeunitPrice(ctx, &params)
	if err != nil {
		return 0, errors.Wrap(err, "Cannot get computeunit price. "+response.GetRuntimeLocation())
	}
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("Cannot get computeunit price. " + response.GetRuntimeLocation())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrap(err, "Cannot get computeunit price. "+response.GetRuntimeLocation())
	}
	var price float64
	err = json.Unmarshal(body, &price)
	if err != nil {
		return 0, errors.Wrap(err, "Cannot get computeunit price. "+response.GetRuntimeLocation())
	}
	return price, nil
}

func AddComputeunitGroupToUser(client *apigen.Client, userID string, groupName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	params := apigen.PostComputeunitUseridParams{GroupName: groupName}
	resp, err := client.PostComputeunitUserid(ctx, userID, &params)
	if err != nil {
		return errors.Wrap(err, "cannot add group to user "+response.GetRuntimeLocation())
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("cannot add group to user " + response.GetRuntimeLocation())
	}
	return nil
}
