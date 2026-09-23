import Card from "./Card";

function CardGrid({ fitur }) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 my-8">
      {fitur?.map((feature, index) => (
        <Card
          key={index}
          title={feature.title}
          description={feature.subtitle}
        />
      ))}
    </div>
  );
}

export default CardGrid;
