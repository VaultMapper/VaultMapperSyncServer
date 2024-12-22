package render

import (
	"bytes"
	"fmt"
	"github.com/NodiumHosting/VaultMapperSyncServer/icons"
	"github.com/NodiumHosting/VaultMapperSyncServer/proto"
	"github.com/fogleman/gg"
	"golang.org/x/image/math/f64"
	"math"
)

const cellSize = 10 // might have to adjust for icon rendering

func RenderVault(cells []*proto.VaultCell) (error, []byte) {
	mapRes := calculateMapResolution(cells)
	cellsPerSide := mapRes*2 + 1
	res := cellSize * cellsPerSide
	dc := gg.NewContext(res+cellSize, res+cellSize)

	for _, cell := range cells {
		drawCell(dc, cell, res)
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
	return f64.Vec3{0, 0, 1}
}

func drawCell(dc *gg.Context, cell *proto.VaultCell, res int) {
	middle := float64(res / 2)

	x := float64(cell.X)
	z := float64(cell.Z)

	fmt.Println("Drawing cell at x:", middle+x*cellSize, "z:", middle+z*cellSize)
	fmt.Println("middle ", middle, "cellSize ", cellSize)

	color := getCellColor(cell)
	dc.SetRGB(color[0], color[1], color[2])

	switch cell.CellType {
	case proto.CellType_CELLTYPE_ROOM:
		dc.DrawRectangle(middle+x*cellSize, middle+z*cellSize, cellSize, cellSize)
		dc.Fill()
	case proto.CellType_CELLTYPE_TUNNEL_X:
		dc.DrawRectangle(middle+x*cellSize, middle+z*cellSize+cellSize/4, cellSize, (cellSize/4)*2+2)
		dc.Fill()
	case proto.CellType_CELLTYPE_TUNNEL_Z:
		dc.DrawRectangle(middle+x*cellSize+cellSize/4, middle+z*cellSize, (cellSize/4)*2+2, cellSize)
		dc.Fill()
	}

	if cell.RoomName == proto.RoomName_ROOMNAME_UNKNOWN {
		return
	}

	icon := icons.GetIcon(&cell.RoomName)
	if icon == nil {
		return
	}

	cellCenterX := int(middle) + int(x*cellSize+cellSize/2)
	cellCenterZ := int(middle) + int(z*cellSize+cellSize/2)

	dc.DrawImageAnchored(icon, cellCenterX, cellCenterZ, 0.5, 0.5)
	dc.Fill()
}

func calculateMapResolution(cells []*proto.VaultCell) int {
	maxX := int32(0)
	maxZ := int32(0)
	for _, cell := range cells {
		if math.Abs(float64(cell.X)) > float64(maxX) {
			maxX = int32(math.Abs(float64(cell.X)))
		}
		if math.Abs(float64(cell.Z)) > float64(maxZ) {
			maxZ = int32(math.Abs(float64(cell.Z)))
		}
	}
	return int(math.Max(float64(maxX), float64(maxZ)))
}
