package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lampadovnikita/StorekeeperTask/pkg/data"
)

type PGStorage struct {
	pool *pgxpool.Pool
}

func NewPGStorage(pool *pgxpool.Pool) *PGStorage {
	return &PGStorage{pool: pool}
}

func (s *PGStorage) GetGatheringInfo(orderIDs []int) (d []data.GatheringInfo, err error) {
	q := `SELECT r.name as rack_name, p.name as product_name, p.id as product_id, o.id as order_id, o.amount,
			  (SELECT array_agg(rks.name)
	 			  FROM storage
	 				  JOIN racks as rks
	 					  ON storage.rack_id = rks.id
	 			  WHERE product_id = s.product_id AND
	 		  		    is_rack_primary = false
		      ) as additional_racks
		      FROM orders as o 
			  	  JOIN products as p
				  	  ON o.product_id = p.id
			  	  JOIN storage as s
				  	  ON s.product_id = p.id
			  	  JOIN racks as r
				  	  ON s.rack_id = r.id
		  	  WHERE o.id in (10, 11, 14, 15) AND
		  	  	  is_rack_primary = true
		  	  ORDER BY r.id`

	rows, err := s.pool.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		gathInfo := data.GatheringInfo{}

		err := rows.Scan(&gathInfo.RackName, &gathInfo.ProductName, &gathInfo.ProductID,
			&gathInfo.OrderID, &gathInfo.Amount, &gathInfo.AdditionalRacks)
		if err != nil {
			return nil, err
		}

		d = append(d, gathInfo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return d, nil
}
