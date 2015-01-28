select array_to_json(array_agg(row_to_json(t)))
from (
SELECT
 datname,
 blks_read,
 tup_returned,
 tup_fetched,
 tup_inserted,
 tup_updated,
 tup_deleted,
 to_char(stats_reset, 'YYYY-MM-DD HH24:MI:SS') as stats_reset

 from pg_stat_database
) t

