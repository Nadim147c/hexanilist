import { useState } from "react"
import Header from "./components/Header"
import Searchbar from "./components/Searchbar"
import HexGridContainer from "./components/HexGridContainer"
import Footer from "./components/Footer"
import { createNodes, fetchAnilist } from "./lib/anilist"
import type { AnilistData } from "./types/anilist"
import { generateHexagonRing } from "./lib/grid"

function getCrop(
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

interface AnilistCache {
  timestamp: number
  data: AnilistData
}

interface Options {
  strokeColor?: string
  strokeWidth?: number
  hexagonRadius?: number
}

async function generate(username: string, opts: Options = {}) {
  let data: AnilistData
  let loadedFromCache = false
  // load from cache if available
  const cached = localStorage.getItem(username)
  if (cached !== null) {
    const cache: AnilistCache = JSON.parse(cached)
    if (Date.now() - cache.timestamp < 10 * 60 * 1000) {
      data = cache.data
      loadedFromCache = true
    }
  }

  if (!loadedFromCache) {
    const [list, err] = await fetchAnilist(username)
    if (err !== null) throw err
    data = list
    // cache the list for 10 minutes
    const cache: AnilistCache = {
      timestamp: Date.now(),
      data: data,
    }
    localStorage.setItem(username, JSON.stringify(cache))
  }

  const nodes = createNodes(data!)
  const hex = generateHexagonRing(nodes.length, opts.hexagonRadius ?? 50)

  if (nodes.length !== hex.hexagons.length) {
    throw new Error(
      `Number of nodes does not match number of hexagons. nodes=${nodes.length}, hexagons=${hex.hexagons.length}`
    )
  }

  const canvas = document.createElement("canvas") as HTMLCanvasElement
  canvas.width = hex.width
  canvas.height = hex.height

  console.debug({ width: hex.width, height: hex.height })

  const ctx = canvas.getContext("2d")
  if (!ctx) throw new Error("Could not get canvas context")

  const imagePromises = nodes.map(node => {
    return new Promise<HTMLImageElement>((resolve, reject) => {
      const img = new Image()
      img.crossOrigin = "anonymous"
      img.onload = () => resolve(img)
      img.onerror = () => reject(new Error(`Failed to load: ${node.image}`))
      img.src = node.image
    })
  })

  const loadedImages = await Promise.all(imagePromises)

  hex.hexagons.forEach((hexagon, index) => {
    const img = loadedImages[index]
    ctx.save()
    hexagon.draw(ctx)
    ctx.clip()

    const hexBox = hexagon.box()
    const crop = getCrop(hexBox.w, hexBox.h, img.width, img.height)
    ctx.drawImage(
      img,
      crop.x,
      crop.y,
      crop.w,
      crop.h,
      hexBox.x,
      hexBox.y,
      hexBox.w,
      hexBox.h
    )

    ctx.restore()

    if (opts.strokeWidth) {
      hexagon.draw(ctx)
      ctx.strokeStyle = opts.strokeColor || "black"
      ctx.lineWidth = opts.strokeWidth
      ctx.stroke()
      ctx.restore()
    }
  })

  // Optional: return the data URL if you need to save it
  return canvas.toDataURL("image/png")
}

export default function App() {
  const [loading, setLoading] = useState<boolean>(false)
  const [result, setResult] = useState<string | null>(null)
  const [username, setUsername] = useState<string | null>(null)
  const [error, setError] = useState<boolean>(false)

  const handleSearch = async (username: string) => {
    setLoading(true)
    setResult(null)
    setError(false)
    setUsername(username)
    try {
      const url = await generate(username, {
        hexagonRadius: 50,
        strokeWidth: 4,
        strokeColor: "black", // or use any hex color
      })
      setResult(url)

      console.log("Searching for:", username)
    } catch (err) {
      console.error("Error fetching data:", err)
      setError(true)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <Header />

      <main className="flex flex-col items-center justify-center space-y-6">
        <h2 className="text-center text-3xl font-black tracking-tighter text-white uppercase md:text-4xl">
          Your AniList Journey,{" "}
          <span className="text-slate-500">Hexagonal</span>
        </h2>
        <Searchbar onSearch={handleSearch} loading={loading} />

        <HexGridContainer
          imageResult={result}
          username={username}
          error={error}
        />
      </main>

      <Footer />
    </>
  )
}
