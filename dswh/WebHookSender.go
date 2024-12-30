package dswh

import (
	"bytes"
	"fmt"
	"github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/NodiumHosting/VaultMapperSyncServer/render"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"os"
)

func SendMap(cells []*proto.VaultCell, vaultID string) {
	err, data := render.RenderVault(cells)
	if err != nil {
		return
	}

	sendWebhookImage(data, vaultID)
}

func sendWebhookImage(data []byte, vaultID string) {
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

	builder := discord.NewWebhookMessageCreateBuilder()
	builder.SetContent("Map of vault " + vaultID)
	builder.AddFile(vaultID+"_map.png", "", bytes.NewReader(data))

	_, err = client.CreateMessage(builder.Build())
	if err != nil {
		fmt.Println(err)
		return
	}
}
