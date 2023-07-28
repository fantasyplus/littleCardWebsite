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
  card_condition VARCHAR(255),
  other VARCHAR(255)
);
```

```
CREATE TABLE IF NOT EXISTS cardIndex (
  person_id INT,
  card_ids VARCHAR(1024),
  FOREIGN KEY (person_id) REFERENCES personInfo(person_id)
);
```

```
CREATE TABLE cardNo{} (
  person_id INT,
  card_name VARCHAR(255),
  card_num INT,
  FOREIGN KEY (person_id) REFERENCES personInfo(person_id)
);
```

```
use non_commercial;
delete from cardno1_1;
delete from cardno1_2;
delete from cardno2;
delete from cardno3;
delete from cardno4;
delete from cardno5;
delete from cardno6;
delete from cardno7;
delete from cardno8;
delete from cardno9;
delete from cardno10;
delete from cardno11;
delete from cardno12;
delete from cardno13;
delete from cardno14;


delete from cardno17;
delete from cardno18;
delete from cardno19_1;
delete from cardno19_2;
delete from cardno19_3;
delete from cardno20;
delete from cardno21;
delete from cardno22;
delete from cardno23;
delete from cardno24;
delete from cardno25;
delete from cardno26;
delete from cardno27;

delete from cardno29;
delete from cardno30;
delete from cardno31;
delete from cardno400;
delete from cardindex;
delete from cardinfo;
delete from personinfo;
```