package postgres

import (
	"database/sql"
	"fmt"

	"mini-asm/internal/config"
	"mini-asm/internal/model"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresStorage implements the Storage interface using PostgreSQL
// This is a concrete implementation that can be swapped with MemoryStorage
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new PostgreSQL storage instance
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func NewPostgresStorageFromConfig(config *config.PostgresConfig) (*PostgresStorage, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &PostgresStorage{db: db}, nil
}

// Create inserts a new asset into the database
func (p *PostgresStorage) Create(asset *model.Asset) error {
	query := `
		INSERT INTO assets (id, name, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := p.db.Exec(
		query,
		asset.ID,
		asset.Name,
		asset.Type,
		asset.Status,
		asset.CreatedAt,
		asset.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}

	return nil
}

// GetAll retrieves all assets from the database
func (p *PostgresStorage) GetAll() ([]*model.Asset, error) {
	query := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		ORDER BY created_at DESC
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		asset := &model.Asset{}
		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Type,
			&asset.Status,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assets, nil
}

// GetByID retrieves a single asset by its ID
func (p *PostgresStorage) GetByID(id string) (*model.Asset, error) {
	query := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		WHERE id = $1
	`

	asset := &model.Asset{}
	err := p.db.QueryRow(query, id).Scan(
		&asset.ID,
		&asset.Name,
		&asset.Type,
		&asset.Status,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	return asset, nil
}

// Update modifies an existing asset in the database
func (p *PostgresStorage) Update(id string, asset *model.Asset) error {
	query := `
		UPDATE assets
		SET name = $1, type = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := p.db.Exec(
		query,
		asset.Name,
		asset.Type,
		asset.Status,
		asset.UpdatedAt,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return model.ErrNotFound
	}

	return nil
}

// Delete removes an asset from the database
func (p *PostgresStorage) Delete(id string) error {
	query := `DELETE FROM assets WHERE id = $1`

	result, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return model.ErrNotFound
	}

	return nil
}

// Filter returns assets matching the given criteria
func (p *PostgresStorage) Filter(assetType, status string) ([]*model.Asset, error) {
	query := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		WHERE 1=1
	`

	var args []interface{}
	argCount := 1

	if assetType != "" {
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, assetType)
		argCount++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	query += " ORDER BY created_at DESC"

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to filter assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		asset := &model.Asset{}
		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Type,
			&asset.Status,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assets, nil
}

// Search finds assets by partial name match
func (p *PostgresStorage) Search(query string) ([]*model.Asset, error) {
	sqlQuery := `
		SELECT id, name, type, status, created_at, updated_at
		FROM assets
		WHERE name ILIKE $1
		ORDER BY created_at DESC
	`

	// Add wildcards for partial matching
	searchPattern := "%" + query + "%"

	rows, err := p.db.Query(sqlQuery, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.Asset
	for rows.Next() {
		asset := &model.Asset{}
		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Type,
			&asset.Status,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return assets, nil
}

func (p *PostgresStorage) GetStatistics() (*model.StatisticsResponse, error) {
	var total int
	err := p.db.QueryRow("SELECT COUNT(*) FROM assets").Scan(&total)
	if err != nil{
		return nil, fmt.Errorf("failed to count total assets: %w", err)
	}

	//Count theo type
	byTypeRows, err := p.db.Query(`
        SELECT type, COUNT(*) as count
        FROM assets
        GROUP BY type
    `)
    if err != nil {
        return nil, fmt.Errorf("failed to group by type: %w", err)
    }
    defer byTypeRows.Close()

	byType := make(map[string]int)
    for byTypeRows.Next() {
        var typeVal string
        var count int
        if err := byTypeRows.Scan(&typeVal, &count); err != nil {
            return nil, fmt.Errorf("failed to scan type row: %w", err)
        }
        byType[typeVal] = count
    }

	if err := byTypeRows.Err(); err != nil {
        return nil, fmt.Errorf("error reading type rows: %w", err)
    }

	//count theo status
	byStatusRows, err := p.db.Query(`
        SELECT status, COUNT(*) as count
        FROM assets
        GROUP BY status
    `)
    if err != nil {
        return nil, fmt.Errorf("failed to group by status: %w", err)
    }
    defer byStatusRows.Close()

    byStatus := make(map[string]int)
    for byStatusRows.Next() {
        var statusVal string
        var count int
        if err := byStatusRows.Scan(&statusVal, &count); err != nil {
            return nil, fmt.Errorf("failed to scan status row: %w", err)
        }
        byStatus[statusVal] = count
    }
    
    if err := byStatusRows.Err(); err != nil {
        return nil, fmt.Errorf("error reading status rows: %w", err)
    }

    return &model.StatisticsResponse{
        Total:    total,
        ByType:   byType,
        ByStatus: byStatus,
    }, nil
}

func (p *PostgresStorage) GetCount(assetType, status string) (int, error) {

    query := "SELECT COUNT(*) FROM assets WHERE 1=1"
    var args []interface{}
    argIndex := 1

    if assetType != "" {
        query += fmt.Sprintf(" AND type = $%d", argIndex)
        args = append(args, assetType)
        argIndex++
    }

    if status != "" {
        query += fmt.Sprintf(" AND status = $%d", argIndex)
        args = append(args, status)
        argIndex++
    }

    var count int
    err := p.db.QueryRow(query, args...).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count assets: %w", err)
    }

    return count, nil
}

func (p *PostgresStorage) CreateBatch(assets []*model.Asset) ([]string, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
		INSERT INTO assets (id, name, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	var ids []string
	for _, asset := range assets {
		_, err := tx.Exec(
			query,
			asset.ID,
			asset.Name,
			asset.Type,
			asset.Status,
			asset.CreatedAt,
			asset.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create asset: %w", err)
		}
		ids = append(ids, asset.ID)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ids, nil
}

func (p *PostgresStorage) DeleteBatch(ids []string) (int, int, error) {
	if len(ids) == 0 {
        return 0, 0, nil
    }

    query := "DELETE FROM assets WHERE id = $1"

    var deletedCount int
    var notFoundCount int

    for _, id := range ids {

        result, err := p.db.Exec(query, id)
        if err != nil {
            return 0, 0, fmt.Errorf("failed to delete asset %s: %w", id, err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            return 0, 0, fmt.Errorf("failed to get rows affected: %w", err)
        }

        if rowsAffected == 0 {
            notFoundCount++
        } else {
            deletedCount++
        }
    }

    return deletedCount, notFoundCount, nil
}
/*
🎓 TEACHING NOTES:

1. Interface Implementation:
   - PostgresStorage implements Storage interface
   - Same methods as MemoryStorage
   - Main.go CHỈ CẦN THAY 1 DÒNG!

2. SQL Queries:
   - Parameterized queries ($1, $2) prevent SQL injection
   - NEVER concatenate user input into SQL!
   - ❌ BAD: query := "SELECT * FROM assets WHERE id = '" + id + "'"
   - ✅ GOOD: query := "SELECT * FROM assets WHERE id = $1"

3. Error Handling:
   - Check sql.ErrNoRows → return model.ErrNotFound
   - Wrap errors with fmt.Errorf("context: %w", err)
   - RowsAffected() to verify UPDATE/DELETE success

4. Scanning Rows:
   - rows.Scan() maps columns to struct fields
   - ORDER MATTERS! Must match SELECT order
   - Don't forget defer rows.Close()

5. Dynamic Query Building:
   - Filter() builds query dynamically based on parameters
   - Track parameter count for $1, $2, $3...
   - WHERE 1=1 trick for easier dynamic AND conditions

6. ILIKE vs LIKE:
   - LIKE: case-sensitive
   - ILIKE: case-insensitive (PostgreSQL specific)
   - % wildcards for partial matching

7. Connection Pool:
   - sql.DB maintains connection pool automatically
   - db.Exec(), db.Query(), db.QueryRow() reuse connections
   - Don't close db in storage methods!

8. Transaction Support (Buổi 4):
   - Current: each operation is auto-committed
   - Future: db.Begin() for multi-step operations

COMPARISON: Memory vs Postgres

MemoryStorage:
- data := make(map[string]*model.Asset)
- Fast: O(1) lookups

PostgresStorage:
- data in database on disk
- Slightly slower but persistent
- Can handle millions of records
- Support for advanced queries (JOIN, aggregation)

KEY POINT: Clean Architecture cho phép swap giữa 2 implementations này
mà KHÔNG THAY ĐỔI business logic!
*/
