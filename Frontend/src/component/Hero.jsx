import useLocalStorageState from "../hooks/useLocalStorageState";
import { useNavigate } from "react-router-dom";

export default function Hero() {
  const navigate = useNavigate();

  const handlePindah = () => {
    navigate('/login');
  };

  return (
    <section className="hero">
      <div className="container mx-auto text-center py-20 px-4">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          Welcome to JobTrack
        </h1>
        <p className="text-lg text-gray-600 mb-8 max-w-xl mx-auto">
          Personal Job and Intern Tracker untuk mempermudah pencatatan lamaran kerja dan magang Anda.
        </p>
        <button onClick = {handlePindah}
          className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition duration-300 cursor-pointer font-medium"
        >
          Mulai Sekarang
        </button>
      </div>
    </section>
  );
}
