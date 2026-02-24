require("oltp_common")

sysbench.cmdline.options.sort_limit =
  {"Maximum rows returned per ORDER BY (0 = unlimited)", 0}

function prepare_statements() end

function event()
  local table_name = "sbtest" .. sysbench.rand.uniform(1, sysbench.opt.tables)
  local limit_clause = ""
  if sysbench.opt.sort_limit > 0 then
    limit_clause = " LIMIT " .. sysbench.opt.sort_limit
  end
  con:query(string.format(
    "SELECT k, COUNT(*) as cnt FROM %s GROUP BY k ORDER BY cnt DESC%s",
    table_name, limit_clause
  ))
end
