export default class Point {
  x: number;
  y: number;
  constructor(x: number, y: number) {
    this.x = x;
    this.y = y;
  }

  move(x: number, y: number) {
    this.x += x;
    this.y += y;
  }

  setPosition(x: number, y: number) {
    this.x = x;
    this.y = y;
  }

  rotate(base: Point, angle: number) {
    const cosTheta = Math.cos(angle);
    const sinTheta = Math.sin(angle);
    const dx = this.x - base.x;
    const dy = this.y - base.y;

    this.x = cosTheta * dx - sinTheta * dy + base.x;
    this.y = sinTheta * dx + cosTheta * dy + base.y;
  }

  distance(d: Point): number {
    return Math.sqrt(Math.pow(this.x - d.x, 2) + Math.pow(this.y - d.y, 2));
  }
}
