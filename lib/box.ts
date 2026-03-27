export default class Box {
  x: number
  y: number
  w: number
  h: number
  constructor(x: number, y: number, w: number, h: number) {
    this.x = x
    this.y = y
    this.w = w
    this.h = h
  }

  start(): [number, number] {
    return [this.x, this.y]
  }

  end(): [number, number] {
    return [this.x + this.w, this.y + this.h]
  }

  size(): [number, number] {
    return [this.w, this.h]
  }
}
