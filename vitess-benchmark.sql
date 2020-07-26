-- MySQL dump 10.13  Distrib 8.0.20, for Linux (x86_64)
--
-- Host: localhost    Database: vitess_benchmark
-- ------------------------------------------------------
-- Server version	8.0.20-cluster

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

--
-- Table structure for table `OLTP`
--

DROP TABLE IF EXISTS `OLTP`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `OLTP` (
  `OLTP_no` int NOT NULL AUTO_INCREMENT,
  `test_no` int DEFAULT NULL,
  `tps` decimal(8,2) DEFAULT NULL,
  `latency` decimal(8,2) DEFAULT NULL,
  `errors` decimal(8,2) DEFAULT NULL,
  `reconnects` decimal(8,2) DEFAULT NULL,
  `time` int DEFAULT NULL,
  `threads` int DEFAULT NULL,
  PRIMARY KEY (`OLTP_no`),
  KEY `test_no` (`test_no`),
  CONSTRAINT `OLTP_ibfk_1` FOREIGN KEY (`test_no`) REFERENCES `benchmark` (`test_no`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `OLTP`
--

LOCK TABLES `OLTP` WRITE;
/*!40000 ALTER TABLE `OLTP` DISABLE KEYS */;
/*!40000 ALTER TABLE `OLTP` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `TPCC`
--

DROP TABLE IF EXISTS `TPCC`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `TPCC` (
  `TPCC_no` int NOT NULL AUTO_INCREMENT,
  `test_no` int DEFAULT NULL,
  `tps` decimal(8,2) DEFAULT NULL,
  `latency` decimal(8,2) DEFAULT NULL,
  `errors` decimal(8,2) DEFAULT NULL,
  `reconnects` decimal(8,2) DEFAULT NULL,
  PRIMARY KEY (`TPCC_no`),
  KEY `test_no` (`test_no`),
  CONSTRAINT `TPCC_ibfk_1` FOREIGN KEY (`test_no`) REFERENCES `benchmark` (`test_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `TPCC`
--

LOCK TABLES `TPCC` WRITE;
/*!40000 ALTER TABLE `TPCC` DISABLE KEYS */;
/*!40000 ALTER TABLE `TPCC` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `benchmark`
--

DROP TABLE IF EXISTS `benchmark`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `benchmark` (
  `test_no` int NOT NULL AUTO_INCREMENT,
  `commit` varchar(100) DEFAULT NULL,
  `DateTime` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`test_no`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `benchmark`
--

LOCK TABLES `benchmark` WRITE;
/*!40000 ALTER TABLE `benchmark` DISABLE KEYS */;
/*!40000 ALTER TABLE `benchmark` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `qps`
--

DROP TABLE IF EXISTS `qps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `qps` (
  `qps_no` int NOT NULL AUTO_INCREMENT,
  `TPCC_no` int DEFAULT NULL,
  `total_qps` decimal(8,2) DEFAULT NULL,
  `reads_qps` decimal(8,2) DEFAULT NULL,
  `writes_qps` decimal(8,2) DEFAULT NULL,
  `other_qps` decimal(8,2) DEFAULT NULL,
  `OLTP_no` int DEFAULT NULL,
  PRIMARY KEY (`qps_no`),
  KEY `TPCC_no` (`TPCC_no`),
  KEY `OLTP_no` (`OLTP_no`),
  CONSTRAINT `qps_ibfk_1` FOREIGN KEY (`TPCC_no`) REFERENCES `TPCC` (`TPCC_no`),
  CONSTRAINT `qps_ibfk_2` FOREIGN KEY (`OLTP_no`) REFERENCES `OLTP` (`OLTP_no`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `qps`
--

LOCK TABLES `qps` WRITE;
/*!40000 ALTER TABLE `qps` DISABLE KEYS */;
/*!40000 ALTER TABLE `qps` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-07-18 21:18:54
