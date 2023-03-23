package main

import (
	"context"
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
	var idGenerator IdGenerator = NewUniqueIdGeneratorService(NewIdsDal(db))
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
		body["id"] = idGenerator.GetId()
		return n.Reply(msg, body)
	}
}

type IdGenerator interface {
	GetId() int64
}

type UniqueIdGeneratorService struct {
	dal *IdsDal
}

func NewUniqueIdGeneratorService(dal *IdsDal) *UniqueIdGeneratorService {
	return &UniqueIdGeneratorService{dal: dal}
}

func (u *UniqueIdGeneratorService) GetId() int64 {
	id, err := u.dal.GetId(context.TODO())
	if err != nil {
		log.Fatalf("error generating id. %v", err)
	}
	return id
}

type IdsDal struct {
	db *sql.DB
}

func NewIdsDal(db *sql.DB) *IdsDal {
	return &IdsDal{db: db}
}

func (d *IdsDal) GetId(ctx context.Context) (int64, error) {
	var id int64
	err := d.db.QueryRowContext(ctx, "insert into ids(sudo_val) values('a') returning id").Scan(&id)
	if err == sql.ErrNoRows {
		return -1, err
	} else if err != nil {
		return -1, err
	}
	return id, nil
}

func (d *IdsDal) GetStringId(ctx context.Context) (string, error) {
	var id string
	err := d.db.QueryRowContext(ctx, "select uuid_in(md5(random()::text || random()::text)::cstring) as id").Scan(&id)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", err
	}
	return id, nil
}
