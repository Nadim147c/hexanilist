import { Mail, Heart, Pizza } from "lucide-react";

export default function Footer() {
  const currentYear = new Date().getFullYear();
  const linkStyle = "text-muted hover:text-accent transition-colors"

  return (
    <footer className="w-full max-w-7xl mx-auto px-6 pb-5 mt-auto">
      <div className="flex flex-col md:flex-row items-center justify-between gap-6 border-t border-white/5 pt-5">

        <div className="text-muted text-sm font-medium order-2 md:order-1">
          HexAniList, <span className="text-lg">&copy;</span> {currentYear}
        </div>

        <div className="flex items-center gap-6 order-1 md:order-2">
          <a
            href="#"
            target="_blank"
            rel="noopener noreferrer"
            className={linkStyle}
          >
            <Mail className="w-5 h-5" />
          </a>
          <a
            href="#"
            target="_blank"
            rel="noopener noreferrer"
            className={linkStyle}
          >
            <Heart className="w-5 h-5" />
          </a>
          <a
            href="#"
            target="_blank"
            rel="noopener noreferrer"
            className={linkStyle}
          >
            <Pizza className="w-5 h-5" />
          </a>

        </div>
        <div className="flex items-center gap-6 text-sm font-medium order-3">
          <span className="text-muted">Made with ❤️</span>
        </div>

      </div>
    </footer>
  );
}