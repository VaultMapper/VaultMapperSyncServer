/*
  Warnings:

  - A unique constraint covering the columns `[playerPlayerUuid]` on the table `ColorCache` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "ColorCache_playerPlayerUuid_key" ON "ColorCache"("playerPlayerUuid");
