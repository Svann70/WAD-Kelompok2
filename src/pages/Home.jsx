import Hero from "../component/Hero";
import CardGrid from "../component/CardGrid";

export default function Home({ features }) {
  return (
    <section className="container mx-auto px-4 pb-12">
      <Hero />
      <div className="max-w-5xl mx-auto mt-4">
        <h2 className="text-xl font-bold text-gray-900 mb-2 text-center">Fitur Aplikasi</h2>
        <p className="text-sm text-gray-500 text-center mb-6">Rancangan fitur yang akan dikembangkan dalam JobTrack.</p>
        <CardGrid fitur={features} />
      </div>
    </section>
  );
}
