# test-social-tournament
Backend developer coding task "Social tournament service"

Basic implementation of API due to requirements [BackendTestTaskWithExamples.pdf](BackendTestTaskWithExamples.pdf)

# Setup

- Create PostgreSQL database and user
- Run script [init.sql](init.sql) in SQL editor to create required objects
- Edit constants in [main.go](main.go) to change default values:

```
SERVICE_PORT  = ":8088"

POSTGRES_HOST = "127.0.0.1"
POSTGRES_PORT = 5432
POSTGRES_USER = "test"
POSTGRES_PASS = "test"
POSTGRES_BASE = "test"
```

 - Compile project and run binary