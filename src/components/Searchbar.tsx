import { SearchIcon, Loader2, LayoutGrid } from "lucide-react"
import type { FormEvent } from "react"

type SearchbarProps = {
  onSearch: (username: string) => void
  loading: boolean
}

export default function Searchbar({ onSearch, loading }: SearchbarProps) {
  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)
    const username = formData.get("username") as string

    if (typeof username === "string" && username.trim()) {
      onSearch(username.trim())
    }
  }

  return (
    <div className="w-full max-w-xl px-4 md:px-0">
      <form
        onSubmit={handleSubmit}
        className="focus-within:border-accent/50 flex items-center gap-2 rounded-xl border border-white/10 bg-[#161B22] p-1.5 transition-colors"
      >
        <div className="flex w-full min-w-0 items-center gap-3 px-3">
          <SearchIcon className="h-4 w-4 shrink-0 text-slate-500 md:h-5 md:w-5" />
          <input
            name="username"
            type="text"
            placeholder="AniList Username (e.g. satoshi_ko)"
            className="w-full bg-transparent text-sm font-medium text-white outline-none placeholder:text-xs placeholder:text-slate-600 sm:placeholder:text-sm md:text-base"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="bg-accent hover:bg-accent-hover flex items-center gap-2 rounded-lg px-2 py-2.5 text-xs font-bold whitespace-nowrap text-white transition-all active:scale-95 disabled:cursor-not-allowed disabled:opacity-50 md:px-4 md:text-sm"
        >
          {loading ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            "Generate Grid"
          )}
          <LayoutGrid className="h-4 w-4 md:h-5 md:w-5" />
        </button>
      </form>
    </div>
  )
}
