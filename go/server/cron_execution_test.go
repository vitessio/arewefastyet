/*
 *
 * Copyright 2021 The Vitess Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

package server

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/uuid"
)

// TestRemoveFromQueueAfterUUIDMutation is a regression test for the queue-zombie
// bug: executeElement mutates element.identifier.UUID between retries, so removing
// the element by its (mutated) current identifier no longer matches the map key it
// was stored under, leaking the entry forever. removeFromQueue must be called with
// the original key so the entry is always removed.
func TestRemoveFromQueueAfterUUIDMutation(t *testing.T) {
	c := qt.New(t)

	queue = make(executionQueue)

	originalKey := executionIdentifier{
		GitRef:         "abc123",
		Source:         "cron_pr",
		Workload:       "tpcc-olap",
		PlannerVersion: "Gen4",
		PullNb:         42,
		UUID:           uuid.NewString(),
	}
	element := &executionQueueElement{
		identifier: originalKey,
		Executing:  true,
	}
	queue[originalKey] = element

	// Simulate what executeElement does on a retry: a fresh execution UUID is
	// assigned to the in-flight element, which changes element.identifier but must
	// not change the key it lives under in the queue map.
	element.identifier.UUID = uuid.NewString()

	// Removing by the mutated identifier (the old, buggy behaviour) would miss.
	delete(queue, element.identifier)
	c.Assert(queue, qt.HasLen, 1, qt.Commentf("deleting by the mutated identifier must not remove the entry"))

	// Removing by the original key (the fix) removes the entry.
	removeFromQueue(originalKey)
	c.Assert(queue, qt.HasLen, 0, qt.Commentf("removeFromQueue with the original key must remove the entry"))
}
