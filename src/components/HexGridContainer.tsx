import { cn } from "../lib/utils"
import { Download, ZoomIn, ZoomOut } from "lucide-react"
import { useState, useEffect } from "react"
import hertaa from "../assets/hertaa.gif"
import crying from "../assets/crying.gif"

type imageProps = {
  imageResult: string | null
  username: string | null
  error?: boolean
  className?: string
}

type StateDisplayProps = {
  gif: string
  alt: string
  text: string
  size?: string
  textClass?: string
}

function StateDisplay({
  gif,
  alt,
  text,
  size = "h-72 w-72",
  textClass = "",
}: StateDisplayProps) {
  return (
    <div className="z-10 flex flex-col items-center space-y-4">
      <div
        className={cn("flex items-center justify-center overflow-hidden", size)}
      >
        <img src={gif} alt={alt} className="h-full w-full object-cover" />
      </div>
      <p className={cn("max-w-xs text-center font-medium", textClass)}>
        {text}
      </p>
    </div>
  )
}

export default function HexGridContainer({
  imageResult,
  username,
  error,
  className,
}: imageProps) {
  const [isExpanded, setIsExpanded] = useState<boolean>(false)
  useEffect(() => {
    setIsExpanded(false)
  }, [imageResult])

  const iconButton =
    "flex h-10 w-10 items-center justify-center rounded-full transition-all hover:scale-110 active:scale-95"

  return (
    <div className={cn("mx-auto w-full max-w-5xl px-4 pb-10", className)}>
      <div
        className={cn(
          "bg-card relative h-full min-h-[200px] w-full",
          "border-border overflow-hidden rounded-3xl border",
          "flex flex-col items-center justify-center p-6 transition-all"
        )}
      >
        {error ? (
          <StateDisplay
            gif={crying}
            size="h-52 w-52"
            alt="Error gif"
            text="Couldn't find you senpai..."
            textClass="text-red-500 animate-pulse font-bold italicss text-xl"
          />
        ) : !imageResult ? (
          <StateDisplay
            gif={hertaa}
            alt="Hertaa kawai gif"
            size="h-52 w-52"
            text="Transform your anime and manga history into a beautiful, minimalist hexagonal visualization."
            textClass="text-muted"
          />
        ) : (
          <>
            <div className="z-10 flex h-full w-full items-center justify-center">
              <img
                src={imageResult}
                alt="Generated Hexagonal Grid"
                className={cn(
                  "h-auto w-auto max-w-full rounded-xl object-contain drop-shadow-xl transition-all duration-300",
                  isExpanded
                    ? "max-h-none scale-105"
                    : "max-h-[500px] scale-100"
                )}
              />
              <div className="absolute right-6 bottom-6 z-20 flex flex-col gap-3">
                <button
                  onClick={() => setIsExpanded(!isExpanded)}
                  className={cn(
                    iconButton,
                    "hidden border border-white/10 bg-white/10 text-white backdrop-blur-md hover:bg-white/20 md:inline-flex"
                  )}
                  title={isExpanded ? "Shrink View" : "Expand to Full Size"}
                >
                  {isExpanded ? (
                    <ZoomOut className="h-5 w-5" />
                  ) : (
                    <ZoomIn className="h-5 w-5" />
                  )}
                </button>
                <button
                  onClick={() => {
                    const link = document.createElement("a")
                    link.download = `${username}-hexgrid.png`
                    link.href = imageResult
                    link.click()
                  }}
                  className={cn(
                    iconButton,
                    "bg-accent hover:bg-accent-hover text-white"
                  )}
                  title="Download Image"
                >
                  <Download className="h-5 w-5" />
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
