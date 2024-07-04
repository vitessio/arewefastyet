import PropTypes from "prop-types";
import { ResponsiveLine } from "@nivo/line";
import { twMerge } from "tailwind-merge";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

const DailySummary = ({ data, setBenchmarktype, benchmarkType }) => {
  const orange = "#E77002";

  const transformedData = [];

  if (data.data !== null) {
    transformedData.push({
      id: "QPSTotal",
      data: data.data.map((item) => ({
        y: item.total_qps.center,
      })),
    });
  }

  const getBenchmarkType = () => {
    setBenchmarktype(data.name);
  };

  return (
    <Card
      className={twMerge(
        "cursor-pointer hover:bg-accent duration-300 hover:scale-105 w-[20vw] h-[15vh]",
        benchmarkType === data.name && "bg-accent brightness-125 border-2"
      )}
      onClick={getBenchmarkType}
    >
      <CardHeader>
        <CardTitle className="font-extralight">{data.name}</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col p-3">
        {data ? (
          <div className="w-full h-1/2 p-2 pt-4 flex flex-nowrap justify-center items-center">
            <ResponsiveLine
              data={transformedData}
              margin={{ top: 20, right: 30, bottom: 50, left: 10 }}
              height={300}
              enableGridX={false}
              enableGridY={false}
              colors={orange}
              axisBottom={null}
              axisLeft={null}
            />
          </div>
        ) : null}
        <div className="flex justify-between flex-1 items-center mt-2">
          <i className="fa-solid fa-arrow-right daily--fa-arrow-right"></i>
        </div>
      </CardContent>
    </Card>
  );
};

DailySummary.propTypes = {
  data: PropTypes.shape({
    name: PropTypes.string.isRequired,
    data: PropTypes.array,
  }),
  setBenchmarktype: PropTypes.func.isRequired,
  isSelected: PropTypes.bool,
};

export default DailySummary;
