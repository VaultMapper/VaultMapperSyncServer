enum CellType {
  NONE,
  ROOM,
  TUNNEL_X,
  TUNNEL_Z,
}

function parseCellType(cellType: string): CellType {
  switch (cellType) {
    case "0":
      return CellType.NONE;
    case "1":
      return CellType.ROOM;
    case "2":
      return CellType.TUNNEL_X;
    case "3":
      return CellType.TUNNEL_Z;
    default:
      throw new Error(`Unknown cell type: ${cellType}`);
  }
}

export default CellType;
export { parseCellType };
