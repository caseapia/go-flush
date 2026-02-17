package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
)

func (r *Repository) SearchAllRanks(ctx context.Context) ([]models.RankStructure, error) {
	var ranks []models.RankStructure
	err := r.db.NewSelect().
		Model(&ranks).
		Relation("Users").
		Relation("Developers").
		Limit(COLUMNS_LIMIT).
		Scan(ctx)
	return ranks, err
}

func (r *Repository) SearchRankByID(ctx context.Context, id int) (*models.RankStructure, error) {
	rank := new(models.RankStructure)
	err := r.db.NewSelect().
		Model(rank).
		Relation("Users").
		Relation("Developers").
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return rank, nil
}

func (r *Repository) SearchRankByName(ctx context.Context, rankName string) (*models.RankStructure, error) {
	rank := new(models.RankStructure)
	err := r.db.NewSelect().
		Model(rank).
		Relation("Users").
		Relation("Developers").
		Where("name = ?", rankName).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return rank, nil
}

func (r *Repository) CreateRank(ctx context.Context, rank *models.RankStructure) error {
	_, err := r.db.NewInsert().
		Model(rank).
		Exec(ctx)

	return err
}

func (r *Repository) DeleteRank(ctx context.Context, rank *models.RankStructure) error {
	_, err := r.db.NewDelete().
		Model(rank).
		WherePK().
		Exec(ctx)
	return err
}

func (r *Repository) EditRank(ctx context.Context, rank *models.RankStructure) (*models.RankStructure, error) {
	_, err := r.db.NewUpdate().
		Model(rank).
		Column("name", "color", "flags").
		WherePK().
		Exec(ctx)

	return rank, err
}
