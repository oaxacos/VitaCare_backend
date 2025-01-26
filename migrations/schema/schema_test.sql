INSERT INTO "users" ("id", "email", "first_name", "last_name", "rol", "dni", "birthdate", "phone", "is_active", "deceased_at")VALUES
('f3b27c73-c844-476e-917f-56278406579f', 'john.doe@example.com', 'John', 'Doe', 'patient', '123456789', '1990-05-15T00:00:00Z', '+1234567890', true, NULL);

INSERT INTO "tokens" ("id", "token", "user_id", "expired_at", "created_at")VALUES
('f99a7d8f-7add-4047-a0c6-fc2cb803bfeb','KjIZ52lspWj0xBHp7zYQxizcCHGT9w5tRsRRUZ8ViGQ=','f3b27c73-c844-476e-917f-56278406579f'	,'2020-01-26 14:13:41.466 -0600', CURRENT_TIMESTAMP);