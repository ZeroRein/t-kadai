# Stage 1: Build Frontend
FROM node:20 as frontend-builder
WORKDIR /app
# We expect the build context to be the project root
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Setup Nginx
FROM nginx:alpine
COPY --from=frontend-builder /app/dist /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
# Certs will be mounted or copied. If mounted via volume in compose, we don't strictly need COPY here 
# but for a self-contained image it's good. In compose we use volume. 
# However, for the 'build' to succeed without volume, the certs might need to exist or be ignored?
# The COPY nginx/certs below assumes they exist in build context. 
# If we generate them on the host, we can copy them.
COPY nginx/certs /etc/nginx/certs

EXPOSE 80 443
CMD ["nginx", "-g", "daemon off;"]
