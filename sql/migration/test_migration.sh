#!/usr/bin/env sh

mysql -u root < ./000_old_schema.sql
mysql -u root < ./001_New_parent_table_for_executions_125.sql
