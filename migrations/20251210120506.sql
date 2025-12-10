-- Rename a column from "executable_tag" to "executable_id"
ALTER TABLE "public"."parameters" RENAME COLUMN "executable_tag" TO "executable_id";
