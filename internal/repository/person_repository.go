package repository

import (
	"context"
	"fmt"
	"errors"
	"reflect"
	"strings"

	"EffectiveMobile/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"EffectiveMobile/config"
)

// небольшой костыль
var ErrNotFound = errors.New("not found")

type PersonRepository struct {
	db *pgxpool.Pool
	cfg *config.Config
}

func NewPersonRepository(db *pgxpool.Pool, cfg *config.Config) *PersonRepository {
	return &PersonRepository{
		db: db,
		cfg: cfg,
	}
}

func (pr *PersonRepository) Insert(ctx context.Context, person models.Person) error {
	query := fmt.Sprintf(`INSERT INTO %s.person (id, name, surname, patronymic, age, gender, nationality) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)`, pr.cfg.DB.Schema)

	_, err := pr.db.Exec(ctx, query, person.ID, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
	if err != nil {
		return fmt.Errorf("PersonRepository.Insert: failed to insert person: %w", err)
	}
	return nil
}

func (pr *PersonRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s.person WHERE id = $1`, pr.cfg.DB.Schema)

	res, err := pr.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("PersonRepository.Delete: failed to delete person: %w", err)
	}

	rowsAffected:= res.RowsAffected()

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (pr *PersonRepository) Update(ctx context.Context, id string, person models.Person) error {
	query := fmt.Sprintf(`UPDATE %s.person 
        SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6 
        WHERE id = $7`, pr.cfg.DB.Schema)

	res, err := pr.db.Exec(ctx, query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, id)
	if err != nil {
		return fmt.Errorf("PersonRepository.Update: failed to update person: %w", err)
	}

	rowsAffected:= res.RowsAffected()

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (pr *PersonRepository) Get(ctx context.Context, filters models.PersonFilters, limit, offset int) (*[]models.Person, error) {
	query := `SELECT id, name, surname, patronymic, age, gender, nationality FROM person WHERE 1=1`
	args := []interface{}{}
	argID := 1

	val := reflect.ValueOf(filters)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name
		if field.IsZero() {
			continue 
		}

		query += fmt.Sprintf(" AND %s = $%d", strings.ToLower(fieldName), argID)
		args = append(args, field.Interface())
		argID++
	}

	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := pr.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PersonRepository.Get: failed to get persons: %w", err)
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality); err != nil {
			return nil, fmt.Errorf("PersonRepository.Get: failed to scan person: %w", err)
		}
		persons = append(persons, p)
	}
	return &persons, nil
}