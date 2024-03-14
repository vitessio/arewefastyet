/*
 *
 * Copyright 2022 The Vitess Authors.
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

package ansible

// This list of constants are the key used to name the variables that are given
// to the Ansible runtime by the Executor.
// These keys are associated with a value that will be used by Ansible to run
// the roles and tasks.
// When executing Ansible, those variables are passed using the `--extra-vars`
// Ansible flag.
// These keys and their values are meant to be stored in Config.ExtraVars.
const (

	// Arewefastyet related keys

	// KeyExecUUID corresponding value in the map is the Execution UUID of the
	// current benchmark.
	KeyExecUUID = "arewefastyet_exec_uuid"

	// KeyExecutionType corresponding value in the map is the type of execution for
	// this benchmark.
	KeyExecutionType = "arewefastyet_execution_type"

	// KeyBenchmarkConfigPath corresponding value in the map is the path to the configuration
	// file that will be used to execute the benchmark by arewefastyet inside the benchmarking
	// server.
	KeyBenchmarkConfigPath = "arewefastyet_configuration_file_path"

	// KeyBenchmarkSecretsPath corresponding value in the map is the path to the secrets
	// file that will be used to execute the benchmark by arewefastyet inside the benchmarking
	// server.
	KeyBenchmarkSecretsPath = "arewefastyet_secrets_file_path"

	// KeyLastIsSame is set to true if the previous execution has the same config as the current one.
	KeyLastIsSame = "arewefastyet_last_exec_is_same"

	// KeyNextIsSame is set to true if the next execution has the same config as the current one.
	KeyNextIsSame = "arewefastyet_next_exec_is_same"

	// Vitess related keys

	// KeyVitessSchema is the path to the Vitess VSchema that will be used for this benchmark.
	KeyVitessSchema = "vitess_vschema_path"

	// KeyVitessVersion corresponding value in the map is the git reference of SHA
	// which benchmarks will be executed.
	KeyVitessVersion = "vitess_git_version"

	// KeyVitessVersionFetchPR corresponding value in the map is the git reference
	// that is used to fetch a Pull Request.
	KeyVitessVersionFetchPR = "vitess_git_version_fetch_pr"

	// KeyVitessVersionPRNumber corresponding value in the map is the Pull Request
	// number that the execution will benchmark.
	KeyVitessVersionPRNumber = "vitess_git_version_pr_nb"

	// KeyVitessVersionName corresponding value in the map is the name of the vitess
	// version on which the benchmarks will be executed. For instance: 'latest', '14',
	// '13', ...
	KeyVitessVersionName = "vitess_version_name"

	// KeyVtgatePlanner corresponding value in the map is the query planner version
	// that will be used to execute the benchmark.
	KeyVtgatePlanner = "vitess_planner_version"

	// KeyExtraFlagsVTGate represents the list of VTGate flag that will be passed down
	// to VTGate during startup. The flags are represented as follows in the string:
	// 		"--flag1 --flag2"
	KeyExtraFlagsVTGate = "extra_vtgate_flags"

	// KeyExtraFlagsVTTablet represents the list of VTTablet flag that will be passed down
	// to VTTablet during startup. The flags are represented as follows in the string:
	// 		"--flag1 --flag2"
	KeyExtraFlagsVTTablet = "extra_vttablet_flags"

	// Runtime related keys

	// KeyGoVersion corresponding value in the map is the golang version to use for
	// the execution of the benchmark.
	KeyGoVersion = "golang_gover"

	// Stats database related keys

	// KeyStatsDBHost corresponding value in the map is the hostname for the stats
	// database.
	KeyStatsDBHost = "stats_remote_db_host"

	// KeyStatsDBName corresponding value in the map is the database name of the stats
	// database.
	KeyStatsDBName = "stats_remote_db_database"

	// KeyStatsDBPort corresponding value in the map is the port used to connect to the
	// stats database.
	KeyStatsDBPort = "stats_remote_db_port"

	// KeyStatsDBUser corresponding value in the map is the user that will be connected
	// to the stats database.
	KeyStatsDBUser = "stats_remote_db_user"

	// KeyStatsDBPassword corresponding value in the map is the password used to login the
	// user to the stats database.
	KeyStatsDBPassword = "stats_remote_db_password"
)
