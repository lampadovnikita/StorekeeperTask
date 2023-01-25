# StorekeeperTask
## Примерная структура таблиц
<img src="/repository_resources/tables.png"  width="610" height="400">  

## Примеры работы программы
<img src="/repository_resources/example1.png"  width="410" height="400">  

<img src="/repository_resources/example2.png"  width="560" height="400">  

## Запрос на получение необходимых в задаче данных
```sql
SELECT r.name AS rack_name, p.name AS product_name, p.id AS product_id, o.id AS order_id, o.amount,
	(SELECT array_agg(rks.name)
	 	FROM storage
	 		JOIN racks AS rks
	 			ON storage.rack_id = rks.id
	 	WHERE product_id = s.product_id AND
	 		  is_rack_primary = false
	) AS additional_racks
	FROM orders AS o 
		JOIN products AS p
			ON o.product_id = p.id
		JOIN storage AS s
			ON s.product_id = p.id
		JOIN racks AS r
			ON s.rack_id = r.id
	WHERE o.id = ANY ('{10, 11, 14, 15}') AND
		  is_rack_primary = true
	ORDER BY r.id
```

## Запросы на создание таблиц
```sql
CREATE TABLE IF NOT EXISTS products (
	id INT PRIMARY KEY,
	name TEXT NOT NULL
	
	CONSTRAINT name_positive_size CHECK (char_length(name) > 0)
);
```

```sql
CREATE TABLE IF NOT EXISTS racks (
	id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	name TEXT NOT NULL
	
	CONSTRAINT name_positive_size CHECK (char_length(name) > 0)
);
```

```sql
CREATE TABLE IF NOT EXISTS orders (
	id INT NOT NULL,
	product_id INT,
	amount INT,
	
	CONSTRAINT id_positive CHECK (id > 0),
	CONSTRAINT fk_product
		FOREIGN KEY(product_id)
			REFERENCES products(id),
	CONSTRAINT amount_positive CHECK (amount > 0)
);
```

```sql
CREATE TABLE IF NOT EXISTS storage (
	product_id INT,
	rack_id INT,
	is_rack_primary BOOL,
	
	CONSTRAINT fk_product
		FOREIGN KEY(product_id)
			REFERENCES products(id),
	CONSTRAINT fk_rack
		FOREIGN KEY(rack_id)
			REFERENCES racks(id)
);
```

## Запросы на вставку данных

```sql
INSERT INTO products (id, name) VALUES (1, 'Ноутбук');
INSERT INTO products (id, name) VALUES (2, 'Телевизор');
INSERT INTO products (id, name) VALUES (3, 'Телефон');
INSERT INTO products (id, name) VALUES (4, 'Системный блок');
INSERT INTO products (id, name) VALUES (5, 'Часы');
INSERT INTO products (id, name) VALUES (6, 'Микрофон');
```

```sql
INSERT INTO racks (name) VALUES ('А');
INSERT INTO racks (name) VALUES ('Б');
INSERT INTO racks (name) VALUES ('В');
INSERT INTO racks (name) VALUES ('З');
INSERT INTO racks (name) VALUES ('Ж');
```

```sql
INSERT INTO orders (id, product_id, amount) VALUES (10, 1, 2);
INSERT INTO orders (id, product_id, amount) VALUES (10, 3, 1);
INSERT INTO orders (id, product_id, amount) VALUES (10, 6, 1);

INSERT INTO orders (id, product_id, amount) VALUES (11, 2, 3);

INSERT INTO orders (id, product_id, amount) VALUES (14, 1, 3);
INSERT INTO orders (id, product_id, amount) VALUES (14, 4, 4);

INSERT INTO orders (id, product_id, amount) VALUES (15, 5, 1);
```

```sql
INSERT INTO storage (product_id, rack_id, is_rack_primary) VALUES (1, 1, true);

INSERT INTO storage (product_id, rack_id, is_rack_primary) VALUES (2, 1, true);

INSERT INTO storage (product_id, rack_id, is_rack_primary) VALUES (3, 2, true);
INSERT INTO storage (product_id, rack_id, is_rack_primary) VALUES (3, 3, false);
INSERT INTO storage (product_id, rack_id, is_rack_primary) VALUES (3, 4, false);

INSERT INTO storage (product_id, rack_id, is_rack_primary) VALUES (4, 5, true);

INSERT INTO storage (product_id, rack_id, is_rack_primary) VALUES (5, 5, true);
INSERT INTO storage (product_id, rack_id, is_rack_primary) VALUES (5, 1, false);

INSERT INTO storage (product_id, rack_id, is_rack_primary) VALUES (6, 5, true);
```
