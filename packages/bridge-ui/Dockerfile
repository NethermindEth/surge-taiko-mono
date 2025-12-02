# syntax=docker/dockerfile:1.6

# Step 1: Build the application
FROM node:20-slim AS builder

ARG BUILDKIT_INLINE_CACHE=1

# Install pnpm
RUN npm install -g pnpm@latest

WORKDIR /app

# Copy dependency files first for better layer caching
COPY package.json pnpm-lock.yaml* ./

# Install dependencies
RUN pnpm install

# Copy source code
COPY . .

# Copy and export config
RUN cp .env.example .env && \
    pnpm export:config

# Build the application
RUN pnpm build

# Step 2: Set up the final production environment
FROM node:20-slim AS runtime

# Install pnpm globally
RUN npm install -g pnpm@latest && \
    npm cache clean --force

WORKDIR /app

# Create non-root user
RUN groupadd -r appuser && useradd -r -g appuser appuser

# Copy only necessary files for production
COPY --from=builder --chown=appuser:appuser /app/package.json ./
COPY --from=builder --chown=appuser:appuser /app/pnpm-lock.yaml ./

# Install production dependencies only
RUN pnpm install --prod && \
    pnpm store prune && \
    npm cache clean --force

# Copy the build output
COPY --from=builder --chown=appuser:appuser /app/build ./build

USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD node -e "require('http').get('http://localhost:3000', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"

# Start the SvelteKit app
CMD ["node", "build"]
