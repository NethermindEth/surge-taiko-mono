# Build stage
FROM node:20-alpine as builder

# Install pnpm
RUN npm install -g pnpm

# Install dependencies required for node-gyp
RUN apk add --no-cache python3 make g++ rng-tools

WORKDIR /app

# Copy package files
COPY package.json ./

# Install dependencies
RUN pnpm install --no-frozen-lockfile

# Copy source code and config files
COPY . .

# Copy .env example if not exists
RUN cp -n .env.example .env || true

# Export configs to .env
RUN pnpm export:config

# Build the application
RUN pnpm run build

CMD ["pnpm", "run", "preview","--host"]
