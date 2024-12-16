/*
  Warnings:

  - You are about to drop the `ColorCache` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the column `token` on the `Player` table. All the data in the column will be lost.
  - Added the required column `color` to the `Player` table without a default value. This is not possible if the table is not empty.

*/
-- DropIndex
DROP INDEX "ColorCache_playerPlayerUuid_key";

-- DropTable
PRAGMA foreign_keys=off;
DROP TABLE "ColorCache";
PRAGMA foreign_keys=on;

-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Player" (
    "playerUuid" TEXT NOT NULL PRIMARY KEY,
    "playerName" TEXT NOT NULL,
    "color" TEXT NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO "new_Player" ("createdAt", "playerName", "playerUuid") SELECT "createdAt", "playerName", "playerUuid" FROM "Player";
DROP TABLE "Player";
ALTER TABLE "new_Player" RENAME TO "Player";
CREATE UNIQUE INDEX "Player_playerUuid_key" ON "Player"("playerUuid");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
