FROM nginx:alpine

# Remove default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom nginx config
COPY nginx/nginx.conf /etc/nginx/conf.d/

# Create directory for SSL certificates
RUN mkdir -p /etc/nginx/ssl

# Copy SSL certificates (you'll need to provide these)
COPY nginx/ssl/cert.pem /etc/nginx/ssl/
COPY nginx/ssl/key.pem /etc/nginx/ssl/ 
