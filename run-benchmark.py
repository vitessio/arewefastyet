# ------------------------------------------------------------------------------------------------------------------------------------
# Copyright 2020 The Vitess Authors.
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
#   - Creates Run ID and runs Benchmark
# -------------------------------------------------------------------------------------------------------------------------------------

import time
import subprocess
import os
from config import inventory_file
from pathlib import Path
import uuid 
import sys

# ------------------------------------------------------ Runs benchmark tasks ---------------------------------------------------------

def tasks():
   print('------------- Initialize VPS ------------------')

   run_id = uuid.uuid4()
   commit = sys.argv[1]

   os.system('python initialize_benchmark.py '+ str(run_id) + ' ' + commit)
   print('------------- Running Benchamrk ------------------')
   os.system('./run '+ Path('./ansible/' + inventory_file()).stem + '-' + str(run_id) + '.yml')
   print('------------- Adding results to the database ------------------')
   os.system('python report.py ' + str(run_id))

# -------------------------------------------------------------------------------------------------------------------------------------

tasks()
