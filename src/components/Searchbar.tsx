import { SearchIcon, Loader2, LayoutGrid } from "lucide-react";
import { FormEvent } from "react";

type SearchbarProps = {
  onSearch: (username: string) => void;
  loading: boolean;
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
    <div className="w-full max-w-xl">
      <form
        onSubmit={handleSubmit}
        className="flex items-center gap-2 p-1.5 bg-[#161B22] border border-white/10 rounded-xl 
                    focus-within:border-accent/50 transition-colors"
      >
        <div className="flex items-center gap-3 px-3 w-full">
          <SearchIcon className="w-5 h-5 text-slate-500" />
          <input
            name="username"
            type="text"
            placeholder="Enter AniList Username (e.g., satoshi_ko)"
            className="bg-transparent w-full text-white placeholder:text-slate-600 outline-none text-sm md:text-base font-medium"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="flex items-center gap-2 bg-accent hover:bg-accent-hover text-white px-4 py-2.5 rounded-lg font-bold 
                    text-sm whitespace-nowrap transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : "Generate Grid"}
          <LayoutGrid className="w-5 h-5" />
        </button>

      </form>
    </div>
  );
}