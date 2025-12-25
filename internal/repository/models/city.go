package models

import (
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func GetCityByID(c pgdb.GetCityByIDRow) entity.City {
	ent := entity.City{
		ID:        c.ID,
		Status:    string(c.Status),
		Name:      c.Name,
		Point:     orb.Point{c.PointLng, c.PointLat},
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.AgglomerationID.Valid {
		id := c.AgglomerationID.UUID
		ent.AgglomerationID = &id
	}
	if c.Slug.Valid {
		s := c.Slug.String
		ent.Slug = &s
	}
	if c.Icon.Valid {
		s := c.Icon.String
		ent.Icon = &s
	}
	if c.Banner.Valid {
		s := c.Banner.String
		ent.Banner = &s
	}

	return ent
}

func GetCityBySlug(c pgdb.GetCityBySlugRow) entity.City {
	ent := entity.City{
		ID:        c.ID,
		Status:    string(c.Status),
		Name:      c.Name,
		Point:     orb.Point{c.PointLng, c.PointLat},
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.AgglomerationID.Valid {
		id := c.AgglomerationID.UUID
		ent.AgglomerationID = &id
	}
	if c.Slug.Valid {
		s := c.Slug.String
		ent.Slug = &s
	}
	if c.Icon.Valid {
		s := c.Icon.String
		ent.Icon = &s
	}
	if c.Banner.Valid {
		s := c.Banner.String
		ent.Banner = &s
	}

	return ent
}

func FilterCitiesNearestRow(rows []pgdb.FilterCitiesNearestRow) map[int64]entity.City {
	result := make(map[int64]entity.City)
	for _, c := range rows {
		ent := entity.City{
			ID:        c.ID,
			Status:    string(c.Status),
			Name:      c.Name,
			Point:     orb.Point{c.PointLng, c.PointLat},
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}

		if c.AgglomerationID.Valid {
			id := c.AgglomerationID.UUID
			ent.AgglomerationID = &id
		}
		if c.Slug.Valid {
			s := c.Slug.String
			ent.Slug = &s
		}
		if c.Icon.Valid {
			s := c.Icon.String
			ent.Icon = &s
		}
		if c.Banner.Valid {
			s := c.Banner.String
			ent.Banner = &s
		}

		result[c.DistanceM] = ent
	}

	return result
}

func FilterCitiesRow(rows []pgdb.FilterCitiesRow) []entity.City {
	result := make([]entity.City, 0, len(rows))

	for _, c := range rows {
		ent := entity.City{
			ID:        c.ID,
			Status:    string(c.Status),
			Name:      c.Name,
			Point:     orb.Point{c.PointLng, c.PointLat},
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}

		if c.AgglomerationID.Valid {
			id := c.AgglomerationID.UUID
			ent.AgglomerationID = &id
		}
		if c.Slug.Valid {
			s := c.Slug.String
			ent.Slug = &s
		}
		if c.Icon.Valid {
			s := c.Icon.String
			ent.Icon = &s
		}
		if c.Banner.Valid {
			s := c.Banner.String
			ent.Banner = &s
		}

		result = append(result, ent)
	}

	return result
}
