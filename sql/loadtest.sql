CREATE TEMP TABLE loadtestresults (
	operation varchar(20),
	count int,
	elapsed numeric(18,3)
);

SELECT loadtest(:writes, 'loadtestdata');

COPY (
SELECT array_to_json(array_agg(row_to_json(t)))
FROM (
	select operation, count, elapsed as results from loadtestresults
)t
) TO stdout;

DROP TABLE loadtestresults;
