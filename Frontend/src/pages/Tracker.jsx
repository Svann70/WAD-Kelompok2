import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import api from "../api/GetApi";

export default function Tracker() {
  // -------------------------------------------------------------
  // 1. STATE UTAMA
  // -------------------------------------------------------------
  const [daftarLamaran, setDaftarLamaran] = useState([]);
  const [kataKunci, setKataKunci] = useState("");
  const [filterStatus, setFilterStatus] = useState("Semua");
  const [tampilkanModal, setTampilkanModal] = useState(false);
  const [halaman, setHalaman] = useState(1);
  const dataPerHalaman = 5;

  // State untuk form input
  const [perusahaan, setPerusahaan] = useState("");
  const [posisi, setPosisi] = useState("");
  const [tipe, setTipe] = useState("Magang");
  const [tanggal, setTanggal] = useState(new Date().toISOString().split("T")[0]);
  const [status, setStatus] = useState("Terkirim");
  const [catatan, setCatatan] = useState("");
  const [idEdit, setIdEdit] = useState(null);

  useEffect(() => {
    const fetchJobs = async () => {
      try {
        const response = await api.get("/jobs");
        const list = Array.isArray(response.data?.data) ? response.data.data : [];
        const jobsData = list.map((item) => ({
          id: item.id,
          perusahaan: item.company,
          posisi: item.title,
          tipe: item.type,
          tanggal: item.created_at ? item.created_at.split("T")[0] : "",
          status: item.status,
          catatan: item.description || "",
        }));
        setDaftarLamaran(jobsData);
      } catch (error) {
        console.error("Gagal mengambil data lamaran:", error);
      }
    };

    fetchJobs();
  }, []);

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
  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!perusahaan.trim() || !posisi.trim()) {
      alert("Nama perusahaan dan posisi wajib diisi!");
      return;
    }

    if (idEdit !== null) {
      try {
        const response = await api.put(`/jobs/${idEdit}`, {
          company: perusahaan,
          title: posisi,
          type: tipe,
          description: catatan,
          status: status,
        });
        const updated = response.data?.data || {};
        setDaftarLamaran((prev) =>
          prev.map((item) =>
            item.id === idEdit
              ? {
                  id: updated.id || idEdit,
                  perusahaan: updated.company || perusahaan,
                  posisi: updated.title || posisi,
                  tipe: updated.type || tipe,
                  tanggal: updated.created_at ? updated.created_at.split("T")[0] : item.tanggal,
                  status: updated.status || status,
                  catatan: updated.description !== undefined ? updated.description : catatan,
                }
              : item
          )
        );
        handleTutupModal();
      } catch (err) {
        console.error("Gagal update job:", err);
      }
    } else {
      try {
        const response = await api.post("/jobs", {
          company: perusahaan,
          title: posisi,
          type: tipe,
          description: catatan,
          status: status,
        });
        const created = response.data?.data || {};
        const itemBaru = {
          id: created.id || Date.now(),
          perusahaan: created.company || perusahaan,
          posisi: created.title || posisi,
          tipe: created.type || tipe,
          tanggal: created.created_at ? created.created_at.split("T")[0] : tanggal,
          status: created.status || status,
          catatan: created.description !== undefined ? created.description : catatan,
        };
        setDaftarLamaran([itemBaru, ...daftarLamaran]);
        handleTutupModal();
      } catch (err) {
        console.error("Gagal menambah job:", err);
      }
    }
  };

  // Menghapus data lamaran
  const handleHapus = async (id, namaPerusahaan) => {
    if (confirm(`Apakah Anda yakin ingin menghapus lamaran di ${namaPerusahaan}?`)) {
      try {
        await api.delete(`/jobs/${id}`);
        setDaftarLamaran((prevData) => prevData.filter((item) => item.id !== id));
      } catch (error) {
        console.error("Gagal menghapus lamaran:", error);
      }
    }
  };

  // Mengubah status secara langsung dari dropdown baris tabel
  const handleUbahStatus = async (id, statusBaru) => {
    try {
      await api.patch(`/jobs/${id}/status`, { status: statusBaru });
      setDaftarLamaran((prevData) =>
        prevData.map((item) => (item.id === id ? { ...item, status: statusBaru } : item))
      );
    } catch (error) {
      console.error("Gagal memperbarui status:", error);
    }
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

  useEffect(() => {
    setHalaman(1);
  }, [kataKunci, filterStatus]);

  const totalHalaman = Math.ceil(lamaranTersaring.length / dataPerHalaman) || 1;
  const indeksAwal = (halaman - 1) * dataPerHalaman;
  const dataTampil = lamaranTersaring.slice(indeksAwal, indeksAwal + dataPerHalaman);

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
                dataTampil.map((item) => (
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

        {/* Kontrol Pagination */}
        {lamaranTersaring.length > 0 && (
          <div className="flex flex-col sm:flex-row items-center justify-between px-4 py-3 border-t border-gray-200 bg-gray-50/50 gap-2">
            <span className="text-xs text-gray-500">
              Menampilkan {indeksAwal + 1} - {Math.min(indeksAwal + dataPerHalaman, lamaranTersaring.length)} dari {lamaranTersaring.length} lamaran
            </span>
            <div className="flex items-center gap-2">
              <button
                type="button"
                disabled={halaman === 1}
                onClick={() => setHalaman((prev) => Math.max(prev - 1, 1))}
                className="px-3 py-1 text-xs font-medium border border-gray-300 rounded-md bg-white hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                &larr; Prev
              </button>
              <span className="text-xs font-semibold text-gray-700">
                Halaman {halaman} dari {totalHalaman}
              </span>
              <button
                type="button"
                disabled={halaman >= totalHalaman}
                onClick={() => setHalaman((prev) => Math.min(prev + 1, totalHalaman))}
                className="px-3 py-1 text-xs font-medium border border-gray-300 rounded-md bg-white hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                Next &rarr;
              </button>
            </div>
          </div>
        )}
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
