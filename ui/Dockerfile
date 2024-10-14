FROM node:20-alpine AS builder

WORKDIR /app

RUN apk add --no-cache libc6-compat
COPY package.json package-lock.json ./
RUN npm ci
COPY . .

ARG NEXT_PUBLIC_API_ADDR
ARG NODE_ENV

ENV NODE_ENV=${NODE_ENV}
ENV NEXT_PUBLIC_API_ADDR=${NEXT_PUBLIC_API_ADDR}

RUN npm run build

FROM node:20-alpine AS runner

WORKDIR /app

RUN mkdir .next

COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

EXPOSE 3000

CMD ["node", "server.js"]
