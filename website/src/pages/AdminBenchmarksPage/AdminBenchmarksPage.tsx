import DataForm from "../../common/DataForm";
import Dropdown from "../../common/Dropdown";
import ExecutionQueue from "./components/ExecutionQueue";
import dummyQueue from "./dummyQueue";

const benchmarkWorkloads = [
  "OLTP-SET",
  "OLTP-READONLY",
  "OLTP-READONLY-OLAP",
  "TPCC",
  "OLTP",
];

export default function AdminBenchmarksPage() {
  return (
    <>
      <section className="p-page my-16">
        <div className="flex flex-col gap-y-8">
          <div className="flex flex-col gap-y-2">
            <h2 className="text-primary text-3xl">Request New Benchmark</h2>
            <p className="text-xs text-front text-opacity-80">
              Enter SHA and select workload to add a new benchmark to the
              execution queue
            </p>
          </div>
          <DataForm.Container
            className="flex items-center gap-x-4"
            onSubmit={(data) => {
              console.log(data);
            }}
          >
            <DataForm.Input
              name="commit"
              className="w-3/5 text-lg px-4 py-2 bg-foreground bg-opacity-5 border border-front border-opacity-20 rounded-md outline-none duration-150 focus:border-primary focus:bg-primary focus:bg-opacity-5"
              placeholder="Enter Commit SHA"
            />
            <Dropdown.Container
              name="workload"
              className="w-[15vw] truncate self-stretch text-lg px-4 py-2 border border-front border-opacity-20 bg-foreground bg-opacity-5 rounded-md"
            >
              {benchmarkWorkloads.map((workload, key) => (
                <Dropdown.Option
                  key={key}
                  className="w-[15vw] relative text-sm border-front border border-t-transparent border-opacity-60 bg-background py-2 after:duration-150 after:absolute-cover 
                  after:bg-foreground after:bg-opacity-0 hover:after:bg-opacity-10 font-medium"
                >
                  {workload}
                </Dropdown.Option>
              ))}
            </Dropdown.Container>
            <DataForm.Input
              type="submit"
              className="cursor-pointer bg-primary px-6 py-2 rounded-md self-stretch font-medium text-back duration-300 hover:brightness-125 hover:scale-105 active:scale-90 active:brightness-90"
              value="Request"
            />
          </DataForm.Container>
        </div>
      </section>

      <figure className="p-page">
        <div className="border border-front border-opacity-10" />
      </figure>

      <ExecutionQueue data={dummyQueue} />
    </>
  );
}
