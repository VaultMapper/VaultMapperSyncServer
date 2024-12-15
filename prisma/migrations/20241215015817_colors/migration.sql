/*
  Warnings:

  - Added the required column `color` to the `Player` table without a default value. This is not possible if the table is not empty.

*/
-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Player" (
    "playerUuid" TEXT NOT NULL PRIMARY KEY,
    "playerName" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "color" TEXT NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO "new_Player" ("createdAt", "playerName", "playerUuid", "token") SELECT "createdAt", "playerName", "playerUuid", "token" FROM "Player";
DROP TABLE "Player";
ALTER TABLE "new_Player" RENAME TO "Player";
CREATE UNIQUE INDEX "Player_playerUuid_key" ON "Player"("playerUuid");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
