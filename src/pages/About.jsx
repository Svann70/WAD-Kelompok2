import { Link } from "react-router-dom";

export default function About() {
  const anggota = [
    { nama: "Nama Anggota 1", nim: "1202220001" },
    { nama: "Nama Anggota 2", nim: "1202220002" },
    { nama: "Nama Anggota 3", nim: "1202220003" }
  ];

  return (
    <div className="container mx-auto px-6 py-16 max-w-2xl">
      <h1 className="text-3xl font-bold text-gray-900 mb-3">Tentang Proyek</h1>
      <p className="text-sm text-gray-600 leading-relaxed mb-6">
        JobTrack adalah proyek tugas kelompok mata kuliah Web Application Development (WAD) untuk membuat aplikasi pelacak lamaran magang dan kerja sederhana.
      </p>

      <div className="border border-gray-200 rounded-xl p-6 bg-white mb-8">
        <h2 className="text-base font-bold text-gray-900 mb-4">Anggota Kelompok 2</h2>
        <div className="space-y-3">
          {anggota.map((item, index) => (
            <div
              key={index}
              className="flex justify-between items-center py-2 border-b border-gray-100 last:border-b-0 text-sm"
            >
              <span className="font-medium text-gray-800">{item.nama}</span>
              <span className="text-gray-500 font-mono text-xs">{item.nim}</span>
            </div>
          ))}
        </div>
      </div>

      <Link to="/" className="text-sm font-medium text-blue-600 hover:underline">
        Kembali ke Home
      </Link>
    </div>
  );
}
