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

import React from "react";
import { Link } from "react-router-dom";
import ErrorImage from "../../assets/error.png";

import "../Error/error.css";

const Error = () => {
  return (
    <div className="error">
      <div>
        <h1>404</h1>
      </div>
      <div className="errorImg">
        <img src={ErrorImage} alt="error" />
      </div>
      <div>
        <h2>OOPS! Something went wrong</h2>
        <Link to="/home">
          <button className="goHome">Go Back</button>{" "}
        </Link>
      </div>
    </div>
  );
};

export default Error;
