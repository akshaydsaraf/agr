package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// Agri :  Define the Agri structure, with 4 properties.  Structure tags are used by encoding/json library
type Agri struct {
	Id             string `json:"id"`
    Owner          string `json:"owner"`
    OType          string `json:"otype"`
    Grain          string `json:"grain"`
    Quantity       string `json:"quantity"`
	Quality        string `json:"quality"`
	Price          string `json:"price"`
}

type agriPrivateDetails struct {
	Owner          string `json:"owner"`
	Ownerc         string `json:"ownerc"`
}

// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("fabagr_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "queryAgri":
		return s.queryAgri(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "createAgri":
		return s.createAgri(APIstub, args)
	case "queryAllAgris":
		return s.queryAllAgris(APIstub)
	case "changeAgriOwner":
		return s.changeAgriOwner(APIstub, args)
	case "getHistoryForAsset":
		return s.getHistoryForAsset(APIstub, args)
	case "queryAgrisByOwner":
		return s.queryAgrisByOwner(APIstub, args)
	case "restictedMethod":
		return s.restictedMethod(APIstub, args)
	case "test":
		return s.test(APIstub, args)
	case "createPrivateAgri":
		return s.createPrivateAgri(APIstub, args)
	case "readPrivateAgri":
		return s.readPrivateAgri(APIstub, args)
	case "updatePrivateData":
		return s.updatePrivateData(APIstub, args)
	case "readAgriPrivateDetails":
		return s.readAgriPrivateDetails(APIstub, args)
	case "createPrivateAgriImplicitForOrg1":
		return s.createPrivateAgriImplicitForOrg1(APIstub, args)
	case "createPrivateAgriImplicitForOrg2":
		return s.createPrivateAgriImplicitForOrg2(APIstub, args)
	case "queryPrivateDataHash":
		return s.queryPrivateDataHash(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}

	// return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryAgri(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	farmAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(farmAsBytes)
}

func (s *SmartContract) readPrivateAgri(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	// collectionAgris, collectionAgriPrivateDetails, _implicit_org_Org1MSP, _implicit_org_Org2MSP
	farmAsBytes, err := APIstub.GetPrivateData(args[0], args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[1] + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if farmAsBytes == nil {
		jsonResp := "{\"Error\":\"Agri private details does not exist: " + args[1] + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(farmAsBytes)
}

func (s *SmartContract) readPrivateAgriIMpleciteForOrg1(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	farmAsBytes, _ := APIstub.GetPrivateData("_implicit_org_Org1MSP", args[0])
	return shim.Success(farmAsBytes)
}

func (s *SmartContract) readAgriPrivateDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	farmAsBytes, err := APIstub.GetPrivateData("collectionAgriPrivateDetails", args[0])

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[0] + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if farmAsBytes == nil {
		jsonResp := "{\"Error\":\"Marble private details does not exist: " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(farmAsBytes)
}

func (s *SmartContract) test(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	farmAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(farmAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	agris := []Agri{
		{Id: "a1", Owner: "Akshay", OType: "1", Grain: "Wheat", Quantity: "300", Quality: "1", Price: "400"},
		{Id: "a2", Owner: "Aditya", OType: "1", Grain: "Rice", Quantity: "300", Quality: "1", Price: "400"},
		{Id: "a3", Owner: "Ameya", OType: "1", Grain: "Barley", Quantity: "300", Quality: "1", Price: "400"},
		{Id: "a4", Owner: "Ajinkya", OType: "1", Grain: "Yellow Peas", Quantity: "300", Quality: "1", Price: "400"},
		{Id: "a5", Owner: "Avinash", OType: "1", Grain: "Wheat", Quantity: "300", Quality: "1", Price: "400"},
	}

	i := 0
	for i < len(agris) {
		farmAsBytes, _ := json.Marshal(agris[i])
		APIstub.PutState("Agri"+strconv.Itoa(i), farmAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createPrivateAgri(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	type agriTransientInput struct {
	Id             string `json:"id"`
    Owner          string `json:"owner"`
    OType          string `json:"otype"`
    Grain          string `json:"grain"`
    Quantity       string `json:"quantity"`
	Quality        string `json:"quality"`
	Price          string `json:"price"`
	Ownerc         string `json:"ownerc"`
	Key            string `json:"key"`
	}
	if len(args) != 0 {
		return shim.Error("1111111----Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	logger.Infof("11111111111111111111111111")

	transMap, err := APIstub.GetTransient()
	if err != nil {
		return shim.Error("222222 -Error getting transient: " + err.Error())
	}

	agriDataAsBytes, ok := transMap["agri"]
	if !ok {
		return shim.Error("agri must be a key in the transient map")
	}
	logger.Infof("********************8   " + string(agriDataAsBytes))

	if len(agriDataAsBytes) == 0 {
		return shim.Error("333333 -marble value in the transient map must be a non-empty JSON string")
	}

	logger.Infof("2222222")

	var farmInput agriTransientInput
	err = json.Unmarshal(agriDataAsBytes, &farmInput)
	if err != nil {
		return shim.Error("44444 -Failed to decode JSON of: " + string(agriDataAsBytes) + "Error is : " + err.Error())
	}

	logger.Infof("3333")

	if len(farmInput.Key) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(farmInput.Id) == 0 {
		return shim.Error("Id field must be a non-empty string")
	}
	if len(farmInput.Owner) == 0 {
		return shim.Error("Owner field must be a non-empty string")
	}
	if len(farmInput.OType) == 0 {
		return shim.Error("Owner type field must be a non-empty string")
	}
	if len(farmInput.Grain) == 0 {
		return shim.Error("Grain field must be a non-empty string")
	}
	if len(farmInput.Quantity) == 0 {
		return shim.Error("Quantity field must be a non-empty string")
	}
	if len(farmInput.Quality) == 0 {
		return shim.Error("Quality field must be a non-empty string")
	}
	if len(farmInput.Price) == 0 {
		return shim.Error("price field must be a non-empty string")
	}
	if len(farmInput.Ownerc) == 0 {
		return shim.Error("Owner control field must be a non-empty string")
	}

	logger.Infof("444444")

	// ==== Check if agri already exists ====
	farmAsBytes, err := APIstub.GetPrivateData("collectionAgris", farmInput.Key)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	} else if farmAsBytes != nil {
		fmt.Println("This Agri already exists: " + farmInput.Key)
		return shim.Error("This Agri already exists: " + farmInput.Key)
	}

	logger.Infof("55555")

	var agri = Agri{Id: farmInput.Id, Owner: farmInput.Owner, OType: farmInput.OType, Grain: farmInput.Grain, Quantity: farmInput.Quantity, Quality: farmInput.Quality, Price: farmInput.Price}

	farmAsBytes, err = json.Marshal(agri)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutPrivateData("collectionAgris", farmInput.Key, farmAsBytes)
	if err != nil {
		logger.Infof("6666666")
		return shim.Error(err.Error())
	}

	agriPrivateDetails := &agriPrivateDetails{Owner: farmInput.Owner, Ownerc: farmInput.Ownerc}

	agriPrivateDetailsAsBytes, err := json.Marshal(agriPrivateDetails)
	if err != nil {
		logger.Infof("77777")
		return shim.Error(err.Error())
	}

	err = APIstub.PutPrivateData("collectionAgriPrivateDetails", farmInput.Key, agriPrivateDetailsAsBytes)
	if err != nil {
		logger.Infof("888888")
		return shim.Error(err.Error())
	}

	return shim.Success(farmAsBytes)
}

func (s *SmartContract) updatePrivateData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	type agriTransientInput struct {
		Owner          string `json:"owner"`
		Ownerc         string `json:"ownerc"`
		Key            string `json:"key"`
	}
	if len(args) != 0 {
		return shim.Error("1111111----Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	logger.Infof("11111111111111111111111111")

	transMap, err := APIstub.GetTransient()
	if err != nil {
		return shim.Error("222222 -Error getting transient: " + err.Error())
	}

	agriDataAsBytes, ok := transMap["agri"]
	if !ok {
		return shim.Error("agri must be a key in the transient map")
	}
	logger.Infof("********************8   " + string(agriDataAsBytes))

	if len(agriDataAsBytes) == 0 {
		return shim.Error("333333 -marble value in the transient map must be a non-empty JSON string")
	}

	logger.Infof("2222222")

	var farmInput agriTransientInput
	err = json.Unmarshal(agriDataAsBytes, &farmInput)
	if err != nil {
		return shim.Error("44444 -Failed to decode JSON of: " + string(agriDataAsBytes) + "Error is : " + err.Error())
	}

	agriPrivateDetails := &agriPrivateDetails{Owner: farmInput.Owner, Ownerc: farmInput.Ownerc}

	agriPrivateDetailsAsBytes, err := json.Marshal(agriPrivateDetails)
	if err != nil {
		logger.Infof("77777")
		return shim.Error(err.Error())
	}

	err = APIstub.PutPrivateData("collectionAgriPrivateDetails", farmInput.Key, agriPrivateDetailsAsBytes)
	if err != nil {
		logger.Infof("888888")
		return shim.Error(err.Error())
	}

	return shim.Success(agriPrivateDetailsAsBytes)

}

func (s *SmartContract) createAgri(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var agri = Agri{Id: args[1] , Owner: args[2], OType: args[3], Grain: args[4], Quantity: args[5], Quality: args[6], Price: args[7]}

	farmAsBytes, _ := json.Marshal(agri)
	APIstub.PutState(args[0], farmAsBytes)

	indexName := "owner~key"
	colorNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{agri.Owner, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	APIstub.PutState(colorNameIndexKey, value)

	return shim.Success(farmAsBytes)
}

func (S *SmartContract) queryAgrisByOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	owner := args[0]

	ownerAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey("owner~key", []string{owner})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer ownerAndIdResultIterator.Close()

	var i int
	var id string

	var agris []byte
	bArrayMemberAlreadyWritten := false

	agris = append([]byte("["))

	for i = 0; ownerAndIdResultIterator.HasNext(); i++ {
		responseRange, err := ownerAndIdResultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		id = compositeKeyParts[1]
		assetAsBytes, err := APIstub.GetState(id)

		if bArrayMemberAlreadyWritten == true {
			newBytes := append([]byte(","), assetAsBytes...)
			agris = append(agris, newBytes...)

		} else {
			// newBytes := append([]byte(","), AgrisAsBytes...)
			agris = append(agris, assetAsBytes...)
		}

		fmt.Printf("Found a asset for index : %s asset id : ", objectType, compositeKeyParts[0], compositeKeyParts[1])
		bArrayMemberAlreadyWritten = true

	}

	agris = append(agris, []byte("]")...)

	return shim.Success(agris)
}

func (s *SmartContract) queryAllAgris(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "AGR0"
	endKey := "AGR999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllAgris:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) restictedMethod(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// get an ID for the client which is guaranteed to be unique within the MSP
	//id, err := cid.GetID(APIstub) -

	// get the MSP ID of the client's identity
	//mspid, err := cid.GetMSPID(APIstub) -

	// get the value of the attribute
	//val, ok, err := cid.GetAttributeValue(APIstub, "attr1") -

	// get the X509 certificate of the client, or nil if the client's identity was not based on an X509 certificate
	//cert, err := cid.GetX509Certificate(APIstub) -

	val, ok, err := cid.GetAttributeValue(APIstub, "role")
	if err != nil {
		// There was an error trying to retrieve the attribute
		shim.Error("Error while retriving attributes")
	}
	if !ok {
		// The client identity does not possess the attribute
		shim.Error("Client identity doesnot posses the attribute")
	}
	// Do something with the value of 'val'
	if val != "approver" {
		fmt.Println("Attribute role: " + val)
		return shim.Error("Only user with role as APPROVER have access this method!")
	} else {
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
		}

		farmAsBytes, _ := APIstub.GetState(args[0])
		return shim.Success(farmAsBytes)
	}

}

func (s *SmartContract) changeAgriOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	farmAsBytes, _ := APIstub.GetState(args[0])
	agri := Agri{}

	json.Unmarshal(farmAsBytes, &agri)
	agri.Owner = args[1]

	farmAsBytes, _ = json.Marshal(agri)
	APIstub.PutState(args[0], farmAsBytes)

	return shim.Success(farmAsBytes)
}

func (t *SmartContract) getHistoryForAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	agriName := args[0]

	resultsIterator, err := stub.GetHistoryForKey(agriName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) createPrivateAgriImplicitForOrg1(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect arguments. Expecting 5 arguments")
	}

	var agri = Agri{Id: args[1] , Owner: args[2], OType: args[3], Grain: args[4], Quantity: args[5], Quality: args[6], Price: args[7]}

	farmAsBytes, _ := json.Marshal(agri)
	// APIstub.PutState(args[0], farmAsBytes)

	err := APIstub.PutPrivateData("_implicit_org_Org1MSP", args[0], farmAsBytes)
	if err != nil {
		return shim.Error("Failed to add asset: " + args[0])
	}
	return shim.Success(farmAsBytes)
}

func (s *SmartContract) createPrivateAgriImplicitForOrg2(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 8 {
		return shim.Error("Incorrect arguments. Expecting 5 arguments")
	}

	var agri = Agri{Id: args[1] , Owner: args[2], OType: args[3], Grain: args[4], Quantity: args[5], Quality: args[6], Price: args[7]}

	farmAsBytes, _ := json.Marshal(agri)
	APIstub.PutState(args[0], farmAsBytes)

	err := APIstub.PutPrivateData("_implicit_org_Org2MSP", args[0], farmAsBytes)
	if err != nil {
		return shim.Error("Failed to add asset: " + args[0])
	}
	return shim.Success(farmAsBytes)
}

func (s *SmartContract) queryPrivateDataHash(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	farmAsBytes, _ := APIstub.GetPrivateDataHash(args[0], args[1])
	return shim.Success(farmAsBytes)
}

// func (s *SmartContract) CreateCarAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	var car Car
// 	err := json.Unmarshal([]byte(args[0]), &car)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	farmAsBytes, err := json.Marshal(car)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	err = APIstub.PutState(car.ID, farmAsBytes)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	return shim.Success(nil)
// }

// func (s *SmartContract) addBulkAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
// 	logger.Infof("Function addBulkAsset called and length of arguments is:  %d", len(args))
// 	if len(args) >= 500 {
// 		logger.Errorf("Incorrect number of arguments in function CreateAsset, expecting less than 500, but got: %b", len(args))
// 		return shim.Error("Incorrect number of arguments, expecting 2")
// 	}

// 	var eventKeyValue []string

// 	for i, s := range args {

// 		key :=s[0];
// 		var car = Car{Make: s[1], Model: s[2], Colour: s[3], Owner: s[4]}

// 		eventKeyValue = strings.SplitN(s, "#", 3)
// 		if len(eventKeyValue) != 3 {
// 			logger.Errorf("Error occured, Please make sure that you have provided the array of strings and each string should be  in \"EventType#Key#Value\" format")
// 			return shim.Error("Error occured, Please make sure that you have provided the array of strings and each string should be  in \"EventType#Key#Value\" format")
// 		}

// 		assetAsBytes := []byte(eventKeyValue[2])
// 		err := APIstub.PutState(eventKeyValue[1], assetAsBytes)
// 		if err != nil {
// 			logger.Errorf("Error coocured while putting state for asset %s in APIStub, error: %s", eventKeyValue[1], err.Error())
// 			return shim.Error(err.Error())
// 		}
// 		// logger.infof("Adding value for ")
// 		fmt.Println(i, s)

// 		indexName := "Event~Id"
// 		eventAndIDIndexKey, err2 := APIstub.CreateCompositeKey(indexName, []string{eventKeyValue[0], eventKeyValue[1]})

// 		if err2 != nil {
// 			logger.Errorf("Error coocured while putting state in APIStub, error: %s", err.Error())
// 			return shim.Error(err2.Error())
// 		}

// 		value := []byte{0x00}
// 		err = APIstub.PutState(eventAndIDIndexKey, value)
// 		if err != nil {
// 			logger.Errorf("Error coocured while putting state in APIStub, error: %s", err.Error())
// 			return shim.Error(err.Error())
// 		}
// 		// logger.Infof("Created Composite key : %s", eventAndIDIndexKey)

// 	}

// 	return shim.Success(nil)
// }

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
