-- migrate:up

ALTER TABLE IF EXISTS "doctor" ADD CONSTRAINT "fk_doctor_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE IF EXISTS "doctor_speciality" ADD CONSTRAINT "fk_doctor_id" FOREIGN KEY ("doctor_id") REFERENCES "doctor" ("id");
ALTER TABLE IF EXISTS "doctor_speciality" ADD CONSTRAINT "fk_speciality_id" FOREIGN KEY ("speciality_id") REFERENCES "medical_specialties" ("id");
ALTER TABLE IF EXISTS "schedule" ADD CONSTRAINT "fk_doctor_id" FOREIGN KEY ("doctor_id") REFERENCES "doctor" ("id");
ALTER TABLE IF EXISTS "medical_ensure" ADD CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE IF EXISTS "address" ADD CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE IF EXISTS "services_in_packages" ADD CONSTRAINT "fk_package_id" FOREIGN KEY ("package_id") REFERENCES "packages" ("id");
ALTER TABLE IF EXISTS "services_in_packages" ADD CONSTRAINT "fk_service_id" FOREIGN KEY ("service_id") REFERENCES "services" ("id");
ALTER TABLE IF EXISTS "medical_appointments" ADD CONSTRAINT "fk_patient_id" FOREIGN KEY ("patient_id") REFERENCES "users" ("id");
ALTER TABLE IF EXISTS "medical_appointments" ADD CONSTRAINT "fk_doctor_id" FOREIGN KEY ("doctor_id") REFERENCES "users" ("id");
ALTER TABLE IF EXISTS "medical_appointments" ADD CONSTRAINT "fk_service_id" FOREIGN KEY ("service_id") REFERENCES "services" ("id");
ALTER TABLE IF EXISTS "medical_appointments" ADD CONSTRAINT "fk_package_id" FOREIGN KEY ("package_id") REFERENCES "packages" ("id");

-- migrate:down
ALTER TABLE IF EXISTS "doctor" DROP CONSTRAINT "fk_doctor_id";
ALTER TABLE IF EXISTS "doctor_speciality" DROP CONSTRAINT "fk_doctor_id";
ALTER TABLE IF EXISTS "doctor_speciality" DROP CONSTRAINT "fk_speciality_id";
ALTER TABLE IF EXISTS "schedule" DROP CONSTRAINT "fk_doctor_id";
ALTER TABLE IF EXISTS "medical_ensure" DROP CONSTRAINT "fk_user_id";
ALTER TABLE IF EXISTS "address" DROP CONSTRAINT "fk_user_id";
ALTER TABLE IF EXISTS "services_in_packages" DROP CONSTRAINT "fk_package_id";
ALTER TABLE IF EXISTS "services_in_packages" DROP CONSTRAINT "fk_service_id";
ALTER TABLE IF EXISTS "medical_appointments" DROP CONSTRAINT "fk_patient_id";
ALTER TABLE IF EXISTS "medical_appointments" DROP CONSTRAINT "fk_doctor_id";
ALTER TABLE IF EXISTS "medical_appointments" DROP CONSTRAINT "fk_service_id";
ALTER TABLE IF EXISTS "medical_appointments" DROP CONSTRAINT "fk_package_id";



