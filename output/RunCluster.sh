#!/bin/bash
export DATABASE_URL="postgresql://postgres:REDACTED@localhost:5432/encvault"

# Function to handle cleanup on script exit
cleanup() {
    echo "Shutting down all servers..."
    # Kill all background jobs
    kill $(jobs -p)
    echo "All servers shut down."
}

# Trap SIGINT (CTRL-C) and call the cleanup function
trap cleanup SIGINT

# Start the registry server
echo "Starting registry server..."
nohup ./run ./configurations/RegistryService.yaml &

# Wait for the registry server to start
sleep 5

# Start the keymaker server
echo "Starting keymaker server..."
nohup ./run ./configurations/KeyMakerService.yaml &

# Wait for the keymaker server to start
sleep 5

# Start the tokengenerator server
echo "Starting tokengenerator server..."
nohup ./run ./configurations/TokenGeneratorService.yaml &

# Wait for the tokengenerator server to start
sleep 5

# Start the user_auth server
echo "Starting user_auth server..."
nohup ./run ./configurations/UserAuthService.yaml &

# Wait for the user_auth server to start
sleep 5

# Start the file upload server
echo "Starting file upload server..."
nohup ./run ./configurations/FileUploadService.yaml &

echo "All servers started successfully."

# Wait indefinitely to keep the script running
wait