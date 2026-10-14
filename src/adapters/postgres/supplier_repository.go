package postgres

import (
	"context"
	"database/sql"
	"uuid"

	"storeapi/src/internal/domain"
)

type SupplierRepository struct{ store *Store }

func NewSupplierRepository(store *Store) *SupplierRepository {
	return &SupplierRepository{store: store}
}

func scanSupplier(scanner interface{ Scan(...any) error }) (domain.Supplier, error) {
	var d SupplierDAO
	if err := scanner.Scan(&d.ID, &d.Name, &d.PhoneNumber, &d.Address.ID, &d.Address.Country, &d.Address.City, &d.Address.Street); err != nil {
		return domain.Supplier{}, err
	}
	return supplierToDomain(d), nil
}

func (r *SupplierRepository) Create(ctx context.Context, supplier domain.Supplier) (domain.Supplier, error) {
	tx, err := r.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.Supplier{}, err
	}
	defer tx.Rollback()
	id := uuid.New()
	addressID := uuid.New()
	if _, err = tx.ExecContext(ctx, `INSERT INTO addresses (id,country,city,street) VALUES ($1,$2,$3,$4)`, addressID, supplier.Address.Country, supplier.Address.City, supplier.Address.Street); err != nil {
		return domain.Supplier{}, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO suppliers (id,name,address_id,phone_number) VALUES ($1,$2,$3,$4)`, id, supplier.Name, addressID, supplier.PhoneNumber); err != nil {
		return domain.Supplier{}, err
	}
	if err = tx.Commit(); err != nil {
		return domain.Supplier{}, err
	}
	supplier.ID = id
	supplier.Address.ID = addressID
	return supplier, nil
}

func (r *SupplierRepository) ChangeAddress(ctx context.Context, id uuid.UUID, address domain.Address) (domain.Supplier, error) {
	tx, err := r.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.Supplier{}, err
	}
	defer tx.Rollback()
	var old string
	if err = tx.QueryRowContext(ctx, `SELECT address_id FROM suppliers WHERE id=$1 FOR UPDATE`, id.String()).Scan(&old); err != nil {
		if err == sql.ErrNoRows {
			return domain.Supplier{}, domain.ErrNotFound
		}
		return domain.Supplier{}, err
	}
	newAddress := uuid.New()
	if _, err = tx.ExecContext(ctx, `INSERT INTO addresses (id,country,city,street) VALUES ($1,$2,$3,$4)`, newAddress, address.Country, address.City, address.Street); err != nil {
		return domain.Supplier{}, err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE suppliers SET address_id=$1 WHERE id=$2`, newAddress, id.String()); err != nil {
		return domain.Supplier{}, err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM addresses WHERE id=$1`, old); err != nil {
		return domain.Supplier{}, err
	}
	if err = tx.Commit(); err != nil {
		return domain.Supplier{}, err
	}
	return r.GetByID(ctx, id)
}

func (r *SupplierRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var addressID string
	if err = tx.QueryRowContext(ctx, `SELECT address_id FROM suppliers WHERE id=$1 FOR UPDATE`, id.String()).Scan(&addressID); err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrNotFound
		}
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM suppliers WHERE id=$1`, id.String()); err != nil {
		return mapPGError(err)
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM addresses WHERE id=$1`, addressID); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *SupplierRepository) List(ctx context.Context) ([]domain.Supplier, error) {
	rows, err := r.store.DB.QueryContext(ctx, `SELECT `+supplierColumns+` FROM suppliers s JOIN addresses a ON a.id=s.address_id ORDER BY s.name,s.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Supplier, 0)
	for rows.Next() {
		v, e := scanSupplier(rows)
		if e != nil {
			return nil, e
		}
		items = append(items, v)
	}
	return items, rows.Err()
}

func (r *SupplierRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Supplier, error) {
	row := r.store.DB.QueryRowContext(ctx, `SELECT `+supplierColumns+` FROM suppliers s JOIN addresses a ON a.id=s.address_id WHERE s.id=$1`, id.String())
	v, err := scanSupplier(row)
	if err == sql.ErrNoRows {
		return domain.Supplier{}, domain.ErrNotFound
	}
	return v, err
}

const supplierColumns = `s.id,s.name,s.phone_number,a.id,a.country,a.city,a.street`
