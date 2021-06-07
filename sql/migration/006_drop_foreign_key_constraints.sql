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

ALTER TABLE macrobenchmark DROP FOREIGN KEY macrobenchmark_ibfk_1;
ALTER TABLE OLTP DROP FOREIGN KEY OLTP_ibfk_1;
ALTER TABLE TPCC DROP FOREIGN KEY TPCC_ibfk_1;
ALTER TABLE qps DROP FOREIGN KEY qps_ibfk_1;
ALTER TABLE qps DROP FOREIGN KEY qps_ibfk_2;
ALTER TABLE microbenchmark DROP FOREIGN KEY microbenchmark_ibfk_1;
ALTER TABLE microbenchmark_details DROP FOREIGN KEY microbenchmark_details_fk_1;
