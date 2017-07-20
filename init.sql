
CREATE TABLE players (
  player_id TEXT NOT NULL PRIMARY KEY,
  balance   NUMERIC(15, 4) NOT NULL DEFAULT 0
);

CREATE TABLE tournaments (
  tournament_id TEXT NOT NULL PRIMARY KEY,
  deposit       NUMERIC(15, 4) NOT NULL,
  status        SMALLINT NOT NULL DEFAULT 1
);

CREATE TABLE participants (
  tournament_id TEXT NOT NULL,
  player_id     TEXT NOT NULL,
  backer_id     TEXT NOT NULL
);

CREATE UNIQUE INDEX participants_unq ON participants (tournament_id, player_id, backer_id);