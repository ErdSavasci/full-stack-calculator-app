# --- Stage 1: The Build Environment ---
FROM node:20-alpine AS builder

WORKDIR /app

# 1. Install dependencies first (for better caching)
COPY package.json package-lock.json ./
RUN npm install

# 2. Copy the rest of the frontend code
COPY . .

# 3. Build the production-ready static files
# This creates a folder named 'build' (or 'dist' if using Vite)
RUN npm run build

# --- Stage 2: The Production Server ---
FROM nginx:stable-alpine

# 4. Copy the static files from the builder stage to Nginx's default folder
COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

# 5. Expose the port Nginx runs on
EXPOSE 80

# 6. Start Nginx
CMD ["nginx", "-g", "daemon off;"]