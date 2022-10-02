-- 1 задание
DO
$$ BEGIN
	raise notice 'Hello World!!'; 
END $$
language plpgsql;

-- 2 задание 
DO
$$ BEGIN
	raise notice 'Current date = %', NOW(); 
END $$
language plpgsql;

-- 3 задание
DO
$$ declare 
	x int;
	y int;
BEGIN
	x=10;
	y=5;
	raise notice 'x-y = %', x-y; 
	raise notice 'x+y = %', x+y;
	raise notice 'x/y = %', x/y; 
	raise notice 'x*y = %', x*y; 
END
$$ language plpgsql; 

-- 4 задание
DO 
$$ DECLARE
  mark int := 5;
BEGIN
    IF mark = 5 THEN 
        RAISE NOTICE 'Отлично!';
    ELSIF mark = 4 
        THEN RAISE NOTICE 'Хорошо';
  	ELSIF mark = 3 
        THEN RAISE NOTICE 'Удовлетворительно';
	ELSIF mark = 2 
        THEN RAISE NOTICE 'Неуд';
  	ELSE 
        RAISE NOTICE 'Введенная оценка не верна';
  END IF;
END
$$ language plpgsql;

DO
$$ DECLARE
  mark int := 5;
BEGIN
  CASE mark
  	WHEN 5 THEN 
		RAISE NOTICE 'Отлично!';
	WHEN 4 THEN
		RAISE NOTICE 'Хорошо';
	WHEN 3 THEN 
		RAISE NOTICE 'Удовлетворительно';
	WHEN 2 THEN 
		RAISE NOTICE 'Неуд';
	ELSE 
		RAISE NOTICE 'Введенная оценка не верна';
  END CASE;
END
$$ language plpgsql;

-- 5 задание
DO 
$$ DECLARE 
	x int := 20;
BEGIN
	WHILE x < 31 LOOP
		RAISE NOTICE '%',x^2;
		x =x + 1;
	END LOOP;
END
$$ language plpgsql;

DO 
$$ BEGIN
	FOR i IN 20..30 LOOP
		RAISE NOTICE '%', i^2;
	END LOOP;
END;
$$ language plpgsql;

DO 
$$ DECLARE 
	x int := 20;
BEGIN
	LOOP
		RAISE NOTICE '%', x^2;
		x = x + 1;
		IF(x = 31) THEN
			EXIT;
		END IF;
	END LOOP;
END
$$ language plpgsql;

-- 6 задание
CREATE OR REPLACE FUNCTION kollats(n int) RETURNS int
AS $$
DECLARE
	cnt int;
BEGIN
	cnt:=0;
	WHILE n<>1 LOOP
		IF n % 2 = 0 THEN
			n = n / 2;
		ELSE
			n = n * 3 + 1;
		END IF;
		cnt = cnt + 1;
	END LOOP;
	RETURN cnt;
END
$$ LANGUAGE plpgsql;
SELECT kollats(2)

CREATE OR REPLACE PROCEDURE kollats_proc(n int)
AS $$
DECLARE
BEGIN
	WHILE n<>1 LOOP
		RAISE NOTICE '%', n;
		IF n % 2 = 0 THEN
			n = n / 2;
		ELSE
			n = n * 3 + 1;
		END IF;
	END LOOP;
END
$$ LANGUAGE plpgsql;
CALL kollats_proc(5)

-- 7 задание
CREATE OR REPLACE FUNCTION luck(n int) RETURNS int
AS $$
DECLARE
	l0 int := 2;
	l1 int := 1;
	tmp int := 0;
BEGIN
	WHILE n <> 2 LOOP
		tmp = l0 + l1;
		l0 = L1;
		l1 = tmp;
		n = n - 1;
	END LOOP;
	RETURN l1;
END
$$ LANGUAGE plpgsql;
SELECT luck(5)

CREATE OR REPLACE PROCEDURE luck_proc(n int)
AS $$
DECLARE
	l0 int := 2;
	l1 int := 1;
	tmp int := 0;
BEGIN
	RAISE NOTICE '%', l0;
	RAISE NOTICE '%', l1;
	WHILE n <> 2 LOOP
		tmp = l0 + l1;
		l0 = L1;
		l1 = tmp;
		n = n - 1;
		RAISE NOTICE '%', l1;
	END LOOP;
END
$$ LANGUAGE plpgsql;
CALL luck_proc(5)

-- 8 задание 
CREATE OR REPLACE FUNCTION get_count_of_people_by_year(year int) RETURNS int
AS $$
DECLARE
	cnt int;
BEGIN
	SELECT COUNT(id) INTO cnt
	FROM people
	WHERE EXTRACT(year from people.birth_date) = get_count_of_people_by_year.year;
	RETURN cnt;
END
$$ LANGUAGE plpgsql;
SELECT get_count_of_people_by_year(1995)

-- 9 задание
CREATE OR REPLACE FUNCTION get_count_of_people_by_eye_color(color varchar) RETURNS int
AS $$
DECLARE
	cnt int;
BEGIN
	SELECT COUNT(eyes) INTO cnt
	FROM people
	WHERE people.eyes = color;
	RETURN cnt;
END
$$ LANGUAGE plpgsql;
SELECT get_count_of_people_by_eye_color('brown')

-- 10 задание
CREATE OR REPLACE FUNCTION youngest_man() RETURNS int
AS $$
DECLARE
	man_id int;
BEGIN
	SELECT id INTO man_id
	FROM people
	ORDER BY birth_date DESC
	LIMIT 1;
	RETURN man_id;
END
$$ LANGUAGE plpgsql;
SELECT youngest_man()

-- 11 задание
CREATE OR REPLACE PROCEDURE iwb()
AS $$
DECLARE
	p people%ROWTYPE;
BEGIN
	FOR p IN 
		SELECT * FROM people
		WHERE people.weight / ((people.growth / 100) ^ 2) >0.01
	LOOP
		RAISE NOTICE 'id: %, name: %, surname: %', p.id, p.name, p.surname;
	END LOOP;
END
$$ LANGUAGE plpgsql;
CALL iwb()

-- 12 задание
BEGIN;
CREATE TABLE blood_relations (
  id SERIAL PRIMARY key,
  relation_type VARCHAR(255) NOT NULL,
  people_id integer REFERENCES people(id));
COMMIT;

-- 13 задание 
CREATE OR REPLACE PROCEDURE 
add_person(IN add_name varchar, add_surname varchar, add_birth_date DATE, add_growth real, add_weight real, add_eyes varchar, add_hair varchar,add_relation_type varchar,add_people_id int)
AS $$
BEGIN
	INSERT INTO people (name, surname, birth_date, growth, weight, eyes, hair)
	VALUES (add_name, add_surname, add_birth_date, add_growth, add_weight, add_eyes, add_hair);
	INSERT INTO blood_relations(relation_type, people_id)
	VALUES (add_relation_type, add_people_id);
END;
$$ LANGUAGE plpgsql;

-- 14 задание
BEGIN;
ALTER TABLE people ADD COLUMN actual_data_time TIMESTAMP;
COMMIT;

-- 15 задание
CREATE OR REPLACE PROCEDURE update_physical_characteristics(id int, updt_growth real, updt_weight real)
AS $$
BEGIN
	UPDATE people
	SET growth = updt_growth, weight = updt_weight, actual_data_time = NOW()
	WHERE people.id = update_physical_characteristics.id;
END;
$$ LANGUAGE plpgsql;
CALL update_physical_characteristics(1,80,80)
.