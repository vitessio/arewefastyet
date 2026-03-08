-- MySQL dump 10.13  Distrib 8.4.8, for macos15 (arm64)
--
-- Host: localhost    Database: arewefastyet
-- ------------------------------------------------------
-- Server version	8.0.34-psdbproxy

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN = 0;

--
-- Table structure for table `execution`
--

DROP TABLE IF EXISTS `execution`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `execution` (
  `uuid` varchar(100) NOT NULL,
  `status` varchar(100) DEFAULT 'created',
  `started_at` datetime DEFAULT NULL,
  `finished_at` datetime DEFAULT NULL,
  `source` varchar(100) DEFAULT NULL,
  `git_ref` varchar(100) DEFAULT NULL,
  `workload` varchar(100) DEFAULT NULL,
  `pull_nb` int DEFAULT '0',
  `go_version` varchar(16) DEFAULT NULL,
  `profile_binary` varchar(20) DEFAULT NULL,
  `profile_mode` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`uuid`),
  KEY `finished_at` (`finished_at`,`status`),
  KEY `started_at` (`started_at` DESC),
  KEY `status` (`status`),
  KEY `git_ref` (`git_ref`),
  KEY `pull_nb` (`pull_nb`),
  KEY `idx_execution_on_source` (`source`),
  KEY `idx_execution_on_profile_binary` (`profile_binary`),
  KEY `finished_at_2` (`finished_at`,`git_ref`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `macrobenchmark`
--

DROP TABLE IF EXISTS `macrobenchmark`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `macrobenchmark` (
  `macrobenchmark_id` int NOT NULL AUTO_INCREMENT,
  `commit` varchar(100) DEFAULT NULL,
  `DateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `exec_uuid` varchar(100) DEFAULT NULL,
  `vtgate_planner_version` varchar(20) NOT NULL DEFAULT 'V3',
  `workload` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`macrobenchmark_id`),
  KEY `exec_uuid` (`exec_uuid`),
  KEY `commit` (`commit`),
  KEY `type_vtgate_planner_version` (`workload`,`vtgate_planner_version`)
) ENGINE=InnoDB AUTO_INCREMENT=63142 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `macrobenchmark_results`
--

DROP TABLE IF EXISTS `macrobenchmark_results`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `macrobenchmark_results` (
  `id` int NOT NULL AUTO_INCREMENT,
  `macrobenchmark_id` int DEFAULT NULL,
  `tps` decimal(8,2) DEFAULT NULL,
  `latency` decimal(8,2) DEFAULT NULL,
  `errors` decimal(8,2) DEFAULT NULL,
  `reconnects` decimal(8,2) DEFAULT NULL,
  `time` int DEFAULT NULL,
  `threads` decimal(8,2) DEFAULT NULL,
  `total_qps` decimal(8,2) DEFAULT NULL,
  `reads_qps` decimal(8,2) DEFAULT NULL,
  `writes_qps` decimal(8,2) DEFAULT NULL,
  `other_qps` decimal(8,2) DEFAULT NULL,
  `queries` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `macrobenchmark_id` (`macrobenchmark_id`)
) ENGINE=InnoDB AUTO_INCREMENT=54062 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `metrics`
--

DROP TABLE IF EXISTS `metrics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `metrics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `exec_uuid` varchar(100) DEFAULT NULL,
  `name` varchar(250) DEFAULT NULL,
  `value` float DEFAULT NULL,
  `description` varchar(250) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `exec_uuid` (`exec_uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=332629 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `microbenchmark`
--

DROP TABLE IF EXISTS `microbenchmark`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `microbenchmark` (
  `microbenchmark_no` int NOT NULL AUTO_INCREMENT,
  `pkg_name` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `git_ref` varchar(255) DEFAULT NULL,
  `exec_uuid` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`microbenchmark_no`),
  KEY `git_ref` (`git_ref`),
  KEY `exec_uuid` (`exec_uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=155766 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `microbenchmark_details`
--

DROP TABLE IF EXISTS `microbenchmark_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `microbenchmark_details` (
  `detail_no` int NOT NULL AUTO_INCREMENT,
  `microbenchmark_no` int DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `bench_type` varchar(255) DEFAULT NULL,
  `n` int DEFAULT NULL,
  `ns_per_op` decimal(22,5) DEFAULT NULL,
  `mb_per_sec` decimal(22,5) DEFAULT NULL,
  `bytes_per_op` decimal(22,5) DEFAULT NULL,
  `allocs_per_op` decimal(22,5) DEFAULT NULL,
  PRIMARY KEY (`detail_no`),
  KEY `microbenchmark_no` (`microbenchmark_no`)
) ENGINE=InnoDB AUTO_INCREMENT=4272143 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `query_plans`
--

DROP TABLE IF EXISTS `query_plans`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `query_plans` (
  `plan_id` int NOT NULL AUTO_INCREMENT,
  `exec_uuid` varchar(100) DEFAULT NULL,
  `macrobenchmark_id` int DEFAULT NULL,
  `key` longtext,
  `plan` longtext,
  `exec_count` bigint DEFAULT NULL,
  `exec_time` bigint DEFAULT NULL,
  `rows` int DEFAULT NULL,
  `errors` int DEFAULT NULL,
  PRIMARY KEY (`plan_id`),
  KEY `macrobenchmark_id` (`macrobenchmark_id`),
  KEY `exec_uuid` (`exec_uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=5164979 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-03-08  4:21:03
