package service

import (
	"context"
	"fmt"

	"EffectiveMobile/internal/models"
	"EffectiveMobile/internal/transport/http/enrichment"

	"github.com/google/uuid"
)

type IPersonRepository interface {
	Get(context.Context, models.PersonFilters, int, int) (*[]models.Person, error)
	Delete(context.Context,string) error
	Insert(context.Context, models.Person) error
	Update(context.Context, string, models.Person) error
}

type PersonService struct {
	PR IPersonRepository
}

func NewPersonService(pr IPersonRepository) *PersonService {
	return &PersonService{
		PR: pr,
	}
}

func (ps *PersonService) Get(ctx context.Context, filters models.PersonFilters, limit, offset int) (*[]models.Person, error) {
	persons, err := ps.PR.Get(ctx, filters, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("PersonService.Get: %w", err)
	}

	return persons, nil
}

func (ps *PersonService) Delete(ctx context.Context, id string) error {
	err := ps.PR.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("PersonService.Delete: %w", err)
	}

	return nil
}

func (ps *PersonService) Insert(ctx context.Context, person models.Person) (*models.Person, error) {
	// Обогащение данных
	age, err := enrichment.GetAge(person.Name)
	if err != nil {
		return nil, fmt.Errorf("PersonService.Insert: failed to fetch age: %w", err)
	}

	gender, err := enrichment.GetGender(person.Name)
	if err != nil {
		return nil, fmt.Errorf("PersonService.Insert: failed to fetch gender: %w", err)
	}

	nationality, err := enrichment.GetNationality(person.Name)
	if err != nil {
		return nil, fmt.Errorf("PersonService.Insert: failed to fetch nationality: %w", err)
	}

	person.Age = age
	person.Gender = gender
	person.Nationality = nationality

	id := uuid.New().String()

	person.ID = id

	err = ps.PR.Insert(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("PersonService: %w", err)
	}

	return &person, nil
}

func (ps *PersonService) Update(ctx context.Context, person models.Person, id string) error {
	err := ps.PR.Update(ctx, id, person)
	if err != nil {
		return fmt.Errorf("PersonService.Update: %w", err)
	}

	return nil
}