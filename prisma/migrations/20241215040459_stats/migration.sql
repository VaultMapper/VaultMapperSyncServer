/*
  Warnings:

  - You are about to drop the column `vaultsRan` on the `Stats` table. All the data in the column will be lost.
  - Added the required column `stat` to the `Stats` table without a default value. This is not possible if the table is not empty.
  - Added the required column `value` to the `Stats` table without a default value. This is not possible if the table is not empty.

*/
-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Stats" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "stat" TEXT NOT NULL,
    "value" INTEGER NOT NULL
);
INSERT INTO "new_Stats" ("id") SELECT "id" FROM "Stats";
DROP TABLE "Stats";
ALTER TABLE "new_Stats" RENAME TO "Stats";
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
