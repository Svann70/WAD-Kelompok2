import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import About from "./pages/About";
import Tracker from "./pages/Tracker";
import Header from "./component/Header";
import Footer from "./component/Footer";
import { featuresData } from "./data/features";

const App = () => {
  return (
    <BrowserRouter>
      <div className="min-h-screen flex flex-col bg-white text-gray-800">
        <Header />
        <main className="flex-1">
          <Routes>
            <Route path="/" element={<Home features={featuresData} />} />
            <Route path="/tracker" element={<Tracker />} />
            <Route path="/about" element={<About />} />
          </Routes>
        </main>
        <Footer />
      </div>
    </BrowserRouter>
  );
};

export default App;
