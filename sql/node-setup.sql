
create user cpmtest login superuser password 'cpmtest';

create user pgpool login superuser password 'pgpool';

create database cpmtest;

grant all privileges on database cpmtest to pgpool;
grant all privileges on database cpmtest to cpmtest;

\c cpmtest;

create extension adminpack;


DROP TABLE IF EXISTS loadtest;
CREATE TABLE loadtest(
	id serial primary key,
	name varchar(200)
);

DROP FUNCTION loadtest(integer, text);
CREATE OR REPLACE FUNCTION loadtest(writes integer, msg text)
    RETURNS VOID LANGUAGE plpgsql AS $PROC$
DECLARE
	timer1 timestamp;
	timer2 timestamp;
	elapsed numeric(18,3);
BEGIN
	timer1 := clock_timestamp();

	IF writes > 0 THEN
		FOR i IN 1 .. writes LOOP
			EXECUTE 'INSERT INTO loadtest (name) VALUES ('''||msg||''')';
		END LOOP;
	END IF;

	timer2 := clock_timestamp();
	elapsed := cast(extract(milliseconds from (timer2 - timer1)) as numeric(18,3));

	EXECUTE 'INSERT INTO loadtestresults (operation, count, elapsed)' ||
       		' VALUES (''inserts'',' || writes|| ',' ||  elapsed || ')';

	timer1 := clock_timestamp();
	IF writes > 0 THEN
		FOR i IN 1 .. writes LOOP
			EXECUTE 'SELECT name FROM loadtest WHERE id = '|| i;
		END LOOP;
	END IF;
	timer2 := clock_timestamp();
	elapsed := cast(extract(milliseconds from (timer2 - timer1)) as numeric(18,3));

	EXECUTE 'INSERT INTO loadtestresults (operation, count, elapsed)' ||
       		' VALUES (''selects'',' || writes|| ',' ||  elapsed || ')';
	
	
	timer1 := clock_timestamp();
	IF writes > 0 THEN
		FOR i IN 1 .. writes LOOP
			EXECUTE 'UPDATE loadtest SET (name) = (''howdy'') WHERE id = '||i;
		END LOOP;
	END IF;
	timer2 := clock_timestamp();
	elapsed := cast(extract(milliseconds from (timer2 - timer1)) as numeric(18,3));

	EXECUTE 'INSERT INTO loadtestresults (operation, count, elapsed)' ||
       		' VALUES (''updates'',' || writes|| ',' ||  elapsed || ')';
	
	timer1 := clock_timestamp();
	IF writes > 0 THEN
		FOR i IN 1 .. writes LOOP
			EXECUTE 'DELETE FROM loadtest  WHERE id = '||i;
		END LOOP;
	END IF;
	timer2 := clock_timestamp();
	elapsed := cast(extract(milliseconds from (timer2 - timer1)) as numeric(18,3));

	EXECUTE 'INSERT INTO loadtestresults (operation, count, elapsed)' ||
       		' VALUES (''deletes'',' || writes|| ',' ||  elapsed || ')';

END;
$PROC$;
