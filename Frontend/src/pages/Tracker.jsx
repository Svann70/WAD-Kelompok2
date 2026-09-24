import { useState, useEffect } from "react";
import { Link } from "react-router-dom";

// Data awal contoh lamaran kerja & magang
const dataAwal = [
  {
    id: 1,
    perusahaan: "GoTo Financial",
    posisi: "Software Engineer Intern",
    tipe: "Magang",
    tanggal: "2026-09-15",
    status: "Interview",
    catatan: "User interview via Google Meet"
  },
  {
    id: 2,
    perusahaan: "Bank Mandiri",
    posisi: "Backend Developer Intern",
    tipe: "Magang",
    tanggal: "2026-09-18",
    status: "Terkirim",
    catatan: "Kirim CV melalui website karir"
  },
  {
    id: 3,
    perusahaan: "Traveloka",
    posisi: "Junior Frontend Engineer",
    tipe: "Full-time",
    tanggal: "2026-09-10",
    status: "Diterima",
    catatan: "Offering letter sudah masuk email"
  }
];

export default function Tracker() {
  // -------------------------------------------------------------
  // 1. STATE UTAMA
  // -------------------------------------------------------------
  // Mengambil data dari localStorage jika ada, jika tidak pakai dataAwal
  const [daftarLamaran, setDaftarLamaran] = useState(() => {
    const dataTersimpan = localStorage.getItem("jobtrack_data");
    return dataTersimpan ? JSON.parse(dataTersimpan) : dataAwal;
  });

  // State untuk pencarian dan filter status
  const [kataKunci, setKataKunci] = useState("");
  const [filterStatus, setFilterStatus] = useState("Semua");

  // State untuk pop-up modal
  const [tampilkanModal, setTampilkanModal] = useState(false);

  // State untuk form input
  const [perusahaan, setPerusahaan] = useState("");
  const [posisi, setPosisi] = useState("");
  const [tipe, setTipe] = useState("Magang");
  const [tanggal, setTanggal] = useState(new Date().toISOString().split("T")[0]);
  const [status, setStatus] = useState("Terkirim");
  const [catatan, setCatatan] = useState("");

  // State untuk mode edit (ID data yang sedang diedit, atau null jika tambah baru)
  const [idEdit, setIdEdit] = useState(null);

  // -------------------------------------------------------------
  // 2. EFEK: SIMPAN OTOMATIS KE LOCALSTORAGE
  // -------------------------------------------------------------
  useEffect(() => {
    localStorage.setItem("jobtrack_data", JSON.stringify(daftarLamaran));
  }, [daftarLamaran]);

  // -------------------------------------------------------------
  // 3. FUNGSI-FUNGSI AKSI
  // -------------------------------------------------------------

  // Reset form ke nilai kosong/default
  const resetForm = () => {
    setPerusahaan("");
    setPosisi("");
    setTipe("Magang");
    setTanggal(new Date().toISOString().split("T")[0]);
    setStatus("Terkirim");
    setCatatan("");
    setIdEdit(null);
  };

  // Buka pop-up untuk tambah lamaran baru
  const handleBukaTambah = () => {
    resetForm();
    setTampilkanModal(true);
  };

  // Buka pop-up untuk edit lamaran yang dipilih
  const handlePilihEdit = (item) => {
    setIdEdit(item.id);
    setPerusahaan(item.perusahaan);
    setPosisi(item.posisi);
    setTipe(item.tipe);
    setTanggal(item.tanggal);
    setStatus(item.status);
    setCatatan(item.catatan || "");
    setTampilkanModal(true);
  };

  // Tutup pop-up modal
  const handleTutupModal = () => {
    resetForm();
    setTampilkanModal(false);
  };

  // Menangani submit form di dalam pop-up
  const handleSubmit = (e) => {
    e.preventDefault();

    if (!perusahaan.trim() || !posisi.trim()) {
      alert("Nama perusahaan dan posisi wajib diisi!");
      return;
    }

    if (idEdit !== null) {
      // Perbarui data yang diedit
      const dataTerupdate = daftarLamaran.map((item) => {
        if (item.id === idEdit) {
          return {
            ...item,
            perusahaan,
            posisi,
            tipe,
            tanggal,
            status,
            catatan
          };
        }
        return item;
      });
      setDaftarLamaran(dataTerupdate);
    } else {
      // Tambah data baru
      const lamaranBaru = {
        id: Date.now(),
        perusahaan,
        posisi,
        tipe,
        tanggal,
        status,
        catatan
      };
      setDaftarLamaran([lamaranBaru, ...daftarLamaran]);
    }

    handleTutupModal();
  };

  // Menghapus data lamaran
  const handleHapus = (id, namaPerusahaan) => {
    const konfirmasi = confirm(`Hapus catatan lamaran di ${namaPerusahaan}?`);
    if (konfirmasi) {
      const sisaData = daftarLamaran.filter((item) => item.id !== id);
      setDaftarLamaran(sisaData);
    }
  };

  // Mengubah status secara langsung dari dropdown baris tabel
  const handleUbahStatus = (id, statusBaru) => {
    const dataTerupdate = daftarLamaran.map((item) => {
      if (item.id === id) {
        return { ...item, status: statusBaru };
      }
      return item;
    });
    setDaftarLamaran(dataTerupdate);
  };

  // -------------------------------------------------------------
  // 4. FILTER DATA DAN HITUNG RINGKASAN
  // -------------------------------------------------------------
  const lamaranTersaring = daftarLamaran.filter((item) => {
    const cocokKataKunci =
      item.perusahaan.toLowerCase().includes(kataKunci.toLowerCase()) ||
      item.posisi.toLowerCase().includes(kataKunci.toLowerCase());

    const cocokStatus =
      filterStatus === "Semua" || item.status === filterStatus;

    return cocokKataKunci && cocokStatus;
  });

  const totalSemua = daftarLamaran.length;
  const totalInterview = daftarLamaran.filter((item) => item.status === "Interview").length;
  const totalDiterima = daftarLamaran.filter((item) => item.status === "Diterima").length;

  return (
    <div className="container mx-auto px-6 py-12 max-w-5xl">
      {/* Header Halaman */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 tracking-tight">
            Job & Intern Tracker
          </h1>
          <p className="text-sm text-gray-600 mt-1">
            Catat dan pantau seluruh perkembangan lamaran kerja serta magang Anda.
          </p>
        </div>
        <Link
          to="/"
          className="text-sm font-medium text-blue-600 hover:underline"
        >
          &larr; Kembali ke Home
        </Link>
      </div>

      {/* Ringkasan Singkat (3 Kartu Bersih) */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
        <div className="border border-gray-200 rounded-xl p-4 bg-white">
          <span className="text-xs text-gray-500 font-medium">Total Lamaran</span>
          <p className="text-2xl font-bold text-gray-900 mt-1">{totalSemua}</p>
        </div>
        <div className="border border-gray-200 rounded-xl p-4 bg-white">
          <span className="text-xs text-amber-600 font-medium">Tahap Interview</span>
          <p className="text-2xl font-bold text-amber-700 mt-1">{totalInterview}</p>
        </div>
        <div className="border border-gray-200 rounded-xl p-4 bg-white">
          <span className="text-xs text-green-600 font-medium">Diterima (Offer)</span>
          <p className="text-2xl font-bold text-green-700 mt-1">{totalDiterima}</p>
        </div>
      </div>

      {/* TOOLBAR: PENCARIAN, FILTER, DAN TOMBOL TAMBAH */}
      <div className="flex flex-col sm:flex-row justify-between items-center gap-3 mb-4">
        <div className="flex flex-col sm:flex-row items-center gap-3 w-full sm:w-auto">
          <input
            type="text"
            placeholder="Cari perusahaan atau posisi..."
            value={kataKunci}
            onChange={(e) => setKataKunci(e.target.value)}
            className="w-full sm:w-64 px-3 py-2 text-sm border border-gray-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          />

          <div className="flex items-center gap-2 w-full sm:w-auto">
            <span className="text-xs text-gray-500 font-medium">Status:</span>
            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              className="w-full sm:w-auto px-3 py-2 text-sm border border-gray-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="Semua">Semua</option>
              <option value="Terkirim">Terkirim</option>
              <option value="Interview">Interview</option>
              <option value="Diterima">Diterima</option>
              <option value="Ditolak">Ditolak</option>
            </select>
          </div>
        </div>

        <button
          onClick={handleBukaTambah}
          className="w-full sm:w-auto bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition cursor-pointer"
        >
          + Tambah Lamaran
        </button>
      </div>

      {/* TABEL PELACAK LAMARAN */}
      <div className="border border-gray-200 rounded-xl overflow-hidden bg-white shadow-xs">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="bg-gray-50 border-b border-gray-200 text-gray-600 text-xs font-semibold">
              <tr>
                <th className="py-3 px-4">Perusahaan</th>
                <th className="py-3 px-4">Posisi</th>
                <th className="py-3 px-4">Tipe</th>
                <th className="py-3 px-4">Tanggal</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Catatan</th>
                <th className="py-3 px-4 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {lamaranTersaring.length === 0 ? (
                <tr>
                  <td colSpan="7" className="py-8 text-center text-gray-500 text-sm">
                    Tidak ada catatan lamaran yang cocok.
                  </td>
                </tr>
              ) : (
                lamaranTersaring.map((item) => (
                  <tr key={item.id} className="hover:bg-gray-50/60 transition">
                    <td className="py-3 px-4 font-medium text-gray-900">
                      {item.perusahaan}
                    </td>
                    <td className="py-3 px-4 text-gray-700">{item.posisi}</td>
                    <td className="py-3 px-4 text-gray-600 text-xs">{item.tipe}</td>
                    <td className="py-3 px-4 text-gray-500 text-xs">{item.tanggal}</td>
                    <td className="py-3 px-4">
                      <select
                        value={item.status}
                        onChange={(e) => handleUbahStatus(item.id, e.target.value)}
                        className="text-xs font-medium px-2 py-1 rounded-md border border-gray-300 bg-white hover:border-gray-400 focus:outline-none focus:ring-1 focus:ring-blue-500 cursor-pointer"
                      >
                        <option value="Terkirim">Terkirim</option>
                        <option value="Interview">Interview</option>
                        <option value="Diterima">Diterima</option>
                        <option value="Ditolak">Ditolak</option>
                      </select>
                    </td>
                    <td className="py-3 px-4 text-gray-600 text-xs max-w-xs truncate">
                      {item.catatan || "-"}
                    </td>
                    <td className="py-3 px-4 text-center">
                      <div className="flex items-center justify-center gap-3 text-xs">
                        <button
                          onClick={() => handlePilihEdit(item)}
                          className="text-blue-600 hover:text-blue-800 font-medium cursor-pointer"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleHapus(item.id, item.perusahaan)}
                          className="text-red-600 hover:text-red-800 font-medium cursor-pointer"
                        >
                          Hapus
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* POP-UP MODAL (UNTUK TAMBAH & EDIT LAMARAN) */}
      {tampilkanModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4">
          <div className="bg-white rounded-xl border border-gray-200 shadow-xl max-w-lg w-full p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-base font-bold text-gray-900">
                {idEdit !== null ? "Edit Data Lamaran" : "Tambah Lamaran Baru"}
              </h2>
              <button
                type="button"
                onClick={handleTutupModal}
                className="text-gray-400 hover:text-gray-600 text-lg cursor-pointer leading-none"
              >
                ✕
              </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-gray-700 mb-1">
                    Nama Perusahaan *
                  </label>
                  <input
                    type="text"
                    required
                    placeholder="Contoh: Tokopedia, GoTo, Bank Mandiri"
                    value={perusahaan}
                    onChange={(e) => setPerusahaan(e.target.value)}
                    className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-gray-700 mb-1">
                    Posisi / Role *
                  </label>
                  <input
                    type="text"
                    required
                    placeholder="Contoh: Frontend Developer Intern"
                    value={posisi}
                    onChange={(e) => setPosisi(e.target.value)}
                    className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-gray-700 mb-1">
                    Tipe Pekerjaan
                  </label>
                  <select
                    value={tipe}
                    onChange={(e) => setTipe(e.target.value)}
                    className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                  >
                    <option value="Magang">Magang (Intern)</option>
                    <option value="Full-time">Full-time</option>
                    <option value="Part-time">Part-time</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-gray-700 mb-1">
                    Tanggal Melamar
                  </label>
                  <input
                    type="date"
                    value={tanggal}
                    onChange={(e) => setTanggal(e.target.value)}
                    className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-gray-700 mb-1">
                    Status Lamaran
                  </label>
                  <select
                    value={status}
                    onChange={(e) => setStatus(e.target.value)}
                    className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
                  >
                    <option value="Terkirim">Terkirim</option>
                    <option value="Interview">Interview</option>
                    <option value="Diterima">Diterima</option>
                    <option value="Ditolak">Ditolak</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-gray-700 mb-1">
                  Catatan / Progres
                </label>
                <input
                  type="text"
                  placeholder="Contoh: Jadwal interview tgl 25, tes online selesai, dsb."
                  value={catatan}
                  onChange={(e) => setCatatan(e.target.value)}
                  className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div className="flex justify-end gap-3 pt-2">
                <button
                  type="button"
                  onClick={handleTutupModal}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-sm text-gray-700 hover:bg-gray-50 transition cursor-pointer"
                >
                  Batal
                </button>
                <button
                  type="submit"
                  className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2 rounded-lg text-sm font-medium transition cursor-pointer"
                >
                  {idEdit !== null ? "Simpan Perubahan" : "Tambah Lamaran"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
