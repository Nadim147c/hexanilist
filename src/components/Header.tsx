import { Star } from "lucide-react";

export default function Header() {
  return (
    <header className="w-full max-w-7xl mx-auto flex items-center justify-between p-6">

      <div className="flex items-center gap-2 group cursor-pointer shrink-0">
        <div className="w-8 h-8 bg-accent rounded-md flex items-center justify-center rotate-3 
                        group-hover:rotate-12 transition-transform duration-300 text-white font-black">
          H
        </div>
        <h1 className="text-2xl font-bold tracking-tighter text-text leading-none">
          Hex<span className="text-muted">AniList</span>
        </h1>
      </div>

      <a
        href="https://github.com/zedihan/hexanilist"
        target="_blank"
        rel="noopener noreferrer"
        className="group flex items-center gap-2 px-4 py-2 rounded-lg border border-border bg-white/5
                   text-slate-400 text-sm font-medium transition-all hover:bg-white/10 hover:text-accent
                   hover:border-white/20 hover:-translate-y-[1px] active:scale-95"
      >
        <Star className="w-4 h-4 text-accent fill-accent/20" />
        <span className="hidden sm:inline">Star on GitHub</span>
      </a>

    </header>
  );
}