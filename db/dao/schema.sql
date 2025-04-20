
CREATE TABLE IF NOT EXISTS `author` (
  `authorId` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`authorId`),
  UNIQUE KEY `authorId_UNIQUE` (`authorId`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `book` (
  `bookId` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `format` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
  `file` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `hash` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,
  `updateTs` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `isbn` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `coverColor` varchar(7) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `edited` tinyint(1) NOT NULL,
  `hasCover` tinyint(1) NOT NULL,
  `coverType` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `blurb` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`bookId`),
  UNIQUE KEY `bookId_UNIQUE` (`bookId`),
  UNIQUE KEY `hash_UNIQUE` (`hash`),
  UNIQUE KEY `file_UNIQUE` (`file`),
  KEY `title` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `bookauthors` (
  `bookAuthorsId` int(11) NOT NULL AUTO_INCREMENT,
  `bookId` int(11) NOT NULL,
  `authorId` int(11) NOT NULL,
  PRIMARY KEY (`bookAuthorsId`),
  UNIQUE KEY `bookAuthorsId_UNIQUE` (`bookAuthorsId`),
  KEY `book` (`bookId`),
  KEY `author` (`authorId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `bookseries` (
  `bookSeriesId` int(11) NOT NULL AUTO_INCREMENT,
  `bookId` int(11) NOT NULL,
  `seriesId` int(11) NOT NULL,
  `sequence` int(11) DEFAULT NULL,
  PRIMARY KEY (`bookSeriesId`),
  UNIQUE KEY `bookSeriesId_UNIQUE` (`bookSeriesId`),
  KEY `book` (`bookId`),
  KEY `series` (`seriesId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `booktags` (
  `bookTagId` int(11) NOT NULL AUTO_INCREMENT,
  `bookId` int(11) NOT NULL,
  `tagId` int(11) NOT NULL,
  PRIMARY KEY (`bookTagId`),
  UNIQUE KEY `bookTagId_UNIQUE` (`bookTagId`),
  KEY `book` (`bookId`),
  KEY `tag` (`tagId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `series` (
  `seriesId` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `issn` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`seriesId`),
  UNIQUE KEY `seriesId_UNIQUE` (`seriesId`),
  UNIQUE KEY `title_UNIQUE` (`title`),
  KEY `issn` (`issn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tag` (
  `tagId` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `color` varchar(7) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`tagId`),
  UNIQUE KEY `tagId_UNIQUE` (`tagId`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
