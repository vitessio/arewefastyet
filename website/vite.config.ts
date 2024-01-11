// Copyright 2023 The Vitess Authors.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import dotenv from "dotenv";

export default defineConfig(({ mode }) => {
  // Load environment variables from the corresponding .env file based on the mode (dev or prod)
  dotenv.config({ path: `./.env.${mode}` });

  return {
    plugins: [react()],
  };
});
