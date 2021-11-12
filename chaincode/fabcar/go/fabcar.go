/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	// "io/ioutil"
	// "os"
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	
)

// SmartContract provides functions for managing a Model
type SmartContract struct {
	contractapi.Contract
}
// // describes organisation
// type Organisation struct {
// 	Name   string `json:"name"`
// 	ID string `json:"id"`
// }

// Model describes basic details of what organisations up a Model
type Model struct {
	Organisation   string `json:"organisation"`
	Status string `json:"status"`
	Owner  string `json:"owner"`
	Organisations []string
 }


// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Model
}

// InitLedger adds a base set of organization to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	models := []Model{
		Model{Organisation: "UHCl", Status: "approved", Owner: "Sateesh",Organisations: []string{"UHCL","NASA","Approver"} },
		Model{Organisation: "NASA", Status: "approved", Owner: "Sateesh",Organisations: []string{"UHCL","NASA","Approver"}},
		Model{Organisation: "UHCl",  Status: "approved", Owner: "Sateesh",Organisations: []string{"UHCL","NASA","Approver"}},
		Model{Organisation: "UHCl", Status: "approved", Owner: "Sateesh",Organisations: []string{"UHCL","NASA","Approver"}},
		Model{Organisation: "UHCl",  Status: "approved", Owner: "Saketh",Organisations: []string{"UHCL","NASA","Approver"}},
		Model{Organisation: "UHCl",  Status: "approved", Owner: "Saketh",Organisations: []string{"UHCL","NASA","Approver"}},
		Model{Organisation: "UHCl",  Status: "pending", Owner: "Saketh",Organisations: []string{"UHCL","NASA","Approver"}},
		Model{Organisation: "UHCl",  Status: "pending", Owner: "Saketh",Organisations: []string{"UHCL","NASA","Approver"}},
		Model{Organisation: "UHCl",  Status: "pending", Owner: "Saketh",Organisations: []string{"UHCL","NASA","Approver"}},
		Model{Organisation: "UHCl", Status: "pending", Owner: "Saketh",Organisations: []string{"UHCL","NASA","Approver"}},
	}

	for i, model := range models {
		modelAsBytes, _ := json.Marshal(model) // marshal is the process of transforming the memory representation of an object to a data format suitable for storage or transmission
		err := ctx.GetStub().PutState("MODEL"+strconv.Itoa(i), modelAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new organization to the world state with given details
func (s *SmartContract) CreateModel(ctx contractapi.TransactionContextInterface, modelID string, organisation string,  status string, owner string, organisations []string) error {
	model := Model{
		Organisation:   organisation,
		Status: status,
		Owner:  owner,
		Organisations: organisations,
	}

	modelAsBytes, _ := json.Marshal(model)
	// err := ctx.GetClientIdentity().AssertAttributeValue("write", "true")
	// if err != nil {
	// 	//if unauthorised, we will throw an error
	// 	return fmt.Errorf("submitting client not authorized to create asset, does not have write permission")
	// }
	return ctx.GetStub().PutState(modelID, modelAsBytes) 
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, modelID string) (*Model, error) {
	modelAsBytes, err := ctx.GetStub().GetState(modelID) 
	
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if modelAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", modelID)
	}

	model := new(Model) 
	_ = json.Unmarshal(modelAsBytes, model) 
	if model.Status == "approved"{ 
		if evaluateViewRule(model){  
			return model, nil
		}
	}else{

		return nil, fmt.Errorf("model %s  is in pending status", modelID)
	}
	
	return nil, fmt.Errorf("you dont have view permissons for %s ", modelID)
	
}

func evaluateViewRule(model *Model) bool {
	if arryContainsString(model.Organisations,model.Organisation){
		return true
	}
	return false
}

func arryContainsString(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

// // QueryAllCars returns all cars found in world state
// func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
// 	startKey := ""
// 	endKey := ""

// 	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	results := []QueryResult{}

// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()

// 		if err != nil {
// 			return nil, err
// 		}

// 		car := new(Model)
// 		_ = json.Unmarshal(queryResponse.Value, car)

// 		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
// 		results = append(results, queryResult)
// 	}

// 	return results, nil
// }

// ChangeCarOwner updates the owner field of car with given id in world state
// func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, carNumber string, newOwner string) error {
// 	car, err := s.QueryCar(ctx, carNumber)

// 	if err != nil {
// 		return err
// 	}

// 	car.Owner = newOwner

// 	carAsBytes, _ := json.Marshal(car)

// 	return ctx.GetStub().PutState(carNumber, carAsBytes)
// }

func (s *SmartContract) UpdateModel(ctx contractapi.TransactionContextInterface, modelID string, organisation string,  status string, owner string, organisations []string) error {
	model := Model{
		Organisation:   organisation, 
		Status: status,
		Owner:  owner,
		Organisations: organisations, 
	}

	modelAsBytes, _ := json.Marshal(model) // getting model from world state
	organisation, ok, err := ctx.GetClientIdentity().GetAttributeValue("organisation") 
	fmt.Print(ok,err)
	modelAsBytesFromWorldState, err := ctx.GetStub().GetState(modelID) 

	if err != nil {
		return fmt.Errorf("%s does not exist", modelID)
	}

	modelForAssertion := new(Model)
	_ = json.Unmarshal(modelAsBytesFromWorldState, modelForAssertion)


if modelForAssertion.Status == "approved"{
	if evaluateUpdateRule(modelForAssertion,organisation){
		return ctx.GetStub().PutState(modelID, modelAsBytes)
	}
}else{

	return  fmt.Errorf("model %s  is in pending status", modelID)
}

return  fmt.Errorf("you dont have update permissons for %s ", modelID)
	
}

func evaluateUpdateRule(model *Model,organisation string) bool {
	//fmt.Print(organisation)
	//fmt.Print(model)
	// jsonFile := `{"Titronics":["NASA"],"UHCL":["NASA","Titronics"]}`
	// jsonFile, err := os.Open("/home/ubuntu/hyperledger/fabric-samples/chaincode/fabcar/go/Organisations.jsondocker exec -it ubuntu_bash bash")
	// // if we os.Open returns an error then handle it
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// byteValue, _ := ioutil.ReadAll(jsonFile)
	// var data map[string]interface{}
	// json.Unmarshal([]byte(byteValue), &data)
	// if Organisations, ok := data[model.Organisation].([]string); ok {
		// Organisations := map[string]interface{} {
		// 	"Titronics":string[]{"NASA"},

		// 	"UHCL":string[]{"NASA","Titronics"}
		// }
		// fmt.Print(Organisations)
				Organisations := []string{"NASA", "Titronics"} // hardcoded lead organizations
		if arryContainsString(Organisations,organisation){ // if user organizations belongs to lead organization.
			return true
		}
	// } else {
	// 	/* not string */
	// }
	
	return false
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
