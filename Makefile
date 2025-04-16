.PHONY: dev dev-air build run-backend start prod deploy-prod

# Development mode - runs Vite dev server and Go server concurrently
dev:
	@echo "Starting development servers..."
	@(cd frontend && npm run dev) & \
	(cd backend && GO_ENV=development go run main.go)

# Development mode with Air hot reloading
dev-air:
	@echo "Starting development servers with Air hot reloading..."
	@(cd frontend && npm run dev) & \
	(cd backend && air)

# Build frontend assets for production
build:
	@echo "Building frontend assets..."
	@cd frontend && npm run build

# Run Go backend server in production mode
run-backend:
	@echo "Starting backend server..."
	@cd backend && GO_ENV=production go run main.go

# Build and start all in production mode
start: build run-backend

# Build and run for production
prod: build
	@echo "Building Go binary..."
	@cd backend && go build -o ../server
	@echo "Starting server..."
	@GO_ENV=production ./server

# Build and prepare files for production deployment
deploy-prod: build
	@echo "Building optimized Go binary..."
	@cd backend && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o ../server
	
	@echo "Creating deployment package..."
	@mkdir -p deploy
	@cp -r frontend/dist deploy/
	@cp -r backend/templates deploy/
	@cp server deploy/
	
	@echo "Creating systemd service file..."
	@echo "[Unit]" > deploy/server.service
	@echo "Description=Go React Islands Server" >> deploy/server.service
	@echo "After=network.target" >> deploy/server.service
	@echo "" >> deploy/server.service
	@echo "[Service]" >> deploy/server.service
	@echo "Type=simple" >> deploy/server.service
	@echo "User=www-data" >> deploy/server.service
	@echo "WorkingDirectory=/opt/server" >> deploy/server.service
	@echo "ExecStart=/opt/server/server" >> deploy/server.service
	@echo "Environment=GO_ENV=production" >> deploy/server.service
	@echo "Restart=on-failure" >> deploy/server.service
	@echo "" >> deploy/server.service
	@echo "[Install]" >> deploy/server.service
	@echo "WantedBy=multi-user.target" >> deploy/server.service
	
	@echo "Deployment package ready at ./deploy"
	@echo "To install on production server, copy the files to /opt/server and run:"
	@echo "  sudo cp ./deploy/server.service /etc/systemd/system/"
	@echo "  sudo systemctl daemon-reload"
	@echo "  sudo systemctl enable server.service"
	@echo "  sudo systemctl start server.service"