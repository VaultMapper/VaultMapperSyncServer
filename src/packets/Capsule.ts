import PacketType from "./PacketType.ts";
import JoinPacket from "./JoinPacket.ts";
import LeavePacket from "./LeavePacket.ts";
import MovePacket from "./MovePacket.ts";
import CellPacket from "./CellPacket.ts";

abstract class AbstractCapsule<T> {
  public readonly type: PacketType;
  public readonly data: T;

  constructor(type: PacketType, data: T) {
    this.type = type;
    this.data = data;
  }
}

class JoinPacketCapsule extends AbstractCapsule<JoinPacket> {
  public override type: PacketType = PacketType.JOIN;

  constructor(data: JoinPacket) {
    super(PacketType.JOIN, data);
  }
}

class LeavePacketCapsule extends AbstractCapsule<LeavePacket> {
  public override type: PacketType = PacketType.LEAVE;

  constructor(data: LeavePacket) {
    super(PacketType.LEAVE, data);
  }
}

class MovePacketCapsule extends AbstractCapsule<MovePacket> {
  public override type: PacketType = PacketType.MOVE;

  constructor(data: MovePacket) {
    super(PacketType.MOVE, data);
  }
}

class CellPacketCapsule extends AbstractCapsule<CellPacket> {
  public override type: PacketType = PacketType.CELL;

  constructor(data: CellPacket) {
    super(PacketType.CELL, data);
  }
}

export type Capsule = JoinPacketCapsule | LeavePacketCapsule | MovePacketCapsule | CellPacketCapsule;
export { JoinPacketCapsule, LeavePacketCapsule, MovePacketCapsule, CellPacketCapsule };
