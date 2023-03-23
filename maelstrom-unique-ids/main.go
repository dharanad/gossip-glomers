package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
)

func main() {
	dbUrl := getConnectionString("postgres", "password", "localhost", "dist_sys", 5432)
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Fatalf("Error opening database connection. %v", err)
	}
	if pErr := db.Ping(); pErr != nil {
		log.Fatalf("Error connecting to db. %v", err)
	}
	n := maelstrom.NewNode()
	var idGenerator IdGenerator = NewUniqueIdGeneratorService(db)
	n.Handle("generate", NewIdGeneratorHandler(idGenerator, n))
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func getConnectionString(username, password, host, dbName string, port int) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbName)
}

func NewIdGeneratorHandler(idGenerator IdGenerator, n *maelstrom.Node) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "generate_ok"
		body["id"] = idGenerator.Generate()
		return n.Reply(msg, body)
	}
}

type IdGenerator interface {
	Generate() int64
}

type UniqueIdGeneratorService struct {
	db *sql.DB
}

func NewUniqueIdGeneratorService(db *sql.DB) *UniqueIdGeneratorService {
	return &UniqueIdGeneratorService{db: db}
}

func (u *UniqueIdGeneratorService) Generate() int64 {
	panic("implement me")
}
