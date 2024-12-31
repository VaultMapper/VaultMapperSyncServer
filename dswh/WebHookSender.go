package dswh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NodiumHosting/VaultMapperSyncServer/models"
	"github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/NodiumHosting/VaultMapperSyncServer/render"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
)

func SendMap(cells []*proto.VaultCell, vaultID string, DB *gorm.DB) {
	err, data := render.RenderVault(cells)
	if err != nil {
		return
	}

	sendWebhookImage(data, vaultID, DB)
}

type MojangResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Properties []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"properties"`
	ProfileActions []interface{} `json:"profileActions"`
}

func GetMinecraftUsername(uuid string) string {
	url := fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", uuid)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var names MojangResponse
	if err := json.NewDecoder(resp.Body).Decode(&names); err != nil {
		return ""
	}

	return names.Name
}

func sendWebhookImage(data []byte, vaultID string, DB *gorm.DB) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK")
	if webhookURL == "" {
		fmt.Println("DISCORD_WEBHOOK not set")
		return
	}

	client, err := webhook.NewWithURL(webhookURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	var players []string
	var playerRecords []models.PlayerVault
	result := DB.Find(&playerRecords, "vault_id = ?", vaultID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("Error retrieving players:", result.Error)
		return
	}

	for _, player := range playerRecords {
		players = append(players, GetMinecraftUsername(player.PlayerUUID)+" ("+player.PlayerUUID+")\n")
	}

	playersList := strings.Join(players, ", ")

	builder := discord.NewWebhookMessageCreateBuilder()
	builder.SetContent("Map of vault " + vaultID + "\nPlayers: " + playersList)
	builder.AddFile(vaultID+"_map.png", "", bytes.NewReader(data))

	_, err = client.CreateMessage(builder.Build())
	if err != nil {
		fmt.Println(err)
		return
	}
}
