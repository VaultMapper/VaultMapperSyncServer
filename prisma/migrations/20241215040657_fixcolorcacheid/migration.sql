-- RedefineTables
PRAGMA defer_foreign_keys=ON;
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_ColorCache" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "color" TEXT NOT NULL,
    "playerPlayerUuid" TEXT NOT NULL,
    CONSTRAINT "ColorCache_playerPlayerUuid_fkey" FOREIGN KEY ("playerPlayerUuid") REFERENCES "Player" ("playerUuid") ON DELETE RESTRICT ON UPDATE CASCADE
);
INSERT INTO "new_ColorCache" ("color", "id", "playerPlayerUuid") SELECT "color", "id", "playerPlayerUuid" FROM "ColorCache";
DROP TABLE "ColorCache";
ALTER TABLE "new_ColorCache" RENAME TO "ColorCache";
CREATE UNIQUE INDEX "ColorCache_playerPlayerUuid_key" ON "ColorCache"("playerPlayerUuid");
PRAGMA foreign_keys=ON;
PRAGMA defer_foreign_keys=OFF;
