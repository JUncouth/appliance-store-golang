package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"uuid"

	"github.com/lib/pq"

	"storeapi/src/internal/domain"
)

type ProductRepository struct{ store *Store }

func NewProductRepository(store *Store) *ProductRepository { return &ProductRepository{store: store} }

func scanProduct(scanner interface{ Scan(...any) error }) (domain.Product, error) {
	var d ProductDAO
	var imageID sql.NullString
	if err := scanner.Scan(&d.ID, &d.Name, &d.Category, &d.Price, &d.AvailableStock, &d.LastUpdateAt, &d.SupplierID, &imageID); err != nil {
		return domain.Product{}, err
	}
	if imageID.Valid {
		imageIDString := imageID.String
		imageIDVal, _ := uuid.Parse(imageIDString)
		d.ImageID = &imageIDVal
	}
	return productToDomain(d), nil
}

func (r *ProductRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	id := uuid.New()
	if product.ImageID != nil {
		var n int
		if err := r.store.DB.QueryRowContext(ctx, `SELECT COUNT(1) FROM images WHERE id=$1`, product.ImageID.String()).Scan(&n); err != nil {
			return domain.Product{}, err
		}
		if n == 0 {
			return domain.Product{}, domain.ErrNotFound
		}
	}
	_, err := r.store.DB.ExecContext(ctx, `INSERT INTO products (id,name,category,price,available_stock,last_update_date,supplier_id,image_id) VALUES ($1,$2,$3,$4,$5,NOW(),$6,$7)`, id, product.Name, product.Category, product.Price, product.AvailableStock, product.SupplierID, product.ImageID)
	if err != nil {
		return domain.Product{}, mapPGError(err)
	}
	product.ID = id
	product.LastUpdateAt = nowUTC()
	return product, nil
}

func (r *ProductRepository) AddProductPhoto(ctx context.Context, id uuid.UUID, imageID uuid.UUID) (domain.Product, error) {
	row := r.store.DB.QueryRowContext(ctx, `UPDATE products p SET image_id=$1, last_update_date=NOW() WHERE id=$2 RETURNING `+productColumns, imageID, id)
	p, err := scanProduct(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Product{}, domain.ErrNotFound
		}
		return domain.Product{}, err
	}
	return p, err
}

func (r *ProductRepository) DecreaseStock(ctx context.Context, id uuid.UUID, quantity int) (domain.Product, error) {
	row := r.store.DB.QueryRowContext(ctx, `UPDATE products p SET available_stock=available_stock-$1,last_update_date=NOW() WHERE id=$2 AND available_stock >= $1 RETURNING `+productColumns, quantity, id.String())
	p, err := scanProduct(row)
	if err == sql.ErrNoRows {
		var stock int
		err2 := r.store.DB.QueryRowContext(ctx, `SELECT available_stock FROM products WHERE id=$1`, id.String()).Scan(&stock)
		if err2 == sql.ErrNoRows {
			return domain.Product{}, domain.ErrNotFound
		}
		if err2 != nil {
			return domain.Product{}, err2
		}
		return domain.Product{}, domain.ErrInsufficientStock
	}
	return p, err
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	row := r.store.DB.QueryRowContext(ctx, `SELECT `+productColumns+` FROM products p WHERE p.id=$1`, id.String())
	p, err := scanProduct(row)
	if err == sql.ErrNoRows {
		return domain.Product{}, domain.ErrNotFound
	}
	return p, err
}

func (r *ProductRepository) ListAvailable(ctx context.Context) ([]domain.Product, error) {
	rows, err := r.store.DB.QueryContext(ctx, `SELECT `+productColumns+` FROM products p WHERE p.available_stock > 0 ORDER BY p.name,p.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Product, 0)
	for rows.Next() {
		p, e := scanProduct(rows)
		if e != nil {
			return nil, e
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.store.DB.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id.String())
	if err != nil {
		return mapPGError(err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func nullableUUID(id *uuid.UUID) any {
	if id == nil {
		return nil
	}
	return id.String()
}

func nowUTC() (t time.Time) { return time.Now().UTC() }

func mapPGError(err error) error {
	if err == nil {
		return nil
	}
	if pqe, ok := err.(*pq.Error); ok && pqe.Code == "23503" {
		return domain.ErrNotFound
	}
	if pqe, ok := err.(*pq.Error); ok && pqe.Code == "23001" {
		return domain.ErrConflict
	}
	return fmt.Errorf("database error: %w", err)
}

const productColumns = `p.id,p.name,p.category,p.price,p.available_stock,p.last_update_date,p.supplier_id,p.image_id`
