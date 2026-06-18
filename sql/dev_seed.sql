-- Dev-only seed data for the local docker-compose `db` service.
--
-- This is NOT production data: it generates a small, internally-consistent set of
-- finished cron benchmarks for the OLTP, TPCC and TPCC-OLAP workloads so the
-- server-rendered pages (Status, Home sparklines, Foreign Keys compare, …) show
-- real content when developing locally. All rows use pull_nb = 0 so the PR page
-- does not try to reach GitHub with the throwaway dev credentials.
--
-- Loaded automatically on first DB init (see docker-compose.yml). To re-seed:
--   docker compose down && docker volume rm arewefastyet_aafy-db && docker compose up -d

USE arewefastyet;

DELIMITER $$

DROP PROCEDURE IF EXISTS seed_dev_data $$

CREATE PROCEDURE seed_dev_data()
BEGIN
  DECLARE wi INT DEFAULT 0;
  DECLARE di INT;
  DECLARE ri INT;
  DECLARE wl VARCHAR(20);
  DECLARE base_qps DOUBLE;
  DECLARE gitref VARCHAR(100);
  DECLARE u VARCHAR(100);
  DECLARE mbid INT;
  DECLARE st DATETIME;
  DECLARE fin DATETIME;
  DECLARE qps DOUBLE;

  -- 3 workloads x 8 commits (one per recent day) x 6 runs each.
  WHILE wi < 3 DO
    SET wl = ELT(wi + 1, 'OLTP', 'TPCC', 'TPCC-OLAP');
    SET base_qps = ELT(wi + 1, 1200, 900, 700);
    SET di = 0;
    WHILE di < 8 DO
      -- Same synthetic commit across workloads for a given day, so the Foreign
      -- Keys page can compare two workloads on one commit.
      SET gitref = SHA1(CONCAT('aafy-ref-', di));
      SET ri = 0;
      WHILE ri < 6 DO
        SET u = UUID();
        SET fin = NOW() - INTERVAL di DAY - INTERVAL (ri * 7) MINUTE;
        SET st = fin - INTERVAL (5 + ri) MINUTE;
        SET qps = base_qps + di * 4 + ri * 3 - (ri MOD 2) * 7;

        INSERT INTO execution
          (uuid, status, started_at, finished_at, source, git_ref, workload, pull_nb, go_version, profile_binary, profile_mode)
        VALUES
          (u, 'finished', st, fin, 'cron', gitref, wl, 0, '1.24.5', NULL, NULL);

        INSERT INTO macrobenchmark (`commit`, exec_uuid, vtgate_planner_version, workload)
        VALUES (gitref, u, 'Gen4', wl);
        SET mbid = LAST_INSERT_ID();

        INSERT INTO macrobenchmark_results
          (macrobenchmark_id, tps, latency, errors, reconnects, `time`, threads, total_qps, reads_qps, writes_qps, other_qps, queries)
        VALUES
          (mbid, qps * 0.5, 2.5 + ri * 0.1, 0, 0, 120, 16, qps, qps * 0.7, qps * 0.2, qps * 0.1, FLOOR(qps * 120));

        INSERT INTO metrics (exec_uuid, name, value, description) VALUES
          (u, 'TotalComponentsCPUTime', 0.00012 + ri * 0.000003, ''),
          (u, 'ComponentsCPUTime.vtgate', 0.00005 + ri * 0.000001, ''),
          (u, 'ComponentsCPUTime.vttablet', 0.00007 + ri * 0.000002, ''),
          (u, 'TotalComponentsMemStatsAllocBytes', 5000000 + ri * 12000, ''),
          (u, 'ComponentsMemStatsAllocBytes.vtgate', 2000000 + ri * 5000, ''),
          (u, 'ComponentsMemStatsAllocBytes.vttablet', 3000000 + ri * 7000, '');

        -- A handful of per-query plans so the Compare Query Plans page renders
        -- real rows. Keys are valid SQL (the server parses them with the Vitess
        -- sqlparser); exec_time varies by commit/day (di) and moves in different
        -- directions per query, so old↔new comparisons show non-zero, mixed deltas.
        INSERT INTO query_plans (exec_uuid, macrobenchmark_id, `key`, `plan`, exec_count, exec_time, `rows`, errors) VALUES
          (u, mbid, 'select * from warehouse where w_id = 1',
            '{"OperatorType":"Route","Variant":"EqualUnique","Keyspace":{"Name":"main","Sharded":true},"Query":"select * from warehouse where w_id = 1"}',
            100, 80 + di * 5, 1, 0),
          (u, mbid, 'insert into history(h_data) values (''seed'')',
            '{"OperatorType":"Insert","Variant":"Sharded","Keyspace":{"Name":"main","Sharded":true},"TargetTabletType":"PRIMARY"}',
            50, 120 + di * 8, 1, 0),
          (u, mbid, 'update customer set c_balance = 1 where c_id = 2',
            '{"OperatorType":"Update","Variant":"Equal","Keyspace":{"Name":"main","Sharded":true}}',
            40, 200 - di * 6, 1, 0),
          (u, mbid, 'delete from new_order where no_o_id = 3',
            '{"OperatorType":"Delete","Variant":"Equal","Keyspace":{"Name":"main","Sharded":true}}',
            30, 150 + di * 3, 0, 0);

        SET ri = ri + 1;
      END WHILE;
      SET di = di + 1;
    END WHILE;
    SET wi = wi + 1;
  END WHILE;

  -- A little status/source variety for the "Previous Executions" table.
  INSERT INTO execution
    (uuid, status, started_at, finished_at, source, git_ref, workload, pull_nb, go_version, profile_binary, profile_mode)
  VALUES
    (UUID(), 'failed',   NOW() - INTERVAL 2 HOUR,    NULL,                                          'cron',      SHA1('aafy-ref-0'), 'OLTP',      0, '1.24.5', NULL, NULL),
    (UUID(), 'started',  NOW() - INTERVAL 10 MINUTE, NULL,                                          'cron_tags', SHA1('aafy-ref-1'), 'TPCC',      0, '1.24.5', NULL, NULL),
    (UUID(), 'finished', NOW() - INTERVAL 1 DAY,     NOW() - INTERVAL 1 DAY + INTERVAL 9 MINUTE,    'cron_tags', SHA1('aafy-ref-2'), 'TPCC-OLAP', 0, '1.24.5', NULL, NULL);
END $$

DELIMITER ;

CALL seed_dev_data();
DROP PROCEDURE seed_dev_data;
