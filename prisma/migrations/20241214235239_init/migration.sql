-- CreateTable
CREATE TABLE "Player" (
    "playerUuid" TEXT NOT NULL PRIMARY KEY,
    "playerName" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CreateIndex
CREATE UNIQUE INDEX "Player_playerUuid_key" ON "Player"("playerUuid");
