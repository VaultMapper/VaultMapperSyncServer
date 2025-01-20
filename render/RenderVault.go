package render

import (
	"bytes"
	"errors"
	"log"
	"math"

	"github.com/NodiumHosting/VaultMapperSyncServer/icons"
	"github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/fogleman/gg"
	"golang.org/x/image/math/f64"
)

const cellSize = 10 // might have to adjust for icon rendering

func RenderVault(cells []*proto.VaultCell) (error, []byte) {
	minX, minZ, maxX, maxZ := calculateMapResolution(cells)
	horizonatalCells := maxX - minX
	verticalCells := maxZ - minZ
	horizontalRes := int(horizonatalCells*cellSize) + cellSize
	verticalRes := int(verticalCells*cellSize) + cellSize

	if (horizontalRes * verticalRes) > 50_000_000 {
		log.Printf("Map is too large to render(%d x %d)\n", horizontalRes, verticalRes)
		return errors.New("Map is too large to render"), nil
	}

	dc := gg.NewContext(horizontalRes, verticalRes)

	startCell := &proto.VaultCell{
		X:         0,
		Z:         0,
		CellType:  proto.CellType_CELLTYPE_ROOM,
		RoomType:  proto.RoomType_ROOMTYPE_START,
		RoomName:  proto.RoomName_ROOMNAME_UNKNOWN,
		Explored:  true,
		Inscribed: false,
		Marked:    false,
	}
	drawCell(dc, startCell, minX, minZ)

	for _, cell := range cells {
		drawCell(dc, cell, minX, minZ)
	}

	var buf bytes.Buffer
	err := dc.EncodePNG(&buf)
	if err != nil {
		return err, nil
	}

	return nil, buf.Bytes()
}

func getCellColor(cell *proto.VaultCell) f64.Vec3 {
	if cell.RoomType == proto.RoomType_ROOMTYPE_START {
		return f64.Vec3{1, 0, 0}
	}
	if cell.Marked {
		return f64.Vec3{1, 0, 1}
	}
	if cell.Inscribed {
		return f64.Vec3{1, 1, 0}
	}
	if cell.RoomType == proto.RoomType_ROOMTYPE_OMEGA {
		return f64.Vec3{0.3333333333333333, 1, 0.3333333333333333}
	}
	if cell.RoomType == proto.RoomType_ROOMTYPE_CHALLENGE {
		return f64.Vec3{0.9411764705882353, 0.6196078431372549, 0}
	}
	if cell.RoomType == proto.RoomType_ROOMTYPE_ORE {
		return f64.Vec3{0, 1, 1}
	}
	return f64.Vec3{0, 0, 1}
}

func drawCell(dc *gg.Context, cell *proto.VaultCell, minX, minZ int32) {
	portalX := float64(cellSize * (-minX))
	portalZ := float64(cellSize * (-minZ))

	x := float64(cell.X)
	z := float64(cell.Z)

	color := getCellColor(cell)
	dc.SetRGB(color[0], color[1], color[2])

	switch cell.CellType {
	case proto.CellType_CELLTYPE_ROOM:
		dc.DrawRectangle(portalX+x*cellSize, portalZ+z*cellSize, cellSize, cellSize)
		dc.Fill()
	case proto.CellType_CELLTYPE_TUNNEL_X:
		dc.DrawRectangle(portalX+x*cellSize, portalZ+z*cellSize+cellSize/4, cellSize, (cellSize/4)*2+2)
		dc.Fill()
	case proto.CellType_CELLTYPE_TUNNEL_Z:
		dc.DrawRectangle(portalX+x*cellSize+cellSize/4, portalZ+z*cellSize, (cellSize/4)*2+2, cellSize)
		dc.Fill()
	}

	if cell.RoomName == proto.RoomName_ROOMNAME_UNKNOWN {
		return
	}

	icon := icons.GetIcon(&cell.RoomName)
	if icon == nil {
		return
	}

	cellCenterX := int(portalX) + int(x*cellSize+cellSize/2)
	cellCenterZ := int(portalZ) + int(z*cellSize+cellSize/2)

	dc.DrawImageAnchored(icon, cellCenterX, cellCenterZ, 0.5, 0.5)
	dc.Fill()
}

func calculateMapResolution(cells []*proto.VaultCell) (int32, int32, int32, int32) {
	minX := int32(math.MaxInt32)
	minZ := int32(math.MaxInt32)
	maxX := int32(0)
	maxZ := int32(0)
	for _, cell := range cells {
		if cell.X < minX {
			minX = cell.X
		}
		if cell.Z < minZ {
			minZ = cell.Z
		}
		if cell.X > maxX {
			maxX = cell.X
		}
		if cell.Z > maxZ {
			maxZ = cell.Z
		}
	}
	// add 1 to each side to avoid rounded corners
	return minX - 1, minZ - 1, maxX + 1, maxZ + 1
}
