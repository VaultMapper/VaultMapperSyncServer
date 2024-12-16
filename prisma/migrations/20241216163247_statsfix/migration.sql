/*
  Warnings:

  - The primary key for the `Stats` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `id` on the `Stats` table. All the data in the column will be lost.

*/
-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Stats" (
    "stat" TEXT NOT NULL PRIMARY KEY,
    "value" INTEGER NOT NULL
);
INSERT INTO "new_Stats" ("stat", "value") SELECT "stat", "value" FROM "Stats";
DROP TABLE "Stats";
ALTER TABLE "new_Stats" RENAME TO "Stats";
CREATE UNIQUE INDEX "Stats_stat_key" ON "Stats"("stat");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
