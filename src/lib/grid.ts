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

  let minX = Infinity,
    maxX = -Infinity
  let minY = Infinity,
    maxY = -Infinity

  hexagons.forEach(hex => {
    minX = Math.min(minX, hex.center.x - hex.radius)
    maxX = Math.max(maxX, hex.center.x + hex.radius)
    minY = Math.min(minY, hex.center.y - hex.radius)
    maxY = Math.max(maxY, hex.center.y + hex.radius)
  })

  const totalWidth = maxX - minX
  const totalHeight = maxY - minY

  const padding = 50

  const moveX = -minX + padding
  const moveY = -minY + padding

  hexagons.forEach(hex => hex.move(moveX, moveY))
  center.setPosition(center.x + moveX, center.y + moveY)

  return {
    width: totalWidth + padding * 2,
    height: totalHeight + padding * 2,
    center,
    hexagons,
  }
}
