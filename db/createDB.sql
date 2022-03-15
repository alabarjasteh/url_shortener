DROP DATABASE IF EXISTS `alibaba_interview`;
CREATE DATABASE `alibaba_interview`;
USE `alibaba_interview`;

CREATE TABLE `pastes` (
  `shortlink` char(7) NOT NULL,
  `originallink` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY(`shortlink`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE INDEX ts ON pastes (created_at);