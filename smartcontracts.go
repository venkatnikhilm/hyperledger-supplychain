package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

// ProductEntity represents the structure of a product in the supply chain
type ProductEntity struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductStatus string `json:"product_status"`
	CurrentOwner string `json:"current_owner"`
	CreatedDate  string `json:"created_date"`
	UpdatedDate  string `json:"updated_date"`
	ProductCategory string `json:"product_category"`
	ProductDescription string `json:"product_description"`
}

// SupplyChainSmartContract defines the smart contract
type SupplyChainSmartContract struct {
	contractapi.Contract
}

// fetchTransactionTimestamp retrieves the current transaction timestamp
func (s *SupplyChainSmartContract) fetchTransactionTimestamp(ctx contractapi.TransactionContextInterface) (string, error) {
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve transaction timestamp: %v", err)
	}
	return time.Unix(timestamp.Seconds, int64(timestamp.Nanos)).Format(time.RFC3339), nil
}

// InitializeLedger adds initial data to the ledger
func (s *SupplyChainSmartContract) InitializeLedger(ctx contractapi.TransactionContextInterface) error {
	timeNow, err := s.fetchTransactionTimestamp(ctx)
	if err != nil {
		return err
	}

	initialProducts := []ProductEntity{
		{ProductID: "prod1", ProductName: "Gaming Laptop", ProductStatus: "Manufactured", CurrentOwner: "TechCorp", CreatedDate: timeNow, UpdatedDate: timeNow, ProductDescription: "A high-performance gaming laptop", ProductCategory: "Electronics"},
		{ProductID: "prod2", ProductName: "5G Smartphone", ProductStatus: "Manufactured", CurrentOwner: "MobileCo", CreatedDate: timeNow, UpdatedDate: timeNow, ProductDescription: "Latest 5G-enabled smartphone", ProductCategory: "Electronics"},
	}

	for _, product := range initialProducts {
		if err := s.saveProduct(ctx, &product); err != nil {
			return err
		}
	}

	return nil
}

// RegisterProduct adds a new product to the ledger
func (s *SupplyChainSmartContract) RegisterProduct(ctx contractapi.TransactionContextInterface, id, name, owner, description, category string) error {
	exists, err := s.CheckProductExistence(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("product with ID %s already exists", id)
	}

	timeNow, err := s.fetchTransactionTimestamp(ctx)
	if err != nil {
		return err
	}

	newProduct := ProductEntity{
		ProductID: id, ProductName: name, ProductStatus: "Manufactured", CurrentOwner: owner, CreatedDate: timeNow, UpdatedDate: timeNow, ProductDescription: description, ProductCategory: category,
	}

	return s.saveProduct(ctx, &newProduct)
}

// ModifyProduct updates existing product details
func (s *SupplyChainSmartContract) ModifyProduct(ctx contractapi.TransactionContextInterface, id, status, owner, description, category string) error {
	productBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("error retrieving product: %v", err)
	}
	if productBytes == nil {
		return fmt.Errorf("product with ID %s does not exist", id)
	}

	var product ProductEntity
	json.Unmarshal(productBytes, &product)

	if status != "" {
		product.ProductStatus = status
	}
	if owner != "" {
		product.CurrentOwner = owner
	}
	if description != "" {
		product.ProductDescription = description
	}
	if category != "" {
		product.ProductCategory = category
	}

	product.UpdatedDate, err = s.fetchTransactionTimestamp(ctx)
	if err != nil {
		return err
	}

	return s.saveProduct(ctx, &product)
}

// TransferOwnership assigns a new owner to the product
func (s *SupplyChainSmartContract) TransferOwnership(ctx contractapi.TransactionContextInterface, id, newOwner string) error {
	return s.ModifyProduct(ctx, id, "", newOwner, "", "")
}

// RetrieveProduct fetches product details based on the product ID
func (s *SupplyChainSmartContract) RetrieveProduct(ctx contractapi.TransactionContextInterface, id string) (*ProductEntity, error) {
	productBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching product details: %v", err)
	}
	if productBytes == nil {
		return nil, fmt.Errorf("product with ID %s not found", id)
	}

	var product ProductEntity
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

// saveProduct is a utility function to add or update a product in the ledger
func (s *SupplyChainSmartContract) saveProduct(ctx contractapi.TransactionContextInterface, product *ProductEntity) error {
	productBytes, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(product.ProductID, productBytes)
}

// CheckProductExistence verifies if a product exists in the ledger
func (s *SupplyChainSmartContract) CheckProductExistence(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	productBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("error checking product existence: %v", err)
	}
	return productBytes != nil, nil
}

// ListAllProducts retrieves all products from the ledger
func (s *SupplyChainSmartContract) ListAllProducts(ctx contractapi.TransactionContextInterface) ([]*ProductEntity, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var allProducts []*ProductEntity
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product ProductEntity
		if err := json.Unmarshal(queryResponse.Value, &product); err != nil {
			return nil, err
		}
		allProducts = append(allProducts, &product)
	}

	return allProducts, nil
}

func main() {
	contract := new(SupplyChainSmartContract)

	chaincode, err := contractapi.NewChaincode(contract)
	if err != nil {
		fmt.Printf("Error creating chaincode instance: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode instance: %s", err.Error())
	}
}