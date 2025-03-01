#!/bin/sh

# Wait for PostgreSQL to be ready
until pg_isready -h db -U postgres; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

# Run the migrations
export PGPASSWORD=password psql -h db -U postgres -d habit_tracker -f /migrations/migrations.sql