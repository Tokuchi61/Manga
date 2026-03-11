SELECT 'CREATE DATABASE manga_test'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'manga_test')\gexec
