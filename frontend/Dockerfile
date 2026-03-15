# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

# Run stage
FROM nginx:stable-alpine

COPY --from=builder /app/dist /usr/share/nginx/html

# Custom nginx config to handle SPA routing
# We use the built-in variable substitution for Nginx on Alpine
RUN echo 'server { \
    listen 80; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html index.htm; \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

# Limit worker processes to 2
RUN sed -i 's/worker_processes  auto;/worker_processes  2;/g' /etc/nginx/nginx.conf

EXPOSE 80

# Railway will map the public port to 80 if we set it in the UI.
# If they provide a $PORT variable, we should listen on it.
CMD ["/bin/sh", "-c", "sed -i \"s/listen 80;/listen ${PORT:-80};/g\" /etc/nginx/conf.d/default.conf && exec nginx -g 'daemon off;'"]
