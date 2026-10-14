package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"uuid"

	"storeapi/src/internal/domain"
)

type ClientRepository struct{ store *Store }

func NewClientRepository(store *Store) *ClientRepository { return &ClientRepository{store: store} }

func scanClient(scanner interface{ Scan(...any) error }) (domain.Client, error) {
	var d ClientDAO
	if err := scanner.Scan(&d.ID, &d.Name, &d.Surname, &d.Birthday, &d.Gender, &d.RegistrationAt, &d.Address.ID, &d.Address.Country, &d.Address.City, &d.Address.Street); err != nil {
		return domain.Client{}, err
	}
	return clientToDomain(d), nil
}

func (r *ClientRepository) Create(ctx context.Context, client domain.Client) (domain.Client, error) {
	tx, err := r.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.Client{}, err
	}
	defer tx.Rollback()
	addressID := uuid.New()
	id := uuid.New()
	if _, err = tx.ExecContext(ctx, `INSERT INTO addresses (id,country,city,street) VALUES ($1,$2,$3,$4)`, addressID, client.Address.Country, client.Address.City, client.Address.Street); err != nil {
		return domain.Client{}, fmt.Errorf("insert address: %w", err)
	}
	reg := client.RegistrationAt
	if reg.IsZero() {
		reg = time.Now().UTC()
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO clients (id,client_name,client_surname,birthday,gender,registration_date,address_id) VALUES ($1,$2,$3,$4,$5,$6,$7)`, id, client.Name, client.Surname, client.Birthday, client.Gender, reg, addressID); err != nil {
		return domain.Client{}, fmt.Errorf("insert client: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return domain.Client{}, err
	}
	client.ID, client.RegistrationAt, client.Address.ID = uuid.UUID(id), reg, addressID
	return client, nil
}

func (r *ClientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.store.DB.ExecContext(ctx, `DELETE FROM clients WHERE id = $1`, id.String())
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ClientRepository) FindByName(ctx context.Context, firstName, lastName string) ([]domain.Client, error) {
	rows, err := r.store.DB.QueryContext(ctx, `SELECT `+clientColumns+` FROM clients c JOIN addresses a ON a.id=c.address_id WHERE c.client_name=$1 AND c.client_surname=$2 ORDER BY c.client_surname,c.client_name,c.id`, firstName, lastName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readClients(rows)
}

func (r *ClientRepository) List(ctx context.Context, limit, offset *int) ([]domain.Client, error) {
	q := `SELECT ` + clientColumns + ` FROM clients c JOIN addresses a ON a.id=c.address_id ORDER BY c.client_surname,c.client_name,c.id`
	args := []any{}
	if limit != nil {
		q += ` LIMIT $1`
		args = append(args, *limit)
	}
	if offset != nil {
		if limit == nil {
			q += ` OFFSET $1`
		} else {
			q += ` OFFSET $2`
		}
		args = append(args, *offset)
	}
	rows, err := r.store.DB.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return readClients(rows)
}

func readClients(rows *sql.Rows) ([]domain.Client, error) {
	clients := make([]domain.Client, 0)
	for rows.Next() {
		c, err := scanClient(rows)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

func (r *ClientRepository) ChangeAddress(ctx context.Context, id uuid.UUID, address domain.Address) (domain.Client, error) {
	tx, err := r.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return domain.Client{}, err
	}
	defer tx.Rollback()
	var oldAddressID string
	if err = tx.QueryRowContext(ctx, `SELECT address_id FROM clients WHERE id=$1 FOR UPDATE`, id.String()).Scan(&oldAddressID); err != nil {
		if err == sql.ErrNoRows {
			return domain.Client{}, domain.ErrNotFound
		}
		return domain.Client{}, err
	}
	newAddressID := uuid.New()
	if _, err = tx.ExecContext(ctx, `INSERT INTO addresses (id,country,city,street) VALUES ($1,$2,$3,$4)`, newAddressID, address.Country, address.City, address.Street); err != nil {
		return domain.Client{}, err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE clients SET address_id=$1 WHERE id=$2`, newAddressID, id.String()); err != nil {
		return domain.Client{}, err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM addresses WHERE id=$1`, oldAddressID); err != nil {
		return domain.Client{}, err
	}
	if err = tx.Commit(); err != nil {
		return domain.Client{}, err
	}
	client, err := r.getByID(ctx, id)
	if err != nil {
		return domain.Client{}, err
	}
	return client, nil
}

func (r *ClientRepository) getByID(ctx context.Context, id uuid.UUID) (domain.Client, error) {
	row := r.store.DB.QueryRowContext(ctx, `SELECT `+clientColumns+` FROM clients c JOIN addresses a ON a.id=c.address_id WHERE c.id=$1`, id.String())
	c, err := scanClient(row)
	if err == sql.ErrNoRows {
		return domain.Client{}, domain.ErrNotFound
	}
	return c, err
}

const clientColumns = `c.id, c.client_name, c.client_surname, c.birthday, c.gender, c.registration_date, a.id, a.country, a.city, a.street`
