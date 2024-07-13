/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/



import StatusHero from "./components/StatusHero";

import PreviousExecutionQueue from "./components/PreviousExecutions/PreviousExecutionQueue";
import { Separator } from "@/components/ui/separator";
import { statusDataTypes } from "@/types";
import useApiCall from "@/utils/Hook";

export default function StatusPage() {
  const {
    data: dataQueue,
    isLoading: isLoadingQueue,
    error: errorQueue,
  } = useApiCall<statusDataTypes>(`${import.meta.env.VITE_API_URL}queue`);
  console.log({dataQueue});
  // const { data: dataPreviousExe, isLoading: isLoadingPreviousExe } = useApiCall(
  //   `${import.meta.env.VITE_API_URL}recent`
  // );

  // const [filters, setFilters] = useState({
  //   type: "",
  //   status: "",
  //   source: "",
  // });

  // const filterData = (data: dataTypes[]) => {
  //   return data.filter((item) => {
  //     const itemType = item.type_of ? item.type_of.toString() : "";
  //     const itemSource = item.source ? item.source.toString() : "";
  //     const itemStatus = item.status ? item.status.toString() : "";

  //     const matchesType =
  //       filters.type === "" ||
  //       filters.type === "All" ||
  //       itemType === filters.type;
  //     const matchesSource =
  //       filters.source === "" ||
  //       filters.source === "All" ||
  //       itemSource === filters.source;
  //     const matchesStatus =
  //       filters.status === "" ||
  //       filters.status === "All" ||
  //       itemStatus === filters.status;

  //     return matchesType && matchesSource && matchesStatus;
  //   });
  // };

  // const handleFilterChange = (name: string, value: string) => {
  //   setFilters((prevFilters) => ({
  //     ...prevFilters,
  //     [name]: value === "" ? "" : value,
  //   }));
  // };

  // const filteredDataQueue = filterData(dataQueue) as dataTypes[];
  // const filteredPreviousDataExe = filterData(dataPreviousExe) as dataTypes[];

  return (
    <>
      <StatusHero />
      <Separator className="w-[4/5]"/>
      <PreviousExecutionQueue />
      <div className="border-accent border mt-5" />

      {/* FILTERS OPTIONS*/}
      {/* <div className="flex flex-col items-center">
        <div className="flex p-5 gap-4 ">
          <Select
            value={filters.type}
            onValueChange={(value) => handleFilterChange("type", value)}
          >
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="All">All</SelectItem>
              {[
                ...new Set(dataPreviousExe.map((item: any) => item.type_of)),
              ].map((type) => (
                <SelectItem key={type} value={type}>
                  {type}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            onValueChange={(value) => handleFilterChange("source", value)}
            value={filters.source}
          >
            <SelectTrigger className="w-[200px]">
              <SelectValue placeholder="Sources" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="All">All</SelectItem>
              {[
                ...new Set(dataPreviousExe.map((item: any) => item.source)),
              ].map((source) => (
                <SelectItem key={source} value={source}>
                  {source}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            onValueChange={(value) => handleFilterChange("status", value)}
            value={filters.status}
          >
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="All">All</SelectItem>
              {[
                ...new Set(dataPreviousExe.map((item: any) => item.status)),
              ].map((status) => (
                <SelectItem key={status} value={status}>
                  {status}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div> */}

        {/* EXECUTION QUEUE */}
        {/* {!isLoadingQueue && dataQueue && dataQueue.length > 0 && (
          <ExecutionQueue data={filteredDataQueue} title={"Execution Queue"} />
        )} */}

        {/* PREVIOUS EXECUTIONS */}
        {/* {!isLoadingPreviousExe &&
          dataPreviousExe &&
          dataPreviousExe.length > 0 && (
            <PreviousExecutions
              data={filteredPreviousDataExe}
              title={"Previous Executions"}
            />
          )} */}

        {/* SHOW LOADER BENEATH IF EITHER IS LOADING */}
        {/* {(isLoadingPreviousExe || isLoadingQueue) && (
          <div className="flex justify-center w-full my-16">
            <RingLoader
              loading={isLoadingPreviousExe || isLoadingQueue}
              color="#E77002"
              size={300}
            />
          </div>
        )}

        {errorQueue && (
          <div className="my-10 text-center text-red-500">{errorQueue}</div>
        )} */}
      {/* </div> */}
    </>
  );
}
