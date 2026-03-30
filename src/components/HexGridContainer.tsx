import { useState } from "react";
import { cn } from "../../lib/utils"
import hertaa from "../assets/hertaa.gif"

export default function HexGridContainer() {
  // handle your states here ..yim
  const [data, setData] = useState(false)

  return (
    <div className="w-full max-w-5xl mx-auto px-4 pb-10">
      <div className={cn(
        "relative min-h-[400px] md:min-h-[600px] w-full bg-card",
        "border border-border rounded-3xl overflow-hidden",
        "flex flex-col items-center justify-center p-6 transition-all"
      )}>

        {!data ? (
          <div className="z-10 flex flex-col items-center space-y-4">
            <div className="w-72 h-72 flex items-center justify-center overflow-hidden">
              <img
                src={hertaa}
                alt="Hertaa kawai gif"
                className="w-full h-full object-cover"
              />
            </div>
            <p className="text-muted font-medium text-center max-w-xs">
              Transform your anime and manga history into a beautiful,
              minimalist hexagonal visualization.
            </p>
          </div>
        ) : (
          <div className="z-10 w-full h-full flex items-center justify-center">

            {/* handle your canvas output here.. */}
            <h1 className="text-white font-black text-2xl">HEXAGON CANVAS READY</h1>
          </div>
        )}

      </div>
    </div>
  )
}