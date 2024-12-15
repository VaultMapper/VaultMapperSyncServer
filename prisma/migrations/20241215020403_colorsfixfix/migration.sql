/*
  Warnings:

  - Made the column `playerPlayerUuid` on table `ColorCache` required. This step will fail if there are existing NULL values in that column.

*/
-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_ColorCache" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT DEFAULT 1,
    "color" TEXT NOT NULL,
    "playerPlayerUuid" TEXT NOT NULL,
    CONSTRAINT "ColorCache_playerPlayerUuid_fkey" FOREIGN KEY ("playerPlayerUuid") REFERENCES "Player" ("playerUuid") ON DELETE RESTRICT ON UPDATE CASCADE
);
INSERT INTO "new_ColorCache" ("color", "id", "playerPlayerUuid") SELECT "color", "id", "playerPlayerUuid" FROM "ColorCache";
DROP TABLE "ColorCache";
ALTER TABLE "new_ColorCache" RENAME TO "ColorCache";
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
