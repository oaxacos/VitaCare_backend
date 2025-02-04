-- migrate:up
-- SEE https://dbdiagram.io/d/vita-care-67959f43263d6cf9a010b168
DO $$ BEGIN
CREATE TYPE "rol" AS ENUM (
  'patient',
  'doctor',
  'admin',
  'secretary'
);
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "email" text UNIQUE NOT NULL,
  "first_name" text,
  "last_name" text,
  "rol" rol DEFAULT 'patient',
  "dni" text UNIQUE NOT NULL,
  "birthdate" timestamptz,
  "phone" text,
  "is_active" bool DEFAULT true,
  "deceased_at" timestamptz DEFAULT null,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "update_at" timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "address" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid,
  "address_line_1" text NOT NULL,
  "address_line_2" text,
  "zip_code" text,
  "country" text,
  "city" text,
  "state" text
);

CREATE TABLE "medical_ensure" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid,
  "name" text,
  "institution" text,
  "social_security_number" text NOT NULL
);

CREATE TABLE "doctor" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid,
  "salary" decimal
);

CREATE TABLE "doctor_speciality" (
  "id" uuid PRIMARY KEY,
  "doctor_id" uuid,
  "speciality_id" uuid
);

CREATE TABLE "medical_specialties" (
  "id" uuid PRIMARY KEY,
  "name" text,
  "description" text,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "update_at" timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "schedule" (
  "id" uuid PRIMARY KEY,
  "doctor_id" uuid,
  "start_at" timestamptz NOT NULL,
  "end_at" timestamptz NOT NULL,
  "is_active" bool DEFAULT true,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "update_at" timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "services" (
  "id" uuid PRIMARY KEY,
  "name" text,
  "description" text,
  "price" decimal NOT NULL,
  "created_by" uuid,
  "deleted_by" uuid,
  "is_active" bool DEFAULT true,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "update_at" timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "packages" (
  "id" uuid PRIMARY KEY,
  "code" text NOT NULL,
  "price" decimal NOT NULL,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "update_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "is_active" bool DEFAULT true,
  "created_by" uuid,
  "deleted_by" uuid
);

CREATE TABLE "services_in_packages" (
  "id" uuid PRIMARY KEY,
  "service_id" uuid,
  "package_id" uuid
);

CREATE TABLE "medical_appointments" (
  "id" uuid PRIMARY KEY,
  "date" timestamptz,
  "patient_id" uuid,
  "doctor_id" uuid,
  "service_id" uuid DEFAULT null,
  "package_id" uuid DEFAULT null,
  "total" decimal,
  "subtotal" decimal,
  "payment_at" timestamptz,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "update_at" timestamptz DEFAULT CURRENT_TIMESTAMP
);


-- migrate:down
DROP TABLE IF EXISTS "users" cascade;
DROP TYPE IF EXISTS "rol";
DROP TABLE IF EXISTS "address";
DROP TABLE IF EXISTS "medical_ensure";
DROP TABLE IF EXISTS "doctor";
DROP TABLE IF EXISTS "doctor_speciality";
DROP TABLE IF EXISTS "medical_specialties";
DROP TABLE IF EXISTS "schedule";
DROP TABLE IF EXISTS "services";
DROP TABLE IF EXISTS "packages";
DROP TABLE IF EXISTS "services_in_packages";
DROP TABLE IF EXISTS "medical_appointments";

