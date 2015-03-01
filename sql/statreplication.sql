select array_to_json(array_agg(row_to_json(t)))
from (
	SELECT
	 pid , usesysid , usename , application_name , client_addr , client_hostname ,
	client_port , to_char(backend_start, 'YYYY-MM-DD HH24:MI-SS') as backend_start , state , sent_location , write_location ,
	flush_location , replay_location , sync_priority , sync_state from pg_stat_replication

) t

