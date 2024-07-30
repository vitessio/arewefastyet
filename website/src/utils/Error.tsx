/*
Copyright 2024 The Vitess Authors.

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
import { Link } from "react-router-dom";
import ErrorImage from "@/assets/error.png";

const Error = () => {
  return (
    <div className="error flex flex-col justify-center items-center mt-12">
      <div className="m-4 flex flex-col items-center">
        <h1 className="text-[3rem] font-bold">404</h1>
      </div>
      <div className="pl-[30px] m-4 flex flex-col items-center">
        <img
          src={ErrorImage}
          alt="error"
          className="p-[10px] w-[260px] h-[225px] errorImgAnimation"
        />
      </div>
      <div className="m-4 flex flex-col items-center">
        <h2 className="font-medium text-center ">OOPS! Something went wrong</h2>
        <Link to="/home">
          <Button className="px-[22px] py-2.5 m-[20px] text-[1rem] rounded-[1rem]  bg-[#e77002]  text-white border-none cursor-pointer flex flex-row items-center gap-[10px]">
            Go Back
          </Button>
        </Link>
      </div>
    </div>
  );
};

export default Error;
