export default function Card({ title, description }) {
  return (
    <div className="p-5 bg-white rounded-xl border border-gray-200 flex flex-col gap-2">
      <h3 className="text-base font-semibold text-gray-900">{title}</h3>
      <p className="text-sm text-gray-600 leading-relaxed">{description}</p>
    </div>
  );
}
