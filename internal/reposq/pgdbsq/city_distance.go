package pgdbsq

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type CityDistance struct {
	City
	DistanceMeters float64 `json:"distance_meters"`
}

func (cd *CityDistance) scan(row sq.RowScanner) error {
	err := row.Scan(
		&cd.ID,
		&cd.AgglomerationID,
		&cd.Status,
		&cd.Slug,
		&cd.Name,
		&cd.Icon,
		&cd.Banner,
		&cd.Point,
		&cd.CreatedAt,
		&cd.UpdatedAt,
		&cd.DistanceMeters,
	)
	if err != nil {
		return fmt.Errorf("scanning city distance: %w", err)
	}
	return nil
}

func (q CityQ) OrderNearest(limit uint, lat, lng float64) CityQ {
	q.selector = q.selector.
		Column(sq.Expr(
			"ST_Distance(c.point, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography) AS distance_meters",
			lng, lat,
		)).
		OrderBy("distance_meters ASC").
		Limit(uint64(limit))

	return q
}

func (q CityQ) FilterWithinRadiusMeters(lat, lng float64, radiusMeters float64) CityQ {
	q.selector = q.selector.Where(sq.Expr(
		"ST_DWithin(c.point, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)",
		lng, lat, radiusMeters,
	))
	q.counter = q.counter.Where(sq.Expr(
		"ST_DWithin(c.point, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)",
		lng, lat, radiusMeters,
	))
	return q
}

func (q CityQ) SelectNearest(ctx context.Context) ([]CityDistance, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", CityTable, err)
	}

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing select query for %s: %w", CityTable, err)
	}
	defer rows.Close()

	var out []CityDistance
	for rows.Next() {
		var cd CityDistance
		if err = cd.scan(rows); err != nil {
			return nil, err
		}
		out = append(out, cd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q CityQ) GetNearest(ctx context.Context) (CityDistance, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return CityDistance{}, fmt.Errorf("building select query for %s: %w", CityTable, err)
	}

	var cd CityDistance
	err = cd.scan(q.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		return CityDistance{}, err
	}

	return cd, nil
}
