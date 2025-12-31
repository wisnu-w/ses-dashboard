# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:18-alpine

WORKDIR /app

# Copy package files for production dependencies
COPY package*.json ./
RUN npm ci --only=production

# Copy built files and server
COPY --from=builder /app/dist ./dist
COPY server.js .

EXPOSE 80

# Run the application
CMD ["npm", "start"]