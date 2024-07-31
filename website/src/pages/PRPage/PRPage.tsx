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

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { PrData } from "@/types";
import useApiCall from "@/utils/Hook";
import { errorApi } from "@/utils/Utils";
import { format, formatDistanceToNow } from "date-fns";
import { Link, useParams } from "react-router-dom";
import RingLoader from "react-spinners/RingLoader";

export default function PRPage() {
  const { pull_nb } = useParams();

  const {
    data: prData,
    isLoading: prLoading,
    error: prError,
  } = useApiCall<PrData>(`${import.meta.env.VITE_API_URL}pr/info/${pull_nb}`);

  if (
    prData?.error ==
    "GET https://api.github.com/repos/vitessio/vitess/pulls/13675: 404 Not Found []"
  ) {
    return (
      <div className="errorPullRequest">Pull request {pull_nb} not found</div>
    );
  }

  return (
    <>
      <section className="flex h-screen flex-col items-center py-[20vh] p-page">
        {prLoading && (
          <div className="loadingSpinner">
            <RingLoader loading={prLoading} color="#E77002" size={300} />
          </div>
        )}

        {prError && (
          <div className="text-red-500 text-center my-2">{errorApi}</div>
        )}

        {!prLoading && prData && (
          <Card className="w-fit border-border">
            <CardHeader>
              <CardTitle className="text-lg md:text-2xl flex flex-row gap-10 justify-between">
                <span>{prData?.Title}</span>
                <Link
                  target="_blank"
                  to={`https://github.com/vitessio/vitess/pull/${pull_nb}`}
                  className="text-primary"
                >
                  #{pull_nb}
                </Link>
              </CardTitle>
              <CardDescription className="py-4">
                <TooltipProvider delayDuration={200}>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <div className="underline text-left">
                        {" "}
                        By {prData?.Author} at{" "}
                        {formatDistanceToNow(prData?.CreatedAt, {
                          addSuffix: true,
                        })}{" "}
                      </div>
                    </TooltipTrigger>
                    <TooltipContent align="start">
                      <p>
                        {format(
                          prData?.CreatedAt,
                          "MMM d, yyyy, h:mm a 'GMT'XXX"
                        )}
                      </p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid w-full items-center gap-4">
                <div className="flex flex-col gap-3">
                  <Label className="md:text-xl font-semibold">Base</Label>
                  {prData?.Base ? (
                    <Link
                      target="_blank"
                      to={`https://github.com/vitessio/vitess/commit/${prData?.Base}`}
                    >
                      <p className="text-xs md:text-lg text-primary">
                        {prData?.Base}
                      </p>
                    </Link>
                  ) : (
                    <p className="text-xs md:text-lg text-red-500">No base</p>
                  )}
                  <Separator />
                  <Label className="md:text-xl font-semibold">Head</Label>
                  {prData?.Head ? (
                    <Link
                      target="_blank"
                      to={`https://github.com/vitessio/vitess/commit/${prData?.Head}`}
                    >
                      <p className="text-xs md:text-lg text-primary">
                        {prData?.Head}
                      </p>
                    </Link>
                  ) : (
                    <p className="text-xs md:text-lg text-red-500">No head</p>
                  )}
                </div>
              </div>
            </CardContent>
            <CardFooter className="flex justify-end">
              <Button
                className="p-4 md:p-8"
                disabled={!prData?.Base || !prData?.Head}
              >
                <Link
                  to={`/compare?old=${prData?.Base}&new=${prData?.Head}`}
                  className="text-lg"
                >
                  Compare with Base Commit
                </Link>
              </Button>
            </CardFooter>
          </Card>
        )}
      </section>
    </>
  );
}
