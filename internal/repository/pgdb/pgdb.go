package pgdb

import "github.com/umisto/cities-svc/internal/domain/entity"

func (a Agglomeration) ToEntity() entity.Agglomeration {
	ent := entity.Agglomeration{
		ID:        a.ID,
		Status:    string(a.Status),
		Name:      a.Name,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
	if a.Icon.Valid {
		ent.Icon = a.Icon.String
	}

	return ent
}

func (c City) ToEntity() entity.City {
	ent := entity.City{
		ID:        c.ID,
		Status:    string(c.Status),
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.AgglomerationID.Valid {
		id := c.AgglomerationID.UUID
		ent.AgglomerationID = &id
	}
	if c.Slug.Valid {
		slug := c.Slug.String
		ent.Slug = &slug
	}
	if c.Icon.Valid {
		icon := c.Icon.String
		ent.Icon = &icon
	}
	if c.Banner.Valid {
		banner := c.Banner.String
		ent.Banner = &banner
	}

	return ent
}
