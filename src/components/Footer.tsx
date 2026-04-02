import { Pizza, Code, Bug } from "lucide-react"

export default function Footer() {
  const currentYear = new Date().getFullYear()

  return (
    <footer className="mx-auto mt-auto w-full max-w-7xl px-6 pb-5">
      <div className="flex flex-col items-center justify-between gap-6 border-t border-white/5 pt-5 md:flex-row">
        <div className="order-2 flex flex-col items-center text-sm leading-tight font-medium md:order-1">
          <div className="text-muted">
            HexAniList, <span className="text-lg leading-none">©</span>{" "}
            {currentYear}
          </div>
          <div className="text-muted/60 -mt-0.5 text-[11px]">
            Not affiliated with AniList
          </div>
        </div>

        <div className="order-1 flex items-center gap-6 md:order-2">
          <LinkIcon
            href="http://github.com/Nadim147c/hexanilist/issues"
            tooltip="Report an issue"
          >
            <Bug className="h-5 w-5" />
          </LinkIcon>
          <LinkIcon
            href="http://github.com/Nadim147c/hexanilist"
            tooltip="Source code"
          >
            <Code className="h-5 w-5" />
          </LinkIcon>
          <LinkIcon href="http://patreon.com/Nadim147c" tooltip="Donate">
            <Pizza className="h-5 w-5" />
          </LinkIcon>
        </div>
        <div className="order-3 flex items-center gap-6 text-sm font-medium">
          <span className="text-muted cursor-default select-none hover:animate-pulse">
            Made with ❤️
          </span>
        </div>
      </div>
    </footer>
  )
}

function LinkIcon({
  href,
  children,
  tooltip,
}: {
  href: string
  children: React.ReactNode
  tooltip: string
}) {
  return (
    <div className="group/link-icon relative">
      <a
        href={href}
        target="_blank"
        rel="noopener noreferrer"
        className="text-muted hover:text-accent transition-colors"
      >
        {children}
      </a>
      <div className="bg-card absolute -top-10 left-1/2 hidden w-max -translate-x-1/2 rounded-lg px-2 py-1 group-hover/link-icon:block">
        <span className="text-muted text-sm">{tooltip}</span>
      </div>
    </div>
  )
}
