/*
  Warnings:

  - You are about to drop the column `color` on the `Player` table. All the data in the column will be lost.

*/
-- CreateTable
CREATE TABLE "ColorCache" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT DEFAULT 1,
    "color" TEXT NOT NULL,
    "playerPlayerUuid" TEXT,
    CONSTRAINT "ColorCache_playerPlayerUuid_fkey" FOREIGN KEY ("playerPlayerUuid") REFERENCES "Player" ("playerUuid") ON DELETE SET NULL ON UPDATE CASCADE
);

-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Player" (
    "playerUuid" TEXT NOT NULL PRIMARY KEY,
    "playerName" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO "new_Player" ("createdAt", "playerName", "playerUuid", "token") SELECT "createdAt", "playerName", "playerUuid", "token" FROM "Player";
DROP TABLE "Player";
ALTER TABLE "new_Player" RENAME TO "Player";
CREATE UNIQUE INDEX "Player_playerUuid_key" ON "Player"("playerUuid");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
