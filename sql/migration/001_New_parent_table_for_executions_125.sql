/*
 *
 * Copyright 2021 The Vitess Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

USE benchmark;

CREATE TABLE IF NOT EXISTS execution (
    `uuid` VARCHAR(100) NOT NULL,
    `status` VARCHAR(100) DEFAULT 'created',
    `started_at` DATETIME NULL,
    `finished_at` DATETIME NULL,
    `source` VARCHAR(100) DEFAULT NULL,
    `git_ref` VARCHAR(100) DEFAULT NULL,
    PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

alter table benchmark rename to macrobenchmark;

alter table macrobenchmark change column test_no macrobenchmark_id int(11) NOT NULL AUTO_INCREMENT;
alter table macrobenchmark add column exec_uuid VARCHAR(100) DEFAULT NULL;
alter table macrobenchmark add constraint  macrobenchmark_ibfk_1 FOREIGN KEY (exec_uuid) REFERENCES execution (uuid);

alter table OLTP change column test_no macrobenchmark_id int(11) DEFAULT NULL;
alter table TPCC change column test_no macrobenchmark_id int(11) DEFAULT NULL;

alter table microbenchmark add column exec_uuid VARCHAR(100) DEFAULT NULL;
alter table microbenchmark add constraint  microbenchmark_ibfk_1 FOREIGN KEY (exec_uuid) REFERENCES execution (uuid);

alter table microbenchmark drop column test_no;
