# Database Schema

## Overview
The application uses PostgreSQL as its primary database. All tables use UUID as primary keys and include timestamp fields for auditing.

## Tables

### messages
Stores message records with content and timestamps.

```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to automatically update updated_at timestamp
CREATE TRIGGER update_messages_updated_at
    BEFORE UPDATE ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

#### Indexes
- `messages_created_at_idx`: Index on created_at for efficient sorting
- `messages_updated_at_idx`: Index on updated_at for efficient sorting

## Functions

### update_updated_at_column()
Trigger function to automatically update the updated_at timestamp.

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
```

## Migrations
Database migrations are managed using golang-migrate. Migration files are stored in the `/migrations` directory.

### Migration Files
- `000001_create_messages_table.up.sql`: Creates the messages table
- `000001_create_messages_table.down.sql`: Drops the messages table

### Running Migrations
```bash
# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down

# Create new migration
make migrate-create name=add_new_table
```

## Backup and Restore
The database can be backed up using pg_dump:

```bash
# Backup
pg_dump -h localhost -U postgres -d messagedb > backup.sql

# Restore
psql -h localhost -U postgres -d messagedb < backup.sql
```

## Performance Considerations
1. Indexes are created on frequently queried columns
2. UUIDs are used instead of sequential IDs for better distribution
3. Timestamp fields use WITH TIME ZONE for proper timezone handling
4. Trigger-based updated_at updates for consistency
