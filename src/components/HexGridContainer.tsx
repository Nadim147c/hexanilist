import { useState } from "react"
import { cn } from "../lib/utils"
import hertaa from "../assets/hertaa.gif"

export default function HexGridContainer() {
  // handle your states here ..yim
  const [data, setData] = useState(false)

  return (
    <div className="mx-auto w-full max-w-5xl px-4 pb-10">
      <div
        className={cn(
          "bg-card relative min-h-[400px] w-full md:min-h-[600px]",
          "border-border overflow-hidden rounded-3xl border",
          "flex flex-col items-center justify-center p-6 transition-all"
        )}
      >
        {!data ? (
          <div className="z-10 flex flex-col items-center space-y-4">
            <div className="flex h-72 w-72 items-center justify-center overflow-hidden">
              <img
                src={hertaa}
                alt="Hertaa kawai gif"
                className="h-full w-full object-cover"
              />
            </div>
            <p className="text-muted max-w-xs text-center font-medium">
              Transform your anime and manga history into a beautiful,
              minimalist hexagonal visualization.
            </p>
          </div>
        ) : (
          <div className="z-10 flex h-full w-full items-center justify-center">
            {/* handle your canvas output here.. */}
            <h1 className="text-2xl font-black text-white">
              HEXAGON CANVAS READY
            </h1>
          </div>
        )}
      </div>
    </div>
  )
}
