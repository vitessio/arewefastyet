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

import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MacroQueriesCompareTable } from "./MacroQueriesCompareTable";
import { columns, MacroQueriesPlan } from "./Columns";
import { FilterConfigs } from "@/types";

// Helper function to create mock query data
const createMockQuery = (index: number): MacroQueriesPlan => {
  return {
    key: `SELECT * FROM table${index}`,
    exec_time_diff: index * 5,
    exec_count_diff: index * 2,
    errors_diff: 0,
    rows_returned_diff: index * -1,
    same_plan: true,
    left: {
      key: `SELECT * FROM table${index}`,
      value: {
        query_type: "SELECT",
        original: `SELECT * FROM table${index} WHERE id = ?`,
        instructions: `{"opcode": "Route", "table": "table${index}"}`,
        exec_count: 100 + index,
        exec_time: 50 + index,
        shard_queries: 1,
        rows_returned: 100 + index,
        rows_affected: 0,
        errors: 0,
        tables_used: `table${index}`,
      },
    },
    right: {
      key: `SELECT * FROM table${index}`,
      value: {
        query_type: "SELECT",
        original: `SELECT * FROM table${index} WHERE id = ?`,
        instructions: `{"opcode": "Route", "table": "table${index}", "optimized": true}`,
        exec_count: 105 + index,
        exec_time: 55 + index,
        shard_queries: 1,
        rows_returned: 99 + index,
        rows_affected: 0,
        errors: 0,
        tables_used: `table${index}`,
      },
    },
  };
};

describe("MacroQueriesCompareTable", () => {
  const filterConfigs: FilterConfigs[] = [
    {
      column: "query",
      title: "Operators",
      options: ["select", "insert", "update", "delete"].map((value) => {
        return { label: value, value: value };
      }),
    },
  ];

  describe("pagination", () => {
    it("displays correct query details when clicking on first page", async () => {
      const user = userEvent.setup();
      // Create 15 queries to ensure pagination (default page size is 10)
      const mockData = Array.from({ length: 15 }, (_, i) => createMockQuery(i));

      render(
        <MacroQueriesCompareTable
          columns={columns}
          data={mockData}
          filterConfigs={filterConfigs}
        />
      );

      // Should be on page 1, showing queries 0-9
      const rows = screen.getAllByRole("row");
      // +1 for header row
      expect(rows.length).toBe(11);

      // Click on the first data row (index 0)
      const firstRow = rows[1];
      await user.click(firstRow);

      // Dialog should show data for query 0
      expect(screen.getByText("50")).toBeInTheDocument(); // left exec_time for query 0
      expect(screen.getByText("55")).toBeInTheDocument(); // right exec_time for query 0
    });

    it("displays correct query details when clicking on second page", async () => {
      const user = userEvent.setup();
      // Create 15 queries to ensure pagination (default page size is 10)
      const mockData = Array.from({ length: 15 }, (_, i) => createMockQuery(i));

      render(
        <MacroQueriesCompareTable
          columns={columns}
          data={mockData}
          filterConfigs={filterConfigs}
        />
      );

      // Navigate to page 2
      const nextPageButton = screen.getByRole("button", { name: /next page/i });
      await user.click(nextPageButton);

      // Should now be on page 2, showing queries 10-14
      const rows = screen.getAllByRole("row");
      // +1 for header row, 5 data rows on page 2
      expect(rows.length).toBe(6);

      // Verify we're showing the correct queries on page 2
      expect(screen.getByText("SELECT * FROM table10")).toBeInTheDocument();
      expect(screen.queryByText("SELECT * FROM table0")).not.toBeInTheDocument();

      // Click on the first data row on page 2 (which is index 10 in the full dataset)
      const firstRowOnPage2 = rows[1];
      await user.click(firstRowOnPage2);

      // Dialog should show data for query 10, NOT query 0
      // Query 10 has exec_time of 50 + 10 = 60 for left, 55 + 10 = 65 for right
      // Query 10 has rows_returned of 100 + 10 = 110 for left, 99 + 10 = 109 for right
      expect(screen.getByText("Statistics")).toBeInTheDocument();
      expect(screen.getByText("110")).toBeInTheDocument(); // left rows_returned for query 10
      expect(screen.getByText("109")).toBeInTheDocument(); // right rows_returned for query 10
    });

    it("displays correct query details when clicking on last row of second page", async () => {
      const user = userEvent.setup();
      // Create 15 queries
      const mockData = Array.from({ length: 15 }, (_, i) => createMockQuery(i));

      render(
        <MacroQueriesCompareTable
          columns={columns}
          data={mockData}
          filterConfigs={filterConfigs}
        />
      );

      // Navigate to page 2
      const nextPageButton = screen.getByRole("button", { name: /next page/i });
      await user.click(nextPageButton);

      // Click on the last data row on page 2 (which is index 14 in the full dataset)
      const rows = screen.getAllByRole("row");
      const lastRowOnPage2 = rows[rows.length - 1];
      await user.click(lastRowOnPage2);

      // Dialog should show data for query 14
      // Query 14 has rows_returned of 100 + 14 = 114 for left, 99 + 14 = 113 for right
      expect(screen.getByText("Statistics")).toBeInTheDocument();
      expect(screen.getByText("114")).toBeInTheDocument(); // left rows_returned for query 14
      expect(screen.getByText("113")).toBeInTheDocument(); // right rows_returned for query 14
    });
  });

  describe("basic functionality", () => {
    it("renders empty state when no data", () => {
      render(
        <MacroQueriesCompareTable
          columns={columns}
          data={[]}
          filterConfigs={filterConfigs}
        />
      );

      expect(screen.getByText("No results.")).toBeInTheDocument();
    });

    it("renders table with data", () => {
      const mockData = Array.from({ length: 3 }, (_, i) => createMockQuery(i));

      render(
        <MacroQueriesCompareTable
          columns={columns}
          data={mockData}
          filterConfigs={filterConfigs}
        />
      );

      // Check that queries are displayed
      expect(screen.getByText("SELECT * FROM table0")).toBeInTheDocument();
      expect(screen.getByText("SELECT * FROM table1")).toBeInTheDocument();
      expect(screen.getByText("SELECT * FROM table2")).toBeInTheDocument();
    });

    it("displays execution time diff badges", () => {
      const mockData = [createMockQuery(0), createMockQuery(5)];

      render(
        <MacroQueriesCompareTable
          columns={columns}
          data={mockData}
          filterConfigs={filterConfigs}
        />
      );

      // Check that diff percentages are displayed
      expect(screen.getByText("0%")).toBeInTheDocument(); // query 0: 0 * 5 = 0%
      expect(screen.getByText("25%")).toBeInTheDocument(); // query 5: 5 * 5 = 25%
    });
  });

  describe("sorting", () => {
    it("displays correct query details after sorting", async () => {
      const user = userEvent.setup();
      // Create queries with different exec_time_diff values
      const mockData = [
        createMockQuery(5), // exec_time_diff: 25
        createMockQuery(1), // exec_time_diff: 5
        createMockQuery(10), // exec_time_diff: 50
      ];

      render(
        <MacroQueriesCompareTable
          columns={columns}
          data={mockData}
          filterConfigs={filterConfigs}
        />
      );

      // Sort by execution time (ascending)
      const sortButton = screen.getByRole("button", { name: /execution time/i });
      await user.click(sortButton);

      // After sorting ascending, first row should be query with exec_time_diff = 5 (index 1)
      const rows = screen.getAllByRole("row");
      const firstRowAfterSort = rows[1];
      await user.click(firstRowAfterSort);

      // Should show data for query 1 (rows_returned: 101 left, 100 right)
      expect(screen.getByText("Statistics")).toBeInTheDocument();
      expect(screen.getByText("101")).toBeInTheDocument();
      expect(screen.getByText("100")).toBeInTheDocument();
    });
  });

  describe("filtering", () => {
    it("displays correct query details after filtering", async () => {
      const user = userEvent.setup();
      const mockData = [
        createMockQuery(0),
        createMockQuery(5),
        createMockQuery(10),
      ];

      render(
        <MacroQueriesCompareTable
          columns={columns}
          data={mockData}
          filterConfigs={filterConfigs}
        />
      );

      // Filter by query containing "table10"
      const searchInput = screen.getByPlaceholderText(/filter executions/i);
      await user.type(searchInput, "table10");

      // Should only show query 10
      const rows = screen.getAllByRole("row");
      expect(rows.length).toBe(2); // header + 1 data row

      // Click on the filtered row
      const filteredRow = rows[1];
      await user.click(filteredRow);

      // Should show data for query 10 (rows_returned: 110 left, 109 right)
      expect(screen.getByText("Statistics")).toBeInTheDocument();
      expect(screen.getByText("110")).toBeInTheDocument();
      expect(screen.getByText("109")).toBeInTheDocument();
    });
  });
});
