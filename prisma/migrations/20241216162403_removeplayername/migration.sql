/*
  Warnings:

  - You are about to drop the column `playerName` on the `Player` table. All the data in the column will be lost.

*/
-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Player" (
    "playerUuid" TEXT NOT NULL PRIMARY KEY,
    "color" TEXT NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO "new_Player" ("color", "createdAt", "playerUuid") SELECT "color", "createdAt", "playerUuid" FROM "Player";
DROP TABLE "Player";
ALTER TABLE "new_Player" RENAME TO "Player";
CREATE UNIQUE INDEX "Player_playerUuid_key" ON "Player"("playerUuid");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
