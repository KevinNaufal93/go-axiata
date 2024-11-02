# GO CRUD Server
Fun CRUD exploration using vanilla GO By Kevin Naufal Permana

# How to run
1. Please use "go get" on these items first
    "github.com/lib/pq"
    "github.com/google/uuid"
2. Make sure your local postgresql live and running
3. create schema to store the tables
4. Run all the migration scripts on your local postgresql schema
5. make sure the configs const in db.go match your local configuration
6. Type "go run ." to start local server, if you intent to make changes make sure to restart
