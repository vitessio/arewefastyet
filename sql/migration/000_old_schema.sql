DROP DATABASE IF EXISTS benchmark;

CREATE DATABASE benchmark;

USE benchmark;

CREATE TABLE `benchmark` (
                             `test_no` int(11) NOT NULL AUTO_INCREMENT,
                             `commit` varchar(100) DEFAULT NULL,
                             `DateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             `source` varchar(100) DEFAULT NULL,
                             PRIMARY KEY (`test_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `OLTP` (
                        `OLTP_no` int(11) NOT NULL AUTO_INCREMENT,
                        `test_no` int(11) DEFAULT NULL,
                        `tps` decimal(8,2) DEFAULT NULL,
                        `latency` decimal(8,2) DEFAULT NULL,
                        `errors` decimal(8,2) DEFAULT NULL,
                        `reconnects` decimal(8,2) DEFAULT NULL,
                        `time` int(11) DEFAULT NULL,
                        `threads` decimal(8,2) DEFAULT NULL,
                        PRIMARY KEY (`OLTP_no`),
                        KEY `test_no` (`test_no`),
                        CONSTRAINT `OLTP_ibfk_1` FOREIGN KEY (`test_no`) REFERENCES `benchmark` (`test_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `TPCC` (
                        `TPCC_no` int(11) NOT NULL AUTO_INCREMENT,
                        `test_no` int(11) DEFAULT NULL,
                        `tps` decimal(8,2) DEFAULT NULL,
                        `latency` decimal(8,2) DEFAULT NULL,
                        `errors` decimal(8,2) DEFAULT NULL,
                        `reconnects` decimal(8,2) DEFAULT NULL,
                        `time` int(11) DEFAULT NULL,
                        `threads` decimal(8,2) DEFAULT NULL,
                        PRIMARY KEY (`TPCC_no`),
                        KEY `test_no` (`test_no`),
                        CONSTRAINT `TPCC_ibfk_1` FOREIGN KEY (`test_no`) REFERENCES `benchmark` (`test_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `qps` (
                       `qps_no` int(11) NOT NULL AUTO_INCREMENT,
                       `TPCC_no` int(11) DEFAULT NULL,
                       `total_qps` decimal(8,2) DEFAULT NULL,
                       `reads_qps` decimal(8,2) DEFAULT NULL,
                       `writes_qps` decimal(8,2) DEFAULT NULL,
                       `other_qps` decimal(8,2) DEFAULT NULL,
                       `OLTP_no` int(11) DEFAULT NULL,
                       PRIMARY KEY (`qps_no`),
                       KEY `TPCC_no` (`TPCC_no`),
                       KEY `OLTP_no` (`OLTP_no`),
                       CONSTRAINT `qps_ibfk_1` FOREIGN KEY (`TPCC_no`) REFERENCES `TPCC` (`TPCC_no`),
                       CONSTRAINT `qps_ibfk_2` FOREIGN KEY (`OLTP_no`) REFERENCES `OLTP` (`OLTP_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `microbenchmark` (
                                  `microbenchmark_no` INT AUTO_INCREMENT,
                                  `test_no` INT NOT NULL,
                                  `pkg_name` VARCHAR(255),
                                  `name` VARCHAR(255),
                                  PRIMARY KEY (`microbenchmark_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `microbenchmark_details` (
                                          `detail_no` INT AUTO_INCREMENT,
                                          `microbenchmark_no` INT,
                                          `name` VARCHAR(255),
                                          `bench_type` VARCHAR(255),
                                          `n` INT,
                                          `ns_per_op` DECIMAL(22,5),
                                          `mb_per_sec` DECIMAL(22,5),
                                          `bytes_per_op` DECIMAL(22,5),
                                          `allocs_per_op` DECIMAL(22,5),
                                          PRIMARY KEY (`detail_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;