CREATE DATABASE IF NOT EXISTS benchmark;

USE benchmark;

CREATE TABLE IF NOT EXISTS execution (
    `uuid` VARCHAR(100) NOT NULL,
    `status` VARCHAR(100) DEFAULT 'created',
    `started_at` TIMESTAMP DEFAULT NULL,
    `finished_at` TIMESTAMP DEFAULT NULL,
    `source` VARCHAR(100) DEFAULT NULL,
    PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS macrobenchmark (
  `macrobenchmark_id` int(11) NOT NULL AUTO_INCREMENT,
  `exec_uuid` VARCHAR(100) DEFAULT NULL,
  `commit` varchar(100) DEFAULT NULL,
  `DateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `source` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`macrobenchmark_id`),
  KEY `exec_uuid` (`exec_uuid`),
  CONSTRAINT `macrobenchmark_ibfk_1` FOREIGN KEY (`exec_uuid`) REFERENCES execution (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `OLTP` (
  `OLTP_no` int(11) NOT NULL AUTO_INCREMENT,
  `macrobenchmark_id` int(11) DEFAULT NULL,
  `tps` decimal(8,2) DEFAULT NULL,
  `latency` decimal(8,2) DEFAULT NULL,
  `errors` decimal(8,2) DEFAULT NULL,
  `reconnects` decimal(8,2) DEFAULT NULL,
  `time` int(11) DEFAULT NULL,
  `threads` decimal(8,2) DEFAULT NULL,
  PRIMARY KEY (`OLTP_no`),
  KEY `macrobenchmark_id` (`macrobenchmark_id`),
  CONSTRAINT `OLTP_ibfk_1` FOREIGN KEY (`macrobenchmark_id`) REFERENCES macrobenchmark (`macrobenchmark_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `TPCC` (
  `TPCC_no` int(11) NOT NULL AUTO_INCREMENT,
  `macrobenchmark_id` int(11) DEFAULT NULL,
  `tps` decimal(8,2) DEFAULT NULL,
  `latency` decimal(8,2) DEFAULT NULL,
  `errors` decimal(8,2) DEFAULT NULL,
  `reconnects` decimal(8,2) DEFAULT NULL,
  `time` int(11) DEFAULT NULL,
  `threads` decimal(8,2) DEFAULT NULL,
  PRIMARY KEY (`TPCC_no`),
  KEY `macrobenchmark_id` (`macrobenchmark_id`),
  CONSTRAINT `TPCC_ibfk_1` FOREIGN KEY (`macrobenchmark_id`) REFERENCES macrobenchmark (`macrobenchmark_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE IF NOT EXISTS `qps` (
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

CREATE TABLE IF NOT EXISTS `microbenchmark` (
    `microbenchmark_no` INT AUTO_INCREMENT,
    `exec_uuid` VARCHAR(100) DEFAULT NULL,
    `pkg_name` VARCHAR(255),
    `name` VARCHAR(255),
    PRIMARY KEY (`microbenchmark_no`),
    KEY `exec_uuid` (`exec_uuid`),
    CONSTRAINT `microbenchmark_ibfk_1` FOREIGN KEY (`exec_uuid`) REFERENCES execution (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `microbenchmark_details` (
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