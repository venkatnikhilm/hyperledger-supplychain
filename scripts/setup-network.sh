#!/bin/bash

# ==========================================
# SUPPLY CHAIN BLOCKCHAIN NETWORK SETUP
# ==========================================

set -e

echo "🔗 Supply Chain Blockchain Network Setup"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker and Docker Compose are installed${NC}"
echo ""

# Function to clean up existing network
cleanup() {
    echo "🧹 Cleaning up existing network..."
    cd network
    docker-compose down -v 2>/dev/null || true
    cd ..
    
    # Remove chaincode containers
    docker rm -f $(docker ps -aq --filter "name=dev-peer") 2>/dev/null || true
    
    # Remove chaincode images
    docker rmi -f $(docker images -q --filter "reference=dev-peer*") 2>/dev/null || true
    
    echo -e "${GREEN}✅ Cleanup complete${NC}"
    echo ""
}

# Function to create required directories
create_directories() {
    echo "📁 Creating required directories..."
    mkdir -p crypto-config
    mkdir -p channel-artifacts
    mkdir -p scripts
    echo -e "${GREEN}✅ Directories created${NC}"
    echo ""
}

# Function to start the network
start_network() {
    echo "🚀 Starting Hyperledger Fabric network..."
    cd network
    docker-compose up -d
    cd ..
    
    # Wait for containers to be ready
    echo "⏳ Waiting for containers to start..."
    sleep 10
    
    # Check if containers are running
    if [ $(docker ps -q --filter "name=peer0.org1.example.com" | wc -l) -eq 1 ] && \
       [ $(docker ps -q --filter "name=orderer.example.com" | wc -l) -eq 1 ]; then
        echo -e "${GREEN}✅ Network is running${NC}"
    else
        echo -e "${RED}❌ Network failed to start. Check docker logs.${NC}"
        exit 1
    fi
    echo ""
}

# Function to create channel
create_channel() {
    echo "📺 Creating supply chain channel..."
    
    docker exec cli peer channel create \
        -o orderer.example.com:7050 \
        -c supplychainchannel \
        -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/supplychainchannel.tx \
        2>/dev/null || echo -e "${YELLOW}⚠️  Channel may already exist or configuration missing${NC}"
    
    sleep 3
    
    echo -e "${GREEN}✅ Channel creation attempted${NC}"
    echo ""
}

# Function to join peer to channel
join_channel() {
    echo "🤝 Joining peer to channel..."
    
    docker exec cli peer channel join \
        -b supplychainchannel.block \
        2>/dev/null || echo -e "${YELLOW}⚠️  Peer may already be joined or channel block missing${NC}"
    
    sleep 3
    
    echo -e "${GREEN}✅ Peer join attempted${NC}"
    echo ""
}

# Function to display network status
show_status() {
    echo "📊 Network Status:"
    echo "===================="
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" --filter "network=supply-chain-network"
    echo ""
}

# Function to show next steps
show_next_steps() {
    echo -e "${GREEN}🎉 Network Setup Complete!${NC}"
    echo ""
    echo "📋 Next Steps:"
    echo "=============="
    echo ""
    echo "1️⃣  Enter the CLI container:"
    echo "   ${YELLOW}docker exec -it cli bash${NC}"
    echo ""
    echo "2️⃣  Inside the container, deploy the chaincode:"
    echo "   ${YELLOW}cd /opt/gopath/src/github.com/hyperledger/fabric/peer${NC}"
    echo "   ${YELLOW}# Then follow commands in network/ChaincodeCommands.txt${NC}"
    echo ""
    echo "3️⃣  Quick test - List all products:"
    echo "   ${YELLOW}peer chaincode query --channelID supplychainchannel -n supplychain -c '{\"function\":\"ListAllProducts\",\"Args\":[]}'${NC}"
    echo ""
    echo "4️⃣  View logs if you encounter issues:"
    echo "   ${YELLOW}docker logs peer0.org1.example.com${NC}"
    echo "   ${YELLOW}docker logs orderer.example.com${NC}"
    echo ""
    echo "📖 For detailed usage, see README.md"
    echo ""
}

# Main execution
main() {
    echo "Select an option:"
    echo "1) Fresh setup (cleanup + start network)"
    echo "2) Start network only"
    echo "3) Stop network"
    echo "4) Show network status"
    echo "5) View logs"
    echo ""
    read -p "Enter your choice (1-5): " choice
    echo ""
    
    case $choice in
        1)
            cleanup
            create_directories
            start_network
            # Uncomment these if you have channel configuration files
            # create_channel
            # join_channel
            show_status
            show_next_steps
            ;;
        2)
            start_network
            show_status
            show_next_steps
            ;;
        3)
            echo "🛑 Stopping network..."
            cd network
            docker-compose down
            cd ..
            echo -e "${GREEN}✅ Network stopped${NC}"
            ;;
        4)
            show_status
            ;;
        5)
            echo "📋 Select log to view:"
            echo "1) Peer"
            echo "2) Orderer"
            echo "3) CLI"
            read -p "Enter choice: " log_choice
            
            case $log_choice in
                1) docker logs peer0.org1.example.com ;;
                2) docker logs orderer.example.com ;;
                3) docker logs cli ;;
                *) echo "Invalid choice" ;;
            esac
            ;;
        *)
            echo -e "${RED}Invalid choice. Please run the script again.${NC}"
            exit 1
            ;;
    esac
}

# Run main function
main