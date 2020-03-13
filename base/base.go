package base

import (
	"encoding/json"
	"errors"
	"fmt"
	roles "github.com/SolarLabRU/fastpay-go-commons/enums/access-role"
	"github.com/SolarLabRU/fastpay-go-commons/enums/state"
	. "github.com/SolarLabRU/fastpay-go-commons/errors"
	. "github.com/SolarLabRU/fastpay-go-commons/models"
	"github.com/SolarLabRU/fastpay-go-commons/requests"
	"github.com/SolarLabRU/fastpay-go-commons/responses"
	. "github.com/SolarLabRU/fastpay-go-commons/validation"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	ChaincodeBankName = "banks"
)

func GetSenderBank(ctx contractapi.TransactionContextInterface) (*Bank, error) {
	clientIdentity := ctx.GetClientIdentity()
	stub := ctx.GetStub()
	mspId, err := clientIdentity.GetMSPID()

	if err != nil {
		return nil, CreateError(ErrorGetMspId, fmt.Sprintf("Невозможно получить MSP ID. %s", err.Error()))
	}

	address, err := GetSenderAddressFromCertificate(clientIdentity)

	if err != nil {
		return nil, err
	}

	return GetBankByRemoteContract(stub, mspId, address)
}

func GetSenderAddressFromCertificate(identity cid.ClientIdentity) (string, error) {
	address, isFound, _ := identity.GetAttributeValue("address")

	address, isFound, _ = func() (string, bool, error) { return "263093b1c21f98c5f9b6433bf9bbb97bb87b6e79", true, nil }() // TODO Убрать

	if !isFound {
		return "", CreateError(ErrorCertificateNotValid, "Отсутвует атрибут address в сертификате")
	}

	return address, nil
}

func InvokeChaincode(stub shim.ChaincodeStubInterface, chaincodeName string, nameFunc string, params interface{}) ([]byte, error) {
	var args [][]byte

	paramsAsBytes, err := json.Marshal(params)

	if err != nil {
		return nil, CreateError(ErrorDefault, fmt.Sprintf("Ошибка парсинга входных параметров. %s", err.Error()))
	}

	args = append(args, []byte(nameFunc))
	args = append(args, paramsAsBytes)

	response := stub.InvokeChaincode(chaincodeName, args, "")

	if response.GetStatus() == 500 { // TODO спарсить код ошибки
		fmt.Println("TODO спарсить код ошибки", response.GetMessage())

		return nil, parseErrorFromAnotherChaincode(response.GetMessage())
	}

	return response.GetPayload(), nil

}

func parseErrorFromAnotherChaincode(message string) error {
	var baseError BaseError

	err := json.Unmarshal([]byte(message), &baseError)

	if err != nil {
		return CreateError(ErrorDefault, fmt.Sprintf("Ошибка при вызове чейнкода: %s", message))
	}

	return CreateError(baseError.Code, fmt.Sprintf("Ошибка при вызове чейнкода: %s", baseError.Message))

}

func GetBankByRemoteContract(stub shim.ChaincodeStubInterface, mspId string, address string) (*Bank, error) {
	request := requests.GetBank{
		Address: address,
		MSPId:   mspId,
	}

	response, err := InvokeChaincode(stub, ChaincodeBankName, "getBankByMspIdAddress", request)
	if err != nil {
		return nil, err
	}

	var bankResponse responses.BankResponse
	err = json.Unmarshal(response, &bankResponse)

	if err != nil {
		return nil, CreateError(ErrorJsonUnmarshal, fmt.Sprintf("Ошибка десерилизации ответа после вызова чейнкода banks. %s", err.Error()))
	}

	return &bankResponse.Data, nil
}

func SenderBankIsAvailable(ctx contractapi.TransactionContextInterface) error {
	bank, _ := GetSenderBank(ctx)
	return SenderBankIsAvailableWithBank(ctx, bank)
}

func SenderBankIsAvailableWithBank(ctx contractapi.TransactionContextInterface, bank *Bank) error {
	if bank == nil {
		var err error = nil
		bank, err = GetSenderBank(ctx)
		if err != nil {
			return err
		}
	}

	if bank == nil || bank.State == state.Available {
		return CreateError(ErrorBankNotAvailable, "Банк отправителя не доступен")
	}

	return nil
}

func CheckAccess(ctx contractapi.TransactionContextInterface, role roles.AccessRole, checkAvailable bool) error {
	return CheckAccessWithBank(ctx, nil, role, checkAvailable)
}

func CheckAccessWithBank(ctx contractapi.TransactionContextInterface, bank *Bank, role roles.AccessRole, checkAvailable bool) error {
	if bank == nil {
		var err error = nil
		bank, err = GetSenderBank(ctx)
		if err != nil {
			return err
		}
	}

	if checkAvailable {
		err := SenderBankIsAvailableWithBank(ctx, bank)
		if err != nil {
			return err
		}
	}

	switch role {
	case roles.Undefined:
		return CreateError(ErrorForbidden, "Права доступа к методу не определены")
	case roles.Regulator:
		if !bank.IsRegulator {
			return CreateError(ErrorForbidden, "Для досупа банк должен быть регулятором")
		}
	case roles.Owner:
		if !bank.IsRegulator {
			return CreateError(ErrorForbidden, "Для досупа банк должен быть владельцем")
		}
	}

	return nil
}

func CreateError(code uint, message string) error {
	baseError := BaseError{
		Code:    code,
		Message: message,
	}

	return createError(&baseError)
}

func CreateErrorWithData(code uint, message, data string) error {
	baseError := BaseError{
		Code:    code,
		Message: message,
		Data:    data,
	}

	return createError(&baseError)
}

func CheckArgs(args string, request interface{}) error {
	err := json.Unmarshal([]byte(args), &request)

	if err != nil {
		return CreateError(ErrorValidateDefault, fmt.Sprintf("Ошибка валидации: %s", err.Error())) // TODO
	}

	err = Validate.Struct(request)
	if err != nil {
		return CreateError(ErrorValidateDefault, fmt.Sprintf("Ошибка валидации: %s", err.Error())) // TODO
	}

	requestInterface, ok := request.(interface{ SetDefaults() })

	if ok {
		requestInterface.SetDefaults()
	}

	return nil
}

func createError(baseError *BaseError) error {
	byteError, err := json.Marshal(baseError)
	if err != nil {
		return errors.New(fmt.Sprintf("{\"code\": %d, \"message\": \"Ошибка при формирование структуры ошибки. %s\", \"data\": \"\"}", ErrorDefault, err.Error()))
	}

	return errors.New(string(byteError))
}
