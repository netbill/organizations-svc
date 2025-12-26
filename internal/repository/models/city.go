package models

import (
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

//func GetCityByID(c pgdbsql.GetCityByIDRow) entity.City {
//	ent := entity.City{
//		ID:        c.ID,
//		Status:    string(c.Status),
//		Name:      c.Name,
//		Point:     orb.Point{c.PointLng, c.PointLat},
//		CreatedAt: c.CreatedAt,
//		UpdatedAt: c.UpdatedAt,
//	}
//
//	if c.AgglomerationID.Valid {
//		id := c.AgglomerationID.UUID
//		ent.AgglomerationID = &id
//	}
//	if c.Slug.Valid {
//		s := c.Slug.String
//		ent.Slug = &s
//	}
//	if c.Icon.Valid {
//		s := c.Icon.String
//		ent.Icon = &s
//	}
//	if c.Banner.Valid {
//		s := c.Banner.String
//		ent.Banner = &s
//	}
//
//	return ent
//}
//
//func GetCityBySlug(c pgdbsql.GetCityBySlugRow) entity.City {
//	ent := entity.City{
//		ID:        c.ID,
//		Status:    string(c.Status),
//		Name:      c.Name,
//		Point:     orb.Point{c.PointLng, c.PointLat},
//		CreatedAt: c.CreatedAt,
//		UpdatedAt: c.UpdatedAt,
//	}
//
//	if c.AgglomerationID.Valid {
//		id := c.AgglomerationID.UUID
//		ent.AgglomerationID = &id
//	}
//	if c.Slug.Valid {
//		s := c.Slug.String
//		ent.Slug = &s
//	}
//	if c.Icon.Valid {
//		s := c.Icon.String
//		ent.Icon = &s
//	}
//	if c.Banner.Valid {
//		s := c.Banner.String
//		ent.Banner = &s
//	}
//
//	return ent
//}
//
//func FilterCitiesNearestRow(rows []pgdbsql.FilterCitiesNearestRow) map[int64]entity.City {
//	result := make(map[int64]entity.City)
//	for _, c := range rows {
//		ent := entity.City{
//			ID:        c.ID,
//			Status:    string(c.Status),
//			Name:      c.Name,
//			Point:     orb.Point{c.PointLng, c.PointLat},
//			CreatedAt: c.CreatedAt,
//			UpdatedAt: c.UpdatedAt,
//		}
//
//		if c.AgglomerationID.Valid {
//			id := c.AgglomerationID.UUID
//			ent.AgglomerationID = &id
//		}
//		if c.Slug.Valid {
//			s := c.Slug.String
//			ent.Slug = &s
//		}
//		if c.Icon.Valid {
//			s := c.Icon.String
//			ent.Icon = &s
//		}
//		if c.Banner.Valid {
//			s := c.Banner.String
//			ent.Banner = &s
//		}
//
//		result[c.DistanceM] = ent
//	}
//
//	return result
//}
//
//func FilterCitiesRow(rows []pgdbsql.FilterCitiesRow) []entity.City {
//	result := make([]entity.City, 0, len(rows))
//
//	for _, c := range rows {
//		ent := entity.City{
//			ID:        c.ID,
//			Status:    string(c.Status),
//			Name:      c.Name,
//			Point:     orb.Point{c.PointLng, c.PointLat},
//			CreatedAt: c.CreatedAt,
//			UpdatedAt: c.UpdatedAt,
//		}
//
//		if c.AgglomerationID.Valid {
//			id := c.AgglomerationID.UUID
//			ent.AgglomerationID = &id
//		}
//		if c.Slug.Valid {
//			s := c.Slug.String
//			ent.Slug = &s
//		}
//		if c.Icon.Valid {
//			s := c.Icon.String
//			ent.Icon = &s
//		}
//		if c.Banner.Valid {
//			s := c.Banner.String
//			ent.Banner = &s
//		}
//
//		result = append(result, ent)
//	}
//
//	return result
//}

func City(c pgdb.City) entity.City {
	ent := entity.City{
		ID:        c.ID,
		Status:    c.Status,
		Name:      c.Name,
		Point:     c.Point,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.AgglomerationID != nil {
		ent.AgglomerationID = c.AgglomerationID
	}
	if c.Slug != nil {
		ent.Slug = c.Slug
	}
	if c.Icon != nil {
		ent.Icon = c.Icon
	}
	if c.Banner != nil {
		ent.Banner = c.Banner
	}

	return ent
}

func CityDistance(cd pgdb.CityDistance) entity.City {
	ent := entity.City{
		ID:        cd.ID,
		Status:    cd.Status,
		Name:      cd.Name,
		Point:     cd.Point,
		CreatedAt: cd.CreatedAt,
		UpdatedAt: cd.UpdatedAt,
	}

	if cd.AgglomerationID != nil {
		ent.AgglomerationID = cd.AgglomerationID
	}
	if cd.Slug != nil {
		ent.Slug = cd.Slug
	}
	if cd.Icon != nil {
		ent.Icon = cd.Icon
	}
	if cd.Banner != nil {
		ent.Banner = cd.Banner
	}

	return ent
}
