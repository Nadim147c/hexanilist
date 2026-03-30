import Hexagon from "./hexagon"
import Point from "./point"

export default class Grid {
  occupied: Map<string, boolean>
  radius: number
  constructor(radius: number) {
    this.occupied = new Map<string, boolean>()
    this.radius = radius
  }

  key(x: number, y: number): string {
    const sx = Math.round(x / this.radius)
    const sy = Math.round(y / this.radius)
    return `${sx},${sy}`
  }

  isOccupied(p: Point): boolean {
    return this.occupied.get(this.key(p.x, p.y)) || false
  }

  markOccupied(p: Point) {
    this.occupied.set(this.key(p.x, p.y), true)
  }
}

export interface HexagonRing {
  height: number
  width: number
  center: Point
  hexagons: Hexagon[]
}

export function generateHexagonRing(
  totalHexagon: number,
  radius: number
): HexagonRing {
  if (totalHexagon === 0) {
    return { height: 0, width: 0, center: new Point(0, 0), hexagons: [] }
  }

  const center = new Point(0, 0)

  const full = new Grid(radius)
  const empty = new Map<Point, boolean>()

  const centerHex = new Hexagon(center, radius, 0)
  full.markOccupied(centerHex.center)
  const hexagons = [centerHex]

  for (const p of centerHex.neighbors()) {
    empty.set(p, true)
  }

  for (let i = 0; i < totalHexagon - 1; i++) {
    let d = Infinity
    let hexCenter = centerHex.center
    for (const point of empty.keys()) {
      const newDist = point.distance(centerHex.center)
      if (newDist < d) {
        hexCenter = point
        d = newDist
      }
    }

    const hex = new Hexagon(hexCenter, radius, 0)
    hexagons.push(hex)

    full.markOccupied(hexCenter)
    for (const p of hex.neighbors()) {
      if (!full.isOccupied(p)) {
        empty.set(p, true)
      }
    }

    for (const p of empty.keys()) {
      if (full.isOccupied(p)) {
        empty.delete(p)
      }
    }
  }

  hexagons.sort((i, j) => {
    return i.center.distance(center) - j.center.distance(center)
  })

  const x = hexagons.reduce((max, hex) => {
    return Math.max(max, Math.abs(hex.box().x + hex.box().w))
  }, 0)

  const y = hexagons.reduce((max, hex) => {
    return Math.max(max, Math.abs(hex.box().y + hex.box().h))
  }, 0)

  const padding = 30
  const moveX = x / 2 + padding
  const moveY = y / 2 + padding

  center.setPosition(moveX, moveY)
  hexagons.forEach(hex => hex.move(moveX, moveY))

  return {
    height: y + padding * 2,
    width: x + padding * 2,
    center,
    hexagons,
  }
}
