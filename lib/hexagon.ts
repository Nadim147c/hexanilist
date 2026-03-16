import Point from "./point"
import Box from "./box"

export default class Hexagon {
  center: Point
  radius: number
  angle: number
  points: Point[]

  constructor(center: Point, radius: number, angle: number) {
    this.center = center
    this.radius = radius
    this.angle = angle
    this.points = new Array<Point>(6)

    for (let i = 0; i < 6; i++) {
      const theta = angle + i * (Math.PI / 3)
      const p = new Point(center.x + radius, center.y)
      p.rotate(center, theta)
      this.points[i] = p
    }
  }

  move(x: number, y: number) {
    this.center.move(x, y)
    for (const point of this.points) {
      point.move(x, y)
    }
  }

  draw(ctx: CanvasRenderingContext2D) {
    ctx.beginPath()
    ctx.moveTo(this.points[0].x, this.points[0].y)
    for (let i = 0; i < 5; i++) {
      const idx = i + 1
      ctx.lineTo(this.points[idx].x, this.points[idx].y)
    }
    ctx.closePath()
  }

  sideLength(): number {
    return this.points[0].distance(this.points[1])
  }

  neighbors(): Point[] {
    const ringRadius = this.radius * Math.cos(Math.PI / 6) * 2
    const n = new Array<Point>(6)
    for (let i = 0; i < 6; i++) {
      const angle = this.angle + Math.PI / 6 + i * (Math.PI / 3)
      const point = new Point(this.center.x + ringRadius, this.center.y)
      point.rotate(this.center, angle)
      n[i] = point
    }
    return n
  }

  box(): Box {
    const xSlice = this.points.map(p => p.x)
    const ySlice = this.points.map(p => p.y)
    const start = new Point(Math.min(...xSlice), Math.min(...ySlice))
    const end = new Point(Math.max(...xSlice), Math.max(...ySlice))
    return new Box(start.x, start.y, end.x - start.x, end.y - start.y)
  }
}
