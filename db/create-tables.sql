DROP TABLE IF EXISTS album;
CREATE TABLE album(
    id         INT AUTO_INCREMENT NOT NULL,
    title      VARCHAR(255) NOT NULL,
    artist     VARCHAR(255) NOT NULL,
    price      DECIMAL(5,2) NOT NULL,
    genre      VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);