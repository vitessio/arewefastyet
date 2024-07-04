import { Link } from "react-router-dom";
import Icon from "@/common/Icon";
import { Button } from "@/components/ui/button";
import DailySummary from "@/common/DailySummary";
import useApiCall from "@/utils/Hook";
import { MacroDataValue } from "@/types";

interface DailySummarydata {
  name: string;
  data : { total_qps: MacroDataValue }[];
}

export default function HomePageHero() {
  const {
    data: dataDailySummary,
    isLoading: isLoadingDailySummary,
    error: errorDailySummary,
  } = useApiCall<DailySummarydata>(`${import.meta.env.VITE_API_URL}daily/summary`);

  return (
    <section className="flex flex-col items-center h-screen p-page">
      <h1 className="text-6xl font-semibold text-center mt-10 leading-normal">
        Benchmarking <br />
        System for <br />
        <span className="text-orange-500"> Vitess</span>
      </h1>
      <div className="flex gap-x-4 mt-10">
        <Button
          asChild
          size={"lg"}
          variant={"default"}
          className="bg-background hover:bg-muted/90 dark:bg-front dark:hover:bg-front/90 dark:text-background text-foreground rounded-lg border"
        >
          <Link
            className="rounded-2xl p-5 flex items-center gap-x-2"
            to="https://vitess.io/blog/2021-07-08-announcing-vitess-arewefastyet"
            target="__blank"
          >
            Blog Post
            <Icon className="text-2xl" icon="bookmark" />
          </Link>
        </Button>
        <Button
          asChild
          size={"lg"}
          variant={"default"}
          className="bg-front text-background hover:bg-front/90 rounded-lg"
        >
          <Link
            className="rounded-2xl p-5 flex items-center gap-x-2"
            to="https://github.com/vitessio/arewefastyet"
            target="__blank"
          >
            GitHub
            <Icon className="text-2xl" icon="github" />
          </Link>
        </Button>
        <Button
          asChild
          size={"lg"}
          variant={"default"}
          className="bg-background hover:bg-muted/90 dark:bg-front dark:hover:bg-front/90 dark:text-background text-foreground rounded-lg border"
        >
          <Link
            className="rounded-2xl p-5 flex items-center gap-x-2"
            to="https://www.vitess.io"
            target="__blank"
          >
            Vitess
            <Icon className="text-2xl" icon="vitess" />
          </Link>
        </Button>
      </div>
      <h2 className="text-2xl font-medium mt-20">
        Historical results on the <Link className="text-orange-500" to="https://github.com/vitessio/arewefastyet/tree/main" target="_blank">main</Link>{" "}
        branch
      </h2>
      {/* <div className="flex gap-x-8 mt-10">
        <DailySummary
          data={oltpData}
          setBenchmarktype={setBenchmarktype}
          benchmarkType={benchmarkType}
        />
        <DailySummary
          data={tpccData}
          setBenchmarktype={setBenchmarktype}
          benchmarkType={benchmarkType}
        />
      </div> */}
      <Link
        to="/daily"
        className="text-orange-500 text-lg mt-10"
      >
        See more historical results {">"}
      </Link>
    </section>
  );
}
