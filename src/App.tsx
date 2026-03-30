import { useState } from "react";
import Header from "./components/Header";
import Searchbar from "./components/Searchbar";
import HexGridContainer from "./components/HexGridContainer";
import Footer from "./components/Footer";

export default function App() {
  const [loading, setLoading] = useState<boolean>(false)

  const handleSearch = async (username: string) => {
    setLoading(true)
    try {
      console.log("Searching for:", username)
    } catch (err) {
      console.error("Error fetching data:", err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <Header />

      <main className="flex flex-col items-center justify-center space-y-6">

        <h2 className="text-3xl md:text-4xl font-black text-white tracking-tighter text-center uppercase">
          Your AniList Journey, <span className="text-slate-500">Hexagonal</span>
        </h2>
        <Searchbar onSearch={handleSearch} loading={loading} />

        <HexGridContainer />
      </main>

      <Footer />
    </>
  )
}