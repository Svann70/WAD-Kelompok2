import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function Header() {
  const navigate = useNavigate();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <header className="border-b border-gray-200 bg-white">
      <div className="container mx-auto flex justify-between items-center px-6 py-4">
        <div className="flex items-center gap-2 cursor-pointer" onClick={() => navigate("/")}>
          <span className="text-xl font-bold text-gray-900 tracking-tight">JobTrack</span>
          <span className="text-xs text-gray-500 border border-gray-200 px-2 py-0.5 rounded">Kelompok 2</span>
        </div>

  
        <button
          className="md:hidden text-sm text-gray-700 px-2 py-1 border border-gray-300 rounded cursor-pointer"
          onClick={() => setIsOpen(!isOpen)}
        >
          {isOpen ? "x" : "☰"}
        </button>

       
        <div className="hidden md:flex items-center gap-6 text-sm text-gray-600">
          <span className="cursor-pointer hover:text-gray-900 transition" onClick={() => navigate("/")}>
            Home
          </span>
          <span className="cursor-pointer hover:text-gray-900 transition" onClick={() => navigate("/about")}>
            About
          </span>
        </div>
      </div>

   
      {isOpen && (
        <div className="md:hidden border-t border-gray-100 px-6 py-3 flex flex-col gap-3 text-sm text-gray-700 bg-gray-50">
          <span className="cursor-pointer py-1" onClick={() => { navigate("/"); setIsOpen(false); }}>
            Home
          </span>
          <span className="cursor-pointer py-1" onClick={() => { navigate("/about"); setIsOpen(false); }}>
            About
          </span>
        </div>
      )}
    </header>
  );
}
