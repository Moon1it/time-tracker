DROP TRIGGER IF EXISTS set_users_updated_at ON users;

DROP FUNCTION IF EXISTS update_updated_at_column;

DROP TABLE IF EXISTS task_histories;

DROP TABLE IF EXISTS tasks;

DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";
