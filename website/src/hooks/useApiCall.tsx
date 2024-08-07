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

import { QueryKey, useQuery } from "@tanstack/react-query";

const fetchData = async (url: string) => {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error("Network response was not ok");
  }
  return response.json();
};

export default function useApiCall<T>({
  url,
  queryKey,
}: {
  url: string | null;
  queryKey: QueryKey;
}): { data: T | undefined; error: Error | null; isLoading: boolean } {
  console.log({ url, queryKey });
  const { data, error, isLoading } = useQuery<T, Error, T, QueryKey>({
    queryKey: [queryKey],
    queryFn: async () => fetchData(url!),
    enabled: !!url, // This ensures the query is only run if the URL is truthy
  });

  return {
    data: data,
    error: error,
    isLoading,
  };
}
