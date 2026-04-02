export function getCrop(
  boxWidth: number,
  boxHeight: number,
  imgWidth: number,
  imgHeight: number
) {
  const boxRatio = boxWidth / boxHeight
  const imgRatio = imgWidth / imgHeight

  if (boxRatio > imgRatio) {
    const w = imgWidth
    const h = imgWidth / boxRatio
    const x = 0
    const y = (imgHeight - h) / 2
    return { x, y, w, h }
  } else {
    const w = imgHeight * boxRatio
    const h = imgHeight
    const x = (imgWidth - w) / 2
    const y = 0
    return { x, y, w, h }
  }
}
