package postgres

import (
	"context"
	"database/sql"
	"uuid"

	"storeapi/src/internal/domain"
)

type ImageRepository struct{ store *Store }

func NewImageRepository(store *Store) *ImageRepository { return &ImageRepository{store: store} }

func scanImage(scanner interface{ Scan(...any) error }) (domain.Image, error) {
	var d ImageDAO
	if err := scanner.Scan(&d.ID, &d.Data); err != nil {
		return domain.Image{}, err
	}
	return imageToDomain(d), nil
}

func (r *ImageRepository) Create(ctx context.Context, image domain.Image) (domain.Image, error) {
	tx, err := r.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.Image{}, err
	}
	defer tx.Rollback()
	imageID := uuid.New()
	if _, err = tx.ExecContext(ctx, `INSERT INTO images (id,image) VALUES ($1,$2)`, imageID, image.Data); err != nil {
		return domain.Image{}, mapPGError(err)
	}
	if err = tx.Commit(); err != nil {
		return domain.Image{}, err
	}
	image.ID = imageID
	return image, nil
}

func (r *ImageRepository) Replace(ctx context.Context, id uuid.UUID, data []byte) (domain.Image, error) {
	res, err := r.store.DB.ExecContext(ctx, `UPDATE images SET image=$1 WHERE id=$2`, data, id.String())
	if err != nil {
		return domain.Image{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.Image{}, domain.ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *ImageRepository) ListAvailable(ctx context.Context) ([]domain.Image, error) {
	rows, err := r.store.DB.QueryContext(ctx, `SELECT * FROM images`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Image, 0)
	for rows.Next() {
		im, e := scanImage(rows)
		if e != nil {
			return nil, e
		}
		items = append(items, im)
	}
	return items, rows.Err()
}

func (r *ImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var productID string
	if err = tx.QueryRowContext(ctx, `SELECT product_id FROM images WHERE id=$1 FOR UPDATE`, id.String()).Scan(&productID); err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrNotFound
		}
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE products SET image_id=NULL WHERE id=$1`, productID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM images WHERE id=$1`, id.String()); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ImageRepository) GetByProductID(ctx context.Context, id uuid.UUID) (domain.Image, error) {
	tx, err := r.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.Image{}, err
	}
	defer tx.Rollback()
	var image domain.Image
	if err = tx.QueryRowContext(ctx, `SELECT image_id FROM products WHERE id=$1`, id).Scan(&image.ID); err != nil {
		if err == sql.ErrNoRows {
			return domain.Image{}, domain.ErrNotFound
		}
		return domain.Image{}, err
	}
	if err = tx.QueryRowContext(ctx, `SELECT image FROM images WHERE id=$1`, image.ID).Scan(&image.Data); err != nil {
		if err == sql.ErrNoRows {
			return domain.Image{}, domain.ErrNotFound
		}
		return domain.Image{}, err
	}
	return image, err
}

func (r *ImageRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Image, error) {
	row := r.store.DB.QueryRowContext(ctx, `SELECT id,image FROM images WHERE id=$1`, id)
	v, err := scanImage(row)
	if err == sql.ErrNoRows {
		return domain.Image{}, domain.ErrNotFound
	}
	return v, err
}
