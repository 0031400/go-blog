PRAGMA foreign_keys = ON;
DROP TABLE IF EXISTS post_tags;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
CREATE TABLE categories (
  uuid VARCHAR(32) PRIMARY KEY,
  name VARCHAR(32) NOT NULL
);
CREATE TABLE posts (
  uuid VARCHAR(32) PRIMARY KEY,
  title VARCHAR(64) NOT NULL,
  content TEXT NOT NULL,
  date VARCHAR(8) NOT NULL,
  brief VARCHAR(200) NOT NULL,
  categoryUUID VARCHAR(32) NOT NULL,
  FOREIGN KEY (categoryUUID) REFERENCES categories(uuid) 
);
CREATE TABLE tags (
  uuid VARCHAR(32) PRIMARY KEY,
  name VARCHAR(32) NOT NULL
);
CREATE TABLE post_tags (
  postUUID VARCHAR(32),
  tagUUID VARCHAR(32),
  PRIMARY KEY (postUUID,tagUUID),
  FOREIGN KEY (postUUID) REFERENCES posts(uuid),
  FOREIGN KEY (tagUUID) REFERENCES tags(uuid)
);

INSERT INTO categories VALUES ('6a70f1d6cf9f4c9b8d7ca36732bbcef7','category one');
INSERT INTO categories VALUES ('20358c152d5c4f3bbe10bd3929a79d1c','category two');

INSERT INTO posts VALUES ('5c9bb8c13fe8420587931d8777a27696','title one','content one','20250608','brief one','6a70f1d6cf9f4c9b8d7ca36732bbcef7');
INSERT INTO posts VALUES ('9658f88584dc435fb08522cd7af9bb53','title two','content two','20250608','brief two','6a70f1d6cf9f4c9b8d7ca36732bbcef7');
INSERT INTO posts VALUES ('90595cb8b6934569b88af06cf7752c7c','title three','content three','20250608','brief three','6a70f1d6cf9f4c9b8d7ca36732bbcef7');
INSERT INTO posts VALUES ('3aeca46a09b34b2485a2e80606a24a32','title four','content four','20250608','brief four','20358c152d5c4f3bbe10bd3929a79d1c');
INSERT INTO posts VALUES ('88eb49f836fb4d9dba0156f7ec24cf8e','title five','content five','20250608','brief five','20358c152d5c4f3bbe10bd3929a79d1c');

INSERT INTO tags VALUES ('53483b8f534a4ce2a566f3db9fd07ea2','tag one');
INSERT INTO tags VALUES ('8148aa3523b44069a7057e3fa40dafba','tag two');

INSERT INTO post_tags VALUES ('5c9bb8c13fe8420587931d8777a27696','53483b8f534a4ce2a566f3db9fd07ea2');
INSERT INTO post_tags VALUES ('9658f88584dc435fb08522cd7af9bb53','53483b8f534a4ce2a566f3db9fd07ea2');
INSERT INTO post_tags VALUES ('90595cb8b6934569b88af06cf7752c7c','53483b8f534a4ce2a566f3db9fd07ea2');
INSERT INTO post_tags VALUES ('3aeca46a09b34b2485a2e80606a24a32','8148aa3523b44069a7057e3fa40dafba');
INSERT INTO post_tags VALUES ('88eb49f836fb4d9dba0156f7ec24cf8e','8148aa3523b44069a7057e3fa40dafba');
INSERT INTO post_tags VALUES ('90595cb8b6934569b88af06cf7752c7c','8148aa3523b44069a7057e3fa40dafba');
INSERT INTO post_tags VALUES ('3aeca46a09b34b2485a2e80606a24a32','53483b8f534a4ce2a566f3db9fd07ea2');
