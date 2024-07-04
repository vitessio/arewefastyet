import { Card, CardContent } from "@/components/ui/card";

export default function Diagram() {
  return (
    <section className="p-page flex flex-col items-center my-8">
      <h2 className="text-4xl font-semibold my-14 text-primary dark:text-front">
        Architecture
      </h2>
      <Card className="w-full max-w-screen-xl my-8 border-border">
        <CardContent className="flex justify-center">
          <img
            className="w-full"
            src="/images/execution-pipeline-dark.png"
            alt="execution pipeline"
          />
        </CardContent>
      </Card>
    </section>
  );
}
