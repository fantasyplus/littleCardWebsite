```
CREATE TABLE IF NOT EXISTS personInfo (
  person_id INT PRIMARY KEY,
  cn VARCHAR(255),
  qq VARCHAR(255),
  UNIQUE KEY unique_cn_qq (cn, qq)
);

```

```
CREATE TABLE IF NOT EXISTS cardInfo (
  card_id VARCHAR(255) PRIMARY KEY,
  card_name VARCHAR(255),
  card_character VARCHAR(255),
  card_type VARCHAR(255),
  card_condition VARCHAR(1024),
  other VARCHAR(1024)
);
```

```
CREATE TABLE IF NOT EXISTS cardIndex (
  person_id INT,
  card_ids VARCHAR(8192),
  FOREIGN KEY (person_id) REFERENCES personInfo(person_id)
);
```

```
CREATE TABLE cardNo{} (
  person_id INT,
  card_name VARCHAR(255),
  card_num INT,
  status VARCHAR(255),
  FOREIGN KEY (person_id) REFERENCES personInfo(person_id)
);
```



rsync -avr --delete dist/ ubuntu@myweb:/home/ubuntu/littleCardWebsite/front/dist/