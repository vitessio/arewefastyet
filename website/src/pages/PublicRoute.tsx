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

import { Route, Routes } from "react-router-dom";

import Layout from "@/pages/Layout";
import Error from "@/utils/Error";
import ComparePage from "./ComparePage/ComparePage";
import DailyPage from "./DailyPage/DailyPage";
import ForeignKeysPage from "./ForeignKeysPage/ForeignKeysPage";
import HistoryPage from "./HistoryPage/HistoryPage";
import HomePage from "./HomePage/HomePage";
import MacroQueriesComparePage from "./MacroQueriesComparePage/MacroQueriesComparePage";
import PRPage from "./PRPage/PRPage";
import PRsPage from "./PRsPage/PRsPage";
import StatusPage from "./StatusPage/StatusPage";

const PublicRoute = () => {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route index element={<HomePage />} />

        <Route path="/home" element={<HomePage />} />
        <Route path="/status" element={<StatusPage />} />
        <Route path="/Daily" element={<DailyPage />} />
        <Route path="/compare" element={<ComparePage />} />
        <Route path="/macrobench/queries/compare" element={<MacroQueriesComparePage />} />
        <Route path="/pr" element={<PRsPage />} />
        <Route path="/pr/:pull_nb" element={<PRPage />} />
        <Route path="/fk" element={<ForeignKeysPage />} />
        <Route path="/history" element={<HistoryPage />} />

        <Route path="*" element={<Error />} />
      </Route>
    </Routes>
  );
};

export default PublicRoute;
