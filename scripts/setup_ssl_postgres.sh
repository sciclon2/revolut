#!/bin/bash

# Create directories if they don't exist
mkdir -p certs

# Generate SSL certificates
openssl req -new -x509 -days 365 -nodes -text -out certs/server.crt -keyout certs/server.key -subj "/CN=postgres"
chmod og-rwx certs/server.key

# Create custom PostgreSQL configuration files
cat <<EOL > certs/postgresql.conf
ssl = on
ssl_cert_file = '/var/lib/postgresql/server.crt'
ssl_key_file = '/var/lib/postgresql/server.key'
EOL

cat <<EOL > certs/pg_hba.conf
# TYPE  DATABASE        USER            ADDRESS                 METHOD
hostssl all             all             all                     md5
EOL

echo "SSL setup completed."
