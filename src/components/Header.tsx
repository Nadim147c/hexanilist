import { Star } from "lucide-react"

export default function Header() {
  return (
    <header className="mx-auto flex w-full max-w-7xl items-center justify-between p-6">
      <div className="group flex shrink-0 cursor-pointer items-center gap-2">
        <div className="bg-accent flex h-8 w-8 rotate-3 items-center justify-center rounded-md font-black text-white transition-transform duration-300 group-hover:rotate-12">
          H
        </div>
        <h1 className="text-text text-2xl leading-none font-bold tracking-tighter">
          Hex<span className="text-muted">AniList</span>
        </h1>
      </div>

      <a
        href="https://github.com/zedxihan/hexanilist"
        target="_blank"
        rel="noopener noreferrer"
        className="group border-border hover:text-accent flex items-center gap-2 rounded-lg border bg-white/5 px-4 py-2 text-sm font-medium text-slate-400 transition-all hover:-translate-y-px hover:border-white/20 hover:bg-white/10 active:scale-95"
      >
        <Star className="text-accent fill-accent/20 h-4 w-4" />
        <span className="hidden sm:inline">Star on GitHub</span>
      </a>
    </header>
  )
}
