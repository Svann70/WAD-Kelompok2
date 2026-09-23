import { Link } from "react-router-dom";

export default function Login () {
    return (
    <div className="container mx-auto px-6 py-20 text-center max-w-xl">
      <div className="bg-white p-8 rounded-xl shadow-md border border-gray-200">
        <h1 className="text-4xl font-bold text-gray-900 mb-3">
        JobTrack
      </h1>
      <p className="text-sm text-gray-600 mb-6">
        Silahkan Login atau buat akun terlebih dahulu.
      </p>
        <form className="flex-column mb-5">
          <input type="email" className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none
  focus:ring-2 focus:ring-blue-500 mb-5" placeholder="Masukkan Email"/>
          <input type="password" className=" w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none
  focus:ring-2 focus:ring-blue-500" placeholder="Masukkan Password"/>
        </form>
        <button type="submit" className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition
  font-medium">
    Login
  </button>
      </div>
      <Link
        to="/"
        className="text-sm font-medium text-blue-600 hover:underline"
      >
        Kembali ke Home
      </Link>
    </div>
  );
}