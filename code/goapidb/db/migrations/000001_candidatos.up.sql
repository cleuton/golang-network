CREATE TABLE IF NOT EXISTS candidatos (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  nome TEXT not null,
  created_at TIMESTAMP default now()
  );