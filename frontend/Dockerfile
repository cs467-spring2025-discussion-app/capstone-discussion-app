# Build stage
FROM node:18 as builder

WORKDIR /app
COPY . .

ENV NODE_ENV=production

RUN npm install --include=dev

# doing this will pickup the .env.production config
RUN npm run build -- --mode production

# Serve stage
FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]