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
import { Dialog } from "@/components/ui/dialog";
import MacroQueriesCompareDialog from "./MacroQueriesCompareDialog";
import { MacroQueriesPlan } from "./Columns";

// Helper function to create mock data
const createMockData = (
  overrides?: Partial<MacroQueriesPlan>
): MacroQueriesPlan => {
  return {
    key: "SELECT * FROM users",
    exec_time_diff: 10,
    exec_count_diff: 5,
    errors_diff: 0,
    rows_returned_diff: -2,
    same_plan: true,
    left: {
      key: "SELECT * FROM users",
      value: {
        query_type: "SELECT",
        original: "SELECT * FROM users WHERE id = ?",
        instructions: '{"opcode": "Route", "table": "users"}',
        exec_count: 100,
        exec_time: 50,
        shard_queries: 1,
        rows_returned: 100,
        rows_affected: 0,
        errors: 0,
        tables_used: "users",
      },
    },
    right: {
      key: "SELECT * FROM users",
      value: {
        query_type: "SELECT",
        original: "SELECT * FROM users WHERE id = ?",
        instructions: '{"opcode": "Route", "table": "users", "optimized": true}',
        exec_count: 105,
        exec_time: 55,
        shard_queries: 1,
        rows_returned: 98,
        rows_affected: 0,
        errors: 0,
        tables_used: "users",
      },
    },
    ...overrides,
  };
};

// Helper to render dialog with proper context
const renderDialog = (data: MacroQueriesPlan) => {
  return render(
    <Dialog open={true}>
      <MacroQueriesCompareDialog data={data} />
    </Dialog>
  );
};

describe("MacroQueriesCompareDialog", () => {
  it("renders statistics table with both left and right data", () => {
    const mockData = createMockData();

    renderDialog(mockData);

    // Check that the statistics table is rendered
    expect(screen.getByText("Statistics")).toBeInTheDocument();
    expect(screen.getByText("Execution Time")).toBeInTheDocument();
    expect(screen.getByText("Rows Returned")).toBeInTheDocument();
    expect(screen.getByText("Errors")).toBeInTheDocument();

    // Check that values are displayed
    expect(screen.getByText("50")).toBeInTheDocument(); // left exec_time
    expect(screen.getByText("55")).toBeInTheDocument(); // right exec_time
    expect(screen.getByText("100")).toBeInTheDocument(); // left rows_returned
    expect(screen.getByText("98")).toBeInTheDocument(); // right rows_returned
  });

  it("renders both query plans when both left and right exist", () => {
    const mockData = createMockData();

    renderDialog(mockData);

    // Check that both query plan sections are rendered
    expect(screen.getByText("Old query plan")).toBeInTheDocument();
    expect(screen.getByText("New query plan")).toBeInTheDocument();
  });

  it("displays N/A when left data is null", () => {
    const mockData = createMockData({ left: null });

    renderDialog(mockData);

    // Check that N/A is displayed in the Old column
    const cells = screen.getAllByText("N/A");
    expect(cells.length).toBeGreaterThanOrEqual(3); // At least 3 N/A for exec_time, rows_returned, errors

    // Check that right values are still displayed
    expect(screen.getByText("55")).toBeInTheDocument();
    expect(screen.getByText("98")).toBeInTheDocument();
  });

  it("does not render old query plan section when left is null", () => {
    const mockData = createMockData({ left: null });

    renderDialog(mockData);

    // Old query plan section should not be rendered
    expect(screen.queryByText("Old query plan")).not.toBeInTheDocument();

    // New query plan should still be rendered
    expect(screen.getByText("New query plan")).toBeInTheDocument();
  });

  it("displays N/A when right data is null", () => {
    const mockData = createMockData({ right: null });

    renderDialog(mockData);

    // Check that N/A is displayed in the New column
    const cells = screen.getAllByText("N/A");
    expect(cells.length).toBeGreaterThanOrEqual(3); // At least 3 N/A for exec_time, rows_returned, errors

    // Check that left values are still displayed
    expect(screen.getByText("50")).toBeInTheDocument();
    expect(screen.getByText("100")).toBeInTheDocument();
  });

  it("does not render new query plan section when right is null", () => {
    const mockData = createMockData({ right: null });

    renderDialog(mockData);

    // New query plan section should not be rendered
    expect(screen.queryByText("New query plan")).not.toBeInTheDocument();

    // Old query plan should still be rendered
    expect(screen.getByText("Old query plan")).toBeInTheDocument();
  });

  it("displays N/A for all values when both left and right are null", () => {
    const mockData = createMockData({ left: null, right: null });

    renderDialog(mockData);

    // Check that all values show N/A
    const cells = screen.getAllByText("N/A");
    expect(cells.length).toBeGreaterThanOrEqual(6); // 6 N/A cells (3 rows × 2 columns)

    // Neither query plan section should be rendered
    expect(screen.queryByText("Old query plan")).not.toBeInTheDocument();
    expect(screen.queryByText("New query plan")).not.toBeInTheDocument();
  });

  it("displays diff values correctly", () => {
    const mockData = createMockData({
      exec_time_diff: 15,
      rows_returned_diff: -5,
      errors_diff: 2,
    });

    renderDialog(mockData);

    // Check that diff values are displayed
    expect(screen.getByText("15")).toBeInTheDocument();
    expect(screen.getByText("-5")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
  });

  it("handles zero values in statistics", () => {
    const mockData = createMockData();
    if (mockData.left && mockData.right) {
      mockData.left.value.errors = 0;
      mockData.right.value.errors = 0;
    }

    renderDialog(mockData);

    // Check that zero values are displayed (there should be at least two 0s for errors)
    const zeroCells = screen.getAllByText("0");
    expect(zeroCells.length).toBeGreaterThanOrEqual(2);
  });
});
