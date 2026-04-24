# ── Dev stage ────────────────────────────────────────────────────────────────
# Runs the Vite dev server with HMR on port 5173.
# VITE_ vars are read at runtime, so supply them via --env-file / env_file —
# no build args required.
FROM node:20-slim AS dev

WORKDIR /app

COPY package.json ./
RUN npm install --ignore-scripts

# Source is mounted as a volume at runtime for hot-module reload;
# the COPY here is a fallback when running without a mount.
COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host"]

# ── Production build stage ────────────────────────────────────────────────────
FROM node:20-slim AS builder

ARG BUILDKIT_INLINE_CACHE=1

# VITE_ vars are embedded at build time.
# Recommended: use docker compose (reads .env automatically).
# Direct builds: docker build --build-arg VITE_L1_RPC_URL=... .
ARG VITE_L1_RPC_URL
ARG VITE_L2_RPC_URL
ARG VITE_BUILDER_RPC_URL
ARG VITE_CATALYST_RPC_URL
ARG VITE_CHAIN_ID
ARG VITE_L1_CHAIN_NAME
ARG VITE_L1_NATIVE_SYMBOL
ARG VITE_L1_NATIVE_NAME
ARG VITE_L1_NATIVE_LOGO
ARG VITE_EXPLORER_URL
ARG VITE_SAFE_PROXY_FACTORY
ARG VITE_SAFE_SINGLETON
ARG VITE_SAFE_MULTISEND
ARG VITE_SAFE_FALLBACK_HANDLER
ARG VITE_L1_VAULT
ARG VITE_USDC_TOKEN
ARG VITE_USDC_DECIMALS
ARG VITE_SIMPLE_DEX
ARG VITE_L1_BRIDGE
ARG VITE_L2_CHAIN_ID
ARG VITE_L2_VAULT
ARG VITE_L2_USDC_TOKEN
ARG VITE_L1_ROUTER
ARG VITE_L1_DEX_WETH
ARG VITE_L1_DEX_NAME
ARG VITE_WALLETCONNECT_PROJECT_ID

WORKDIR /app

COPY package.json ./
RUN npm install --ignore-scripts

COPY . .

# Write a .env file from build args so Vite picks them up
RUN printf '%s\n' \
    "VITE_L1_RPC_URL=${VITE_L1_RPC_URL}" \
    "VITE_L2_RPC_URL=${VITE_L2_RPC_URL}" \
    "VITE_BUILDER_RPC_URL=${VITE_BUILDER_RPC_URL}" \
    "VITE_CATALYST_RPC_URL=${VITE_CATALYST_RPC_URL}" \
    "VITE_CHAIN_ID=${VITE_CHAIN_ID}" \
    "VITE_L1_CHAIN_NAME=${VITE_L1_CHAIN_NAME}" \
    "VITE_L1_NATIVE_SYMBOL=${VITE_L1_NATIVE_SYMBOL}" \
    "VITE_L1_NATIVE_NAME=${VITE_L1_NATIVE_NAME}" \
    "VITE_L1_NATIVE_LOGO=${VITE_L1_NATIVE_LOGO}" \
    "VITE_EXPLORER_URL=${VITE_EXPLORER_URL}" \
    "VITE_SAFE_PROXY_FACTORY=${VITE_SAFE_PROXY_FACTORY}" \
    "VITE_SAFE_SINGLETON=${VITE_SAFE_SINGLETON}" \
    "VITE_SAFE_MULTISEND=${VITE_SAFE_MULTISEND}" \
    "VITE_SAFE_FALLBACK_HANDLER=${VITE_SAFE_FALLBACK_HANDLER}" \
    "VITE_L1_VAULT=${VITE_L1_VAULT}" \
    "VITE_USDC_TOKEN=${VITE_USDC_TOKEN}" \
    "VITE_USDC_DECIMALS=${VITE_USDC_DECIMALS}" \
    "VITE_SIMPLE_DEX=${VITE_SIMPLE_DEX}" \
    "VITE_L1_BRIDGE=${VITE_L1_BRIDGE}" \
    "VITE_L2_CHAIN_ID=${VITE_L2_CHAIN_ID}" \
    "VITE_L2_VAULT=${VITE_L2_VAULT}" \
    "VITE_L2_USDC_TOKEN=${VITE_L2_USDC_TOKEN}" \
    "VITE_L1_ROUTER=${VITE_L1_ROUTER}" \
    "VITE_L1_DEX_WETH=${VITE_L1_DEX_WETH}" \
    "VITE_L1_DEX_NAME=${VITE_L1_DEX_NAME}" \
    "VITE_WALLETCONNECT_PROJECT_ID=${VITE_WALLETCONNECT_PROJECT_ID}" \
    > .env

RUN npm run build

# Stage 2: Serve static output with nginx
FROM nginx:1.27-alpine AS runtime

# Remove default nginx static assets
RUN rm -rf /usr/share/nginx/html/*

# nginx.conf uses ${BUILDER_URL} which is substituted at container start by
# the nginx image's built-in envsubst mechanism (files in /etc/nginx/templates/
# with a .template suffix are processed and written to /etc/nginx/conf.d/).
COPY nginx.conf /etc/nginx/templates/default.conf.template

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]