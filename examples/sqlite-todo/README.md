# SQLite Todo API Example 📝

A simple Todo API demonstrating ReboloLang's SQLite database support.

## Features

- ✅ Full CRUD operations
- ✅ SQLite database with WAL mode
- ✅ RESTful API design
- ✅ Automatic migrations
- ✅ Clean architecture

## Running the Example

1. Navigate to the example directory:
```bash
cd examples/sqlite-todo
```

2. Run the application:
```bash
go run main.go
```

The API will start on `http://localhost:3000`

## API Endpoints

### Get All Todos
```bash
curl http://localhost:3000/todos
```

### Create a Todo
```bash
curl -X POST http://localhost:3000/todos \
  -d "title=Learn ReboloLang"
```

### Get a Single Todo
```bash
curl http://localhost:3000/todos/1
```

### Update a Todo
```bash
curl -X PUT http://localhost:3000/todos/1 \
  -d "title=Master ReboloLang" \
  -d "completed=true"
```

### Delete a Todo
```bash
curl -X DELETE http://localhost:3000/todos/1
```

## Database

The application uses SQLite with:
- **Standard database/sql** package from Go
- **WAL mode** for better concurrency
- **Shared cache** for performance
- **Automatic schema creation** using SQL

Database file: `todos.db` (created automatically)

> The framework doesn't impose any ORM - you can use GORM, sqlx, or any other ORM you prefer!

## Switching to PostgreSQL or MySQL

Just update `config.yml`:

### PostgreSQL
```yaml
database:
  driver: "postgres"
  url: "postgres://user:pass@localhost/todos?sslmode=disable"
```

### MySQL
```yaml
database:
  driver: "mysql"
  url: "user:pass@tcp(localhost:3306)/todos?parseTime=true"
```

The code works exactly the same! 🎉

## Code Highlights

### Model Definition
```go
type Todo struct {
    ID        int64     `bun:"id,pk,autoincrement"`
    Title     string    `bun:"title,notnull"`
    Completed bool      `bun:"completed,notnull,default:false"`
    CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
```

### Database Access
```go
db := app.DB()  // Get *sql.DB instance

// Use standard database/sql
rows, err := db.QueryContext(ctx, 
    "SELECT id, title, completed, created_at FROM todos ORDER BY created_at DESC")
```

### Migration
```go
_, err := db.NewCreateTable().
    Model((*Todo)(nil)).
    IfNotExists().
    Exec(ctx)
```

## Learn More

- [Database Support Documentation](../../DATABASE_SUPPORT.md)
- [ReboloLang Architecture](../../ARCHITECTURE.md)

