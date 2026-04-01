import { Mail, Heart, Pizza } from "lucide-react"

export default function Footer() {
  const currentYear = new Date().getFullYear()
  const linkStyle = "text-muted hover:text-accent transition-colors"

  return (
    <footer className="mx-auto mt-auto w-full max-w-7xl px-6 pb-5">
      <div className="flex flex-col items-center justify-between gap-6 border-t border-white/5 pt-5 md:flex-row">
        <div className="order-2 flex flex-col text-sm leading-tight font-medium items-center md:order-1">
          <div className="text-muted">
            HexAniList, <span className="text-lg leading-none">©</span>{" "}
            {currentYear}
          </div>
          <div className="text-muted/60 -mt-0.5 text-[11px]">
            Not affiliated with AniList
          </div>
        </div>

        <div className="order-1 flex items-center gap-6 md:order-2">
          <a
            href="#"
            target="_blank"
            rel="noopener noreferrer"
            className={linkStyle}
          >
            <Mail className="h-5 w-5" />
          </a>
          <a
            href="#"
            target="_blank"
            rel="noopener noreferrer"
            className={linkStyle}
          >
            <Heart className="h-5 w-5" />
          </a>
          <a
            href="#"
            target="_blank"
            rel="noopener noreferrer"
            className={linkStyle}
          >
            <Pizza className="h-5 w-5" />
          </a>
        </div>
        <div className="order-3 flex items-center gap-6 text-sm font-medium">
          <span className="text-muted">Made with ❤️</span>
        </div>
      </div>
    </footer>
  )
}
