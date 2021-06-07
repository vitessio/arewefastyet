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

--
-- Table structure for table `execution`
--

DROP TABLE IF EXISTS `execution`;
CREATE TABLE `execution` (
                             `uuid` varchar(100) NOT NULL,
                             `status` varchar(100) DEFAULT 'created',
                             `started_at` datetime DEFAULT NULL,
                             `finished_at` datetime DEFAULT NULL,
                             `source` varchar(100) DEFAULT NULL,
                             `git_ref` varchar(100) DEFAULT NULL,
                             `type` varchar(100) DEFAULT '',
                             `pull_nb` int(11) DEFAULT 0,
                             PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Table structure for table `macrobenchmark`
--

DROP TABLE IF EXISTS `macrobenchmark`;
CREATE TABLE `macrobenchmark` (
                                  `macrobenchmark_id` int(11) NOT NULL AUTO_INCREMENT,
                                  `commit` varchar(100) DEFAULT NULL,
                                  `DateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                  `source` varchar(100) DEFAULT NULL,
                                  `exec_uuid` varchar(100) DEFAULT NULL,
                                  `vtgate_planner_version` varchar(20) NOT NULL DEFAULT 'V3',
                                  PRIMARY KEY (`macrobenchmark_id`)
#                                   KEY `macrobenchmark_ibfk_1` (`exec_uuid`),
#                                   CONSTRAINT `macrobenchmark_ibfk_1` FOREIGN KEY (`exec_uuid`) REFERENCES `execution` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=939 DEFAULT CHARSET=utf8;

--
-- Table structure for table `OLTP`
--

DROP TABLE IF EXISTS `OLTP`;
CREATE TABLE `OLTP` (
                        `OLTP_no` int(11) NOT NULL AUTO_INCREMENT,
                        `macrobenchmark_id` int(11) DEFAULT NULL,
                        `tps` decimal(8,2) DEFAULT NULL,
                        `latency` decimal(8,2) DEFAULT NULL,
                        `errors` decimal(8,2) DEFAULT NULL,
                        `reconnects` decimal(8,2) DEFAULT NULL,
                        `time` int(11) DEFAULT NULL,
                        `threads` decimal(8,2) DEFAULT NULL,
                        PRIMARY KEY (`OLTP_no`)
#                         KEY `test_no` (`macrobenchmark_id`),
#                         CONSTRAINT `OLTP_ibfk_1` FOREIGN KEY (`macrobenchmark_id`) REFERENCES `macrobenchmark` (`macrobenchmark_id`)
) ENGINE=InnoDB AUTO_INCREMENT=511 DEFAULT CHARSET=utf8;

--
-- Table structure for table `TPCC`
--

DROP TABLE IF EXISTS `TPCC`;
CREATE TABLE `TPCC` (
                        `TPCC_no` int(11) NOT NULL AUTO_INCREMENT,
                        `macrobenchmark_id` int(11) DEFAULT NULL,
                        `tps` decimal(8,2) DEFAULT NULL,
                        `latency` decimal(8,2) DEFAULT NULL,
                        `errors` decimal(8,2) DEFAULT NULL,
                        `reconnects` decimal(8,2) DEFAULT NULL,
                        `time` int(11) DEFAULT NULL,
                        `threads` decimal(8,2) DEFAULT NULL,
                        PRIMARY KEY (`TPCC_no`)
#                         KEY `test_no` (`macrobenchmark_id`),
#                         CONSTRAINT `TPCC_ibfk_1` FOREIGN KEY (`macrobenchmark_id`) REFERENCES `macrobenchmark` (`macrobenchmark_id`)
) ENGINE=InnoDB AUTO_INCREMENT=356 DEFAULT CHARSET=utf8;

--
-- Table structure for table `qps`
--

DROP TABLE IF EXISTS `qps`;
CREATE TABLE `qps` (
                       `qps_no` int(11) NOT NULL AUTO_INCREMENT,
                       `TPCC_no` int(11) DEFAULT NULL,
                       `total_qps` decimal(8,2) DEFAULT NULL,
                       `reads_qps` decimal(8,2) DEFAULT NULL,
                       `writes_qps` decimal(8,2) DEFAULT NULL,
                       `other_qps` decimal(8,2) DEFAULT NULL,
                       `OLTP_no` int(11) DEFAULT NULL,
                       PRIMARY KEY (`qps_no`)
#                        KEY `TPCC_no` (`TPCC_no`),
#                        KEY `OLTP_no` (`OLTP_no`),
#                        CONSTRAINT `qps_ibfk_1` FOREIGN KEY (`TPCC_no`) REFERENCES `TPCC` (`TPCC_no`),
#                        CONSTRAINT `qps_ibfk_2` FOREIGN KEY (`OLTP_no`) REFERENCES `OLTP` (`OLTP_no`)
) ENGINE=InnoDB AUTO_INCREMENT=865 DEFAULT CHARSET=utf8;

--
-- Table structure for table `microbenchmark`
--

DROP TABLE IF EXISTS `microbenchmark`;
CREATE TABLE `microbenchmark` (
                                  `microbenchmark_no` int(11) NOT NULL AUTO_INCREMENT,
                                  `pkg_name` varchar(255) DEFAULT NULL,
                                  `name` varchar(255) DEFAULT NULL,
                                  `git_ref` varchar(255) DEFAULT NULL,
                                  `exec_uuid` varchar(100) DEFAULT NULL,
                                  PRIMARY KEY (`microbenchmark_no`)
#                                   KEY `microbenchmark_ibfk_1` (`exec_uuid`),
#                                   CONSTRAINT `microbenchmark_ibfk_1` FOREIGN KEY (`exec_uuid`) REFERENCES `execution` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=14619 DEFAULT CHARSET=utf8;

--
-- Table structure for table `microbenchmark_details`
--

DROP TABLE IF EXISTS `microbenchmark_details`;
CREATE TABLE `microbenchmark_details` (
                                          `detail_no` int(11) NOT NULL AUTO_INCREMENT,
                                          `microbenchmark_no` int(11) DEFAULT NULL,
                                          `name` varchar(255) DEFAULT NULL,
                                          `bench_type` varchar(255) DEFAULT NULL,
                                          `n` int(11) DEFAULT NULL,
                                          `ns_per_op` decimal(22,5) DEFAULT NULL,
                                          `mb_per_sec` decimal(22,5) DEFAULT NULL,
                                          `bytes_per_op` decimal(22,5) DEFAULT NULL,
                                          `allocs_per_op` decimal(22,5) DEFAULT NULL,
                                          PRIMARY KEY (`detail_no`)
#                                           KEY `microbenchmark_details_fk_1` (`microbenchmark_no`),
#                                           CONSTRAINT `microbenchmark_details_fk_1` FOREIGN KEY (`microbenchmark_no`) REFERENCES `microbenchmark` (`microbenchmark_no`)
) ENGINE=InnoDB AUTO_INCREMENT=358174 DEFAULT CHARSET=utf8;
