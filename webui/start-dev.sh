#!/bin/bash
# Start backend
echo "Starting Backend..."
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!

# Wait for backend to start
sleep 2

# Start frontend
echo "Starting Frontend..."
cd ../frontend
npm run dev &
FRONTEND_PID=$!

# Handle shutdown
trap "kill $BACKEND_PID $FRONTEND_PID" SIGINT

wait
