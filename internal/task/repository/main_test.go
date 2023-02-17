package repository

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AssylzhanZharzhanov/axxonsoft-test-service/pkg/database/postgres"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	database  *gorm.DB
	testSQLDB *sql.DB
)

func TestMain(m *testing.M) {

	var (
		databaseName = "test"
		username     = "admin"
		password     = "test"
	)

	// Define if we will use reaper to clean up resources.
	//
	skipReaper := os.Getenv("TESTCONTAINERS_RYUK_DISABLED") != ""
	log.Printf("flag SkipReaper = %v", skipReaper)

	// Setup Database.
	//
	ctx := context.Background()
	exposedPort := "5432"
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.2",
		ExposedPorts: []string{exposedPort},
		WaitingFor:   wait.ForListeningPort(nat.Port(exposedPort)),
		Env: map[string]string{
			"POSTGRES_DB":       databaseName,
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": password,
		},
		SkipReaper: skipReaper, // if you have problems with reaper rights just skip this
	}

	// Start test container.
	//
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container %v", err)
	}

	defer func() { _ = container.Terminate(ctx) }()

	// Setup test dependencies.
	//
	ip, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to start container host %v", err)
	}
	port, err := container.MappedPort(ctx, nat.Port(exposedPort))
	if err != nil {
		log.Fatalf("Failed to start container port %v", err)
	}
	postgreAddress := fmt.Sprintf("%s:%s", ip, port.Port())
	fmt.Printf("postgresql address: %s\n", postgreAddress)

	postgreDSN := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", username, password, ip, port.Port(), databaseName)
	fmt.Printf("postgresql DSN: %s\n", postgreDSN)

	// Initialize database connection.
	//
	database, err = postgres.NewConnection(postgreDSN)
	if err != nil {
		log.Fatalf("Cannot create database %v", err)
	}

	testSQLDB, err = database.DB()
	if err != nil {
		log.Fatalf("Cannot establish connection with database %v", err)
	}
	if err := testSQLDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database %v", err)
	}
	defer func() { _ = testSQLDB.Close() }()

	exitVal := m.Run()
	os.Exit(exitVal)
}

func runTestSetup(t *testing.T, setup func() error) {
	err := setup()
	if err != nil {
		t.Errorf("Setup faild with error: %s", err)
		t.FailNow()
	}
}

func executeTestStatement(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.
	}()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	queryResult, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return queryResult, err
}

func getVersion(match []string, i int, filename string) (int, error) {
	if len(match) > i {
		return strconv.Atoi(match[i])
	}
	return 0, fmt.Errorf("incorrect file versioning pattern: %s", filename)
}

func setupDatabase() error {

	type version struct {
		Major int
		Minor int
		Patch int
	}
	type migration struct {
		Version version
		Content []string
	}

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	migrationsDir := strings.ReplaceAll(currentDir, "internal/task/repository", "migrations/postgresql")
	pattern := regexp.MustCompile(`(.*\/)?((\d+).(\d+).(\d+).*\.sql)$`)
	var migrations []*migration
	err = filepath.Walk(migrationsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking filepath: %v", err)
		}
		if info != nil && info.IsDir() {
			return nil
		}
		match := pattern.FindStringSubmatch(path)
		if match == nil {
			return fmt.Errorf("file is ignored, naming pattern is wrong: %s", path)
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var contents []string
		for scanner.Scan() {
			contents = append(contents, scanner.Text())
		}
		major, err := getVersion(match, 3, path)
		if err != nil {
			return err
		}
		minor, err := getVersion(match, 4, path)
		if err != nil {
			return err
		}
		patch, err := getVersion(match, 5, path)
		if err != nil {
			return err
		}
		migration := &migration{
			Version: version{
				Major: major,
				Minor: minor,
				Patch: patch,
			},
			Content: contents,
		}
		migrations = append(migrations, migration)
		return nil
	})
	if err != nil {
		return err
	}

	// Sort migrations in the correct order.
	//
	sort.Slice(migrations, func(i, j int) bool {
		if migrations[i].Version.Major == migrations[j].Version.Major {
			if migrations[i].Version.Minor == migrations[j].Version.Minor {
				return migrations[i].Version.Patch < migrations[j].Version.Patch
			}
			return migrations[i].Version.Minor < migrations[j].Version.Minor
		}
		return migrations[i].Version.Major < migrations[j].Version.Major
	})

	for _, m := range migrations {
		statements, err := parseStatements(m.Content)
		if err != nil {
			return err
		}
		for _, stmt := range statements {
			_, err := executeTestStatement(testSQLDB, stmt)
			if err != nil {
				return fmt.Errorf("execution failed for version: %v err: %v, statement: %v", m.Version, err, stmt)
			}
		}
	}

	return nil
}

func purgeDatabase(db *sql.DB) {
	_, err := executeTestStatement(db, "DROP SCHEMA public CASCADE;")
	if err != nil {
		panic(fmt.Sprintf("Failed to cleanup the database, err: %v", err))
	}
	_, err = executeTestStatement(db, "CREATE SCHEMA public;")
	if err != nil {
		panic(fmt.Sprintf("Failed to cleanup the database, err: %v", err))
	}
}

func parseStatements(lines []string) (queries []string, err error) {
	query := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			continue
		}
		query += line + "\n"
		if strings.HasSuffix(line, ";") {
			queries = append(queries, query)
			query = ""
		}
	}
	return queries, nil
}
