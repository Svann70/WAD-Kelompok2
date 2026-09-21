import { Link } from "react-router-dom";

export default function Tracker() {
  return (
    <div className="container mx-auto px-6 py-20 text-center max-w-xl">
      <h1 className="text-2xl font-bold text-gray-900 mb-3">
        Sheet Pelacak Lamaran
      </h1>
      <p className="text-sm text-gray-600 mb-6">
        Halaman ini sedang dalam tahap awal persiapan. Modul tabel pelacak lamaran akan dikembangkan pada tahap selanjutnya.
      </p>
      <Link
        to="/"
        className="text-sm font-medium text-blue-600 hover:underline"
      >
        Kembali ke Home
      </Link>
    </div>
  );
}
