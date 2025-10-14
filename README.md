# 🔗 Hyperledger Fabric Supply Chain Blockchain

A production-ready blockchain solution for supply chain management using Hyperledger Fabric and Go smart contracts. Track products from manufacturing to delivery with transparent, immutable records.

## 📋 Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Quick Start](#quick-start)
- [Detailed Setup](#detailed-setup)
- [Usage Examples](#usage-examples)
- [API Reference](#api-reference)
- [Troubleshooting](#troubleshooting)
- [Development](#development)

---

## 🎯 Overview

This project implements a permissioned blockchain network for supply chain management. It allows multiple organizations to:
- Register products on a shared, tamper-proof ledger
- Transfer ownership across the supply chain
- Track product status and history
- Query product information in real-time
- Maintain complete audit trails

**Key Benefits:**
- ✅ Transparency across all parties
- ✅ Immutable transaction history
- ✅ Reduced disputes and fraud
- ✅ Real-time product tracking
- ✅ Trustless verification

---

## ✨ Features

### Smart Contract Functions
- **InitializeLedger** - Populate ledger with sample data
- **RegisterProduct** - Add new products to the blockchain
- **ModifyProduct** - Update product status, description, or category
- **TransferOwnership** - Change product ownership
- **RetrieveProduct** - Query specific product details
- **CheckProductExistence** - Verify if a product exists
- **ListAllProducts** - Get all products in the supply chain

### Technical Features
- Timestamp tracking (created/updated dates)
- Unique product ID validation
- Partial update support
- Error handling and validation
- Range query support

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Application Layer                  │
│              (CLI / API / Web Interface)            │
└─────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────┐
│              Smart Contract (Chaincode)             │
│         Supply Chain Business Logic (Go)            │
└─────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────┐
│            Hyperledger Fabric Network               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐         │
│  │ Orderer  │  │  Peer0   │  │ CouchDB  │         │
│  │          │  │  Org1    │  │          │         │
│  └──────────┘  └──────────┘  └──────────┘         │
└─────────────────────────────────────────────────────┘
```

**Components:**
- **Orderer**: Orders transactions and creates blocks
- **Peer**: Maintains the ledger and executes chaincode
- **CouchDB**: World state database (optional)
- **Smart Contract**: Business logic written in Go
- **CLI**: Command-line interface for network interaction

---

## 📦 Prerequisites

### Required Software
- **Docker** (v20.10+) - [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose** (v2.0+) - Usually included with Docker Desktop
- **Go** (v1.20+) - [Install Go](https://golang.org/doc/install)
- **Git** - For cloning the repository

### System Requirements
- 8GB RAM minimum (16GB recommended)
- 20GB free disk space
- Linux/MacOS/Windows with WSL2

### Verify Installation
```bash
docker --version
docker-compose --version
go version
git --version
```

---

## 📁 Project Structure

```
supply-chain-blockchain/
│
├── chaincode/
│   └── smartcontract.go          # Smart contract code
│
├── network/
│   ├── docker-compose.yaml       # Network configuration
│   └── ChaincodeCommands.txt     # Deployment commands
│
├── scripts/
│   ├── setup-network.sh          # Network setup script
│   └── deploy-chaincode.sh       # Chaincode deployment script
│
├── crypto-config/                # Certificate files (generated)
├── channel-artifacts/            # Channel configs (generated)
│
└── README.md                     # This file
```

---

## 🚀 Quick Start

### Step 1: Clone and Navigate
```bash
git clone <your-repo-url>
cd supply-chain-blockchain
```

### Step 2: Generate Crypto Materials
```bash
# Create directories
mkdir -p crypto-config channel-artifacts

# For a quick test, you can use the fabric test network to generate these
# Or use cryptogen tool (requires Hyperledger Fabric binaries)
```

### Step 3: Start the Network
```bash
cd network
docker-compose up -d
```

### Step 4: Verify Network is Running
```bash
docker ps
```
You should see containers: `orderer.example.com`, `peer0.org1.example.com`, `couchdb0`, `cli`

### Step 5: Deploy Chaincode
```bash
# Enter the CLI container
docker exec -it cli bash

# Inside the container, run these commands:
cd /opt/gopath/src/github.com/hyperledger/fabric/peer

# Follow commands from network/ChaincodeCommands.txt
```

### Step 6: Test the Smart Contract
```bash
# Initialize ledger
peer chaincode invoke \
    -o orderer.example.com:7050 \
    --channelID supplychainchannel \
    -n supplychain \
    -c '{"function":"InitializeLedger","Args":[]}'

# Query all products
peer chaincode query \
    --channelID supplychainchannel \
    -n supplychain \
    -c '{"function":"ListAllProducts","Args":[]}'
```

---

## 🔧 Detailed Setup

### 1. Network Configuration

Before starting, ensure you have the required crypto materials:

```bash
# Using Fabric test network (recommended for development)
curl -sSL https://bit.ly/2ysbOFE | bash -s

# This downloads Fabric binaries and Docker images
# Then generate crypto materials:
./bin/cryptogen generate --config=./crypto-config.yaml
```

### 2. Create Channel

```bash
# Generate genesis block
./bin/configtxgen -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block

# Create channel transaction
./bin/configtxgen -profile Channel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID supplychainchannel
```

### 3. Start Network

```bash
docker-compose -f network/docker-compose.yaml up -d
```

### 4. Create and Join Channel

```bash
docker exec -it cli bash

# Inside CLI container:
peer channel create -o orderer.example.com:7050 -c supplychainchannel -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/channel.tx

peer channel join -b supplychainchannel.block
```

### 5. Deploy Chaincode

Follow the commands in `network/ChaincodeCommands.txt` step by step.

---

## 💡 Usage Examples

### Register a New Product

```bash
peer chaincode invoke \
    -o orderer.example.com:7050 \
    --channelID supplychainchannel \
    -n supplychain \
    -c '{
        "function":"RegisterProduct",
        "Args":[
            "LAPTOP001",
            "Gaming Laptop Pro",
            "TechManufacturing Inc",
            "High-performance laptop with RTX 4080",
            "Electronics"
        ]
    }'
```

**Result:** Product is permanently recorded on the blockchain with timestamp and owner.

---

### Transfer Product Ownership

```bash
# Manufacturer → Distributor
peer chaincode invoke \
    -o orderer.example.com:7050 \
    --channelID supplychainchannel \
    -n supplychain \
    -c '{
        "function":"TransferOwnership",
        "Args":["LAPTOP001", "GlobalDistributors LLC"]
    }'

# Distributor → Retailer
peer chaincode invoke \
    -o orderer.example.com:7050 \
    --channelID supplychainchannel \
    -n supplychain \
    -c '{
        "function":"TransferOwnership",
        "Args":["LAPTOP001", "BestBuy"]
    }'

# Retailer → Customer
peer chaincode invoke \
    -o orderer.example.com:7050 \
    --channelID supplychainchannel \
    -n supplychain \
    -c '{
        "function":"TransferOwnership",
        "Args":["LAPTOP001", "John Doe"]
    }'
```

---

### Update Product Status

```bash
peer chaincode invoke \
    -o orderer.example.com:7050 \
    --channelID supplychainchannel \
    -n supplychain \
    -c '{
        "function":"ModifyProduct",
        "Args":[
            "LAPTOP001",
            "In Transit",
            "",
            "",
            ""
        ]
    }'
```

**Tip:** Use empty strings `""` for fields you don't want to update.

---

### Query Product Information

```bash
peer chaincode query \
    --channelID supplychainchannel \
    -n supplychain \
    -c '{"function":"RetrieveProduct","Args":["LAPTOP001"]}'
```

**Sample Response:**
```json
{
    "product_id": "LAPTOP001",
    "product_name": "Gaming Laptop Pro",
    "product_status": "In Transit",
    "current_owner": "GlobalDistributors LLC",
    "created_date": "2025-10-14T10:30:00Z",
    "updated_date": "2025-10-14T11:45:00Z",
    "product_category": "Electronics",
    "product_description": "High-performance laptop with RTX 4080"
}
```

---

### List All Products

```bash
peer chaincode query \
    --channelID supplychainchannel \
    -n supplychain \
    -c '{"function":"ListAllProducts","Args":[]}'
```

---

## 📚 API Reference

### RegisterProduct
**Description:** Register a new product on the blockchain  
**Parameters:**
- `id` (string): Unique product identifier
- `name` (string): Product name
- `owner` (string): Initial owner
- `description` (string): Product description
- `category` (string): Product category

**Returns:** Success/error message

---

### ModifyProduct
**Description:** Update existing product details  
**Parameters:**
- `id` (string): Product ID
- `status` (string): New status (or "" to skip)
- `owner` (string): New owner (or "" to skip)
- `description` (string): New description (or "" to skip)
- `category` (string): New category (or "" to skip)

**Returns:** Success/error message

---

### TransferOwnership
**Description:** Change product owner  
**Parameters:**
- `id` (string): Product ID
- `newOwner` (string): New owner name

**Returns:** Success/error message

---

### RetrieveProduct
**Description:** Get product details  
**Parameters:**
- `id` (string): Product ID

**Returns:** ProductEntity JSON object

---

### CheckProductExistence
**Description:** Check if product exists  
**Parameters:**
- `id` (string): Product ID

**Returns:** Boolean (true/false)

---

### ListAllProducts
**Description:** Get all products in ledger  
**Parameters:** None

**Returns:** Array of ProductEntity objects

---

## 🐛 Troubleshooting

### Network Won't Start

**Problem:** `docker-compose up` fails  
**Solution:**
```bash
# Stop all containers
docker-compose down

# Remove volumes
docker volume prune

# Restart
docker-compose up -d
```

---

### Chaincode Installation Fails

**Problem:** `peer lifecycle chaincode install` error  
**Solution:**
- Check Go module path is correct
- Ensure chaincode has `go.mod` file
- Verify peer container can access chaincode directory

```bash
# Inside CLI container:
ls -la /opt/gopath/src/github.com/chaincode
```

---

### Transaction Timeout

**Problem:** `timeout expired while waiting for transaction`  
**Solution:**
- Increase timeout: `export CORE_CHAINCODE_EXECUTETIMEOUT=300s`
- Check peer and orderer logs
- Verify network connectivity

---

### Product Already Exists Error

**Problem:** Trying to register duplicate product ID  
**Solution:** Use `CheckProductExistence` before registration or choose a different ID

---

### View Logs

```bash
# Peer logs
docker logs peer0.org1.example.com

# Orderer logs
docker logs orderer.example.com

# Chaincode logs
docker logs <chaincode_container>
```

---

## 🔨 Development

### Modifying the Smart Contract

1. Edit `chaincode/smartcontract.go`
2. Increment version number in deployment commands
3. Package and install new version
4. Approve and commit (increment sequence number)

### Adding New Functions

```go
func (s *SupplyChainSmartContract) YourNewFunction(
    ctx contractapi.TransactionContextInterface,
    param1 string,
    param2 string,
) error {
    // Your logic here
    return nil
}
```

### Testing

```bash
# Unit tests (create test files)
go test ./chaincode/...

# Integration testing
peer chaincode invoke ... # Test each function
```

---

## 📄 License

MIT License - Feel free to use this project for learning or commercial purposes.

---

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

---

## 📞 Support

- **Issues:** Open a GitHub issue
- **Documentation:** [Hyperledger Fabric Docs](https://hyperledger-fabric.readthedocs.io/)
- **Community:** [Hyperledger Discord](https://discord.com/invite/hyperledger)

---

## 🎓 Learning Resources

- [Hyperledger Fabric Documentation](https://hyperledger-fabric.readthedocs.io/)
- [Blockchain Basics](https://hyperledger-fabric.readthedocs.io/en/latest/blockchain.html)
- [Go Tutorial](https://go.dev/tour/)

---

**Built with ❤️ using Hyperledger Fabric**