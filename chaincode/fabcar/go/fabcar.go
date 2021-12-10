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
		Model{Organisation: "UHCL", Status: "approved", Owner: "Sateesh",Project_ID:"Project1"},
		Model{Organisation: "NASA", Status: "approved", Owner: "Sateesh",Project_ID:"Project1"},
		Model{Organisation: "UHCL",  Status: "approved", Owner: "Sateesh",Project_ID:"Project1"},
		Model{Organisation: "UHCL", Status: "approved", Owner: "Sateesh",Project_ID:"Project1"},
		Model{Organisation: "UHCL",  Status: "approved", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCL",  Status: "approved", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCL",  Status: "pending", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCL",  Status: "pending", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCL",  Status: "pending", Owner: "Saketh",Project_ID:"Project1"},
		Model{Organisation: "UHCL", Status: "pending", Owner: "Saketh",Project_ID:"Project1"},
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
// Create Project Smart Contract will put the set of organizations in couchdb
func (s *SmartContract) CreateProject(ctx contractapi.TransactionContextInterface, Project_ID string, 
	  organisations []string) error {
	project := Project{
		Project_ID: Project_ID,
		Organisations: organisations,
	}

	projectAsBytes, _ := json.Marshal(project)
	
	return ctx.GetStub().PutState(Project_ID, projectAsBytes) 
}

func (s *SmartContract) CreateModel(ctx contractapi.TransactionContextInterface, modelID string,  status string,
	 owner string, Project_ID string) error {
	organisation, ok, err := ctx.GetClientIdentity().GetAttributeValue("organisation") 
	fmt.Print(ok,err)
	model := Model{
		Organisation: organisation,														
		Status: status,
		Owner:  owner,
		Project_ID:Project_ID,
	}
	
	if evaluateCreaterule(ctx){
		if checkExistenceOfOrganisationInProject(ctx, organisation, Project_ID){
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

func checkExistenceOfOrganisationInProject(ctx contractapi.TransactionContextInterface,  organisation string, Project_ID string) bool {
	orgs,err := GetOrganisations(ctx, Project_ID)
	if err != nil {
		fmt.Printf("project doesnt exists %s", Project_ID)
		return false
   }
	if arryContainsString(orgs,organisation){
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




func (s *SmartContract) QueryModel(ctx contractapi.TransactionContextInterface, modelID string) (*Model, error) {
    organisation, ok, err := ctx.GetClientIdentity().GetAttributeValue("organisation") 
	fmt.Print(ok,err)
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
        if evaluateViewRule(ctx,organisation, model.Project_ID){ 
			return model, nil
 
        }
		
    }else{

        return nil, fmt.Errorf("model %s  is in pending status", modelID)
    }
    
    return nil, fmt.Errorf("you dont have view permissons for %s ", modelID)
    
}

func evaluateViewRule(ctx contractapi.TransactionContextInterface,organisation string, Project_ID string) bool {
	role, ok, err := ctx.GetClientIdentity().GetAttributeValue("role") 
	fmt.Print(ok,err)
	if role == "Stakeholder"{
		return checkExistenceOfOrganisationInProject(ctx,organisation, Project_ID)
	}
	return false
    

}



func (s *SmartContract) UpdateModel(ctx contractapi.TransactionContextInterface, modelID string, status string, owner string, Project_ID string ) error {
    organisation, ok, err := ctx.GetClientIdentity().GetAttributeValue("organisation") 
    fmt.Print(ok,err)
	model := Model{
        Organisation:   organisation, 
        Status: status,
        Owner:  owner,
        Project_ID:Project_ID,
    }

    modelAsBytes, _ := json.Marshal(model)
    
    modelAsBytesFromWorldState, err := ctx.GetStub().GetState(modelID) 

    if err != nil {
        return fmt.Errorf("%s does not exist", modelID)
    }

    modelForAssertion := new(Model)
    _ = json.Unmarshal(modelAsBytesFromWorldState, modelForAssertion)

if modelForAssertion.Status == "approved"{
    if evaluateUpdateRule(ctx,modelForAssertion,organisation){
        return ctx.GetStub().PutState(modelID, modelAsBytes)
    }
}else{

    return  fmt.Errorf("model %s  is in pending status", modelID)
}

return  fmt.Errorf("you dont have update permissons for %s ", modelID)
    
}

func evaluateUpdateRule(ctx contractapi.TransactionContextInterface,model *Model,organisation string) bool {
    Organisations, error:= GetOrganisations(ctx,model.Project_ID)
    role, ok, err := ctx.GetClientIdentity().GetAttributeValue("role")  
    fmt.Print(ok,err)
    if error != nil {
        fmt.Printf("%s project does not exist", model.Project_ID)
    }

        if arryContainsString(Organisations,organisation) && role == "CISE" { 
            return true
        }
    
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
