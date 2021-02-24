# -------------------------------------------------------------------------------------------------------------------------------------------------
# Copyright 2021 The Vitess Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#    http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# demonstrates to:
#   - to call 'python run-benchmark.py' at specified time of the day
#
# Arguments: python scheduler.py <time>
# -----------------------------------------------------------------------------------------------------------------------------------------------------

import schedule
import time
import os
import uuid 
import sys

# ---------------------------------------------------------------- Calls run benchmark ----------------------------------------------------------------

def job():
    commit = 'HEAD'
    run_id = uuid.uuid4()
    os.system('python run-benchmark.py ' + commit + ' ' + str(run_id) + ' scheduler' + ' &')

# -----------------------------------------------------------------------------------------------------------------------------------------------------

# Runs everyday at <specified time>
schedule.every().day.at(sys.argv[1]).do(job)

while True:
    schedule.run_pending()
    time.sleep(1)
