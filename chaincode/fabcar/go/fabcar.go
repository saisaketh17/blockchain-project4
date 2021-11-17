package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	
)

// SmartContract provides functions for managing a Model
type SmartContract struct {
	contractapi.Contract
}

// Model describes basic details of what organisations up a Model
type Model struct {
	Organisation   string `json:"organisation"`
	Status string `json:"status"`
	Owner  string `json:"owner"`
	Project_ID string `json:"Project_ID"`
	// Organisations []string
 }

 type Project struct  {
	 Project_ID string `json:"Project_ID"`
	Organisations []string
 }



// InitLedger adds a base set of organization to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	models := []Model{
		Model{Organisation: "UHCl", Status: "approved", Owner: "Sateesh",Project_ID:"Project1" },
		Model{Organisation: "NASA", Status: "approved", Owner: "Sateesh",Project_ID:"Project1"},
		Model{Organisation: "UHCl",  Status: "approved", Owner: "Sateesh",Project_ID:"Project1"},
		Model{Organisation: "UHCl", Status: "approved", Owner: "Sateesh",Project_ID:"Project1"},
		Model{Organisation: "UHCl",  Status: "approved", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCl",  Status: "approved", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCl",  Status: "pending", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCl",  Status: "pending", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCl",  Status: "pending", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCl", Status: "pending", Owner: "Saketh",Project_ID:"Project1"},
	}

	for i, model := range models {
		modelAsBytes, _ := json.Marshal(model) 
		err := ctx.GetStub().PutState("MODEL"+strconv.Itoa(i), modelAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) CreateModel(ctx contractapi.TransactionContextInterface, modelID string,  status string, owner string, Project_ID string) error {
	organisation, ok, err := ctx.GetClientIdentity().GetAttributeValue("organisation") 
	fmt.Print(ok,err)
	model := Model{
		Organisation: organisation,														
		Status: status,
		Owner:  owner,
		Project_ID:Project_ID,
	}
	
	if evaluateCreaterule(ctx){
		if checkExistenceOfOrganisationInProject(ctx,&model){
			modelAsBytes, _ := json.Marshal(model)
			return ctx.GetStub().PutState(modelID, modelAsBytes) 
		}else{
			return  fmt.Errorf("Your organistaion doenst belong to project Organisations")
		}
		
	}
	
	return  fmt.Errorf("you dont have enough permissions to create a model")
	
}

func evaluateCreaterule(ctx contractapi.TransactionContextInterface)bool{
	role, ok, err := ctx.GetClientIdentity().GetAttributeValue("role") 
	fmt.Print(ok,err)
	if role == "CISE"{
		return true
	}
	return false
}
func checkExistenceOfOrganisationInProject(ctx contractapi.TransactionContextInterface,model *Model) bool {
	orgs,err := GetOrganisations(ctx,model.Project_ID )
	if err != nil {
		fmt.Printf("project doesnt exists %s", model.Project_ID)
		return false
   }
	if arryContainsString(orgs,model.Organisation){
		return true
	}
	return false
}
func GetOrganisations(ctx contractapi.TransactionContextInterface, Project_ID string) ([]string,error) {
	projectAsBytes, err := ctx.GetStub().GetState(Project_ID) 
	
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if projectAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", Project_ID)
	}

	project := new(Project) 
	_ = json.Unmarshal(projectAsBytes, project) 

	return project.Organisations,nil
	
}

func arryContainsString(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func (s *SmartContract) CreateProject(ctx contractapi.TransactionContextInterface, Project_ID string, organisations []string) error {
	project := Project{
		Project_ID: Project_ID,
		Organisations: organisations,
	}

	projectAsBytes, _ := json.Marshal(project)
	
	return ctx.GetStub().PutState(Project_ID, projectAsBytes) 
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
