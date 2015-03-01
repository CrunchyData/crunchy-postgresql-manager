select row_to_json(t)
from (
SELECT
  to_char(now(), 'mm/dd/yy HH12:MI:SS') now,
  to_char(block_size::numeric * buffers_alloc / (1024 * 1024 * seconds), 'FM999999999999D9999') AS alloc_mbps,
  to_char(block_size::numeric * buffers_checkpoint / (1024 * 1024 * seconds), 'FM999999999999D9999') AS checkpoint_mbps,
  to_char(block_size::numeric * buffers_clean / (1024 * 1024 * seconds), 'FM999999999999D9999') AS clean_mbps,
  to_char(block_size::numeric * buffers_backend/ (1024 * 1024 * seconds), 'FM999999999999D9999') AS backend_mbps,
  to_char(block_size::numeric * (buffers_checkpoint + buffers_clean + buffers_backend) / (1024 * 1024 * seconds), 'FM999999999999D9999') AS write_mbps  
FROM
(
SELECT now() AS sample,now() - stats_reset AS uptime,EXTRACT(EPOCH FROM now()) - extract(EPOCH FROM stats_reset) AS seconds, b.*,p.setting::integer AS block_size FROM pg_stat_bgwriter b,pg_settings p WHERE p.name='block_size'
) bgw

) t

