# Stage 1: Install dependencies
FROM node:20-alpine AS deps
WORKDIR /app
RUN corepack enable
COPY web/package.json web/pnpm-lock.yaml ./web/
RUN pnpm -C web install --frozen-lockfile

# Stage 2: Build the application
FROM node:20-alpine AS build
WORKDIR /app
RUN corepack enable
COPY --from=deps /app/web/node_modules ./web/node_modules
COPY web ./web

# Pass build-time environment variables
# Note: Next.js inlines NEXT_PUBLIC_ vars at build time
ARG NEXT_PUBLIC_API_BASE_URL
ENV NEXT_PUBLIC_API_BASE_URL=${NEXT_PUBLIC_API_BASE_URL}

RUN pnpm -C web build

# Stage 3: Production Server (Standalone)
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy standalone build
# Next.js "standalone" output creates a minimal set of files in .next/standalone
COPY --from=build --chown=nextjs:nodejs /app/web/.next/standalone ./
COPY --from=build --chown=nextjs:nodejs /app/web/.next/static ./web/.next/static
COPY --from=build --chown=nextjs:nodejs /app/web/public ./web/public

USER nextjs

EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Standalone mode runs via 'node server.js'
CMD ["node", "web/server.js"]
