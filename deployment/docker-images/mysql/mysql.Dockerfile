FROM mysql:8.0

# Set environment variables for MySQL
ENV MYSQL_DATABASE=stadiumhub
ENV MYSQL_USER=stadiumhubuser
ENV MYSQL_PASSWORD=root
ENV MYSQL_ROOT_PASSWORD=root

# Expose MySQL port
EXPOSE 3306

# Optional: Add initialization scripts if needed
# COPY ./init.sql /docker-entrypoint-initdb.d/
