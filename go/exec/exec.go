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

package exec

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/vitessio/arewefastyet/go/storage"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/exec/stats"
	"github.com/vitessio/arewefastyet/go/infra/ansible"
)

const (
	// keyExecUUID is the name of the key passed to each Ansible playbook
	// the value of the key points to an Execution UUID.
	keyExecUUID = "arewefastyet_exec_uuid"

	// keyExecSource is the name of the key that stores the name of the
	// execution's trigger.
	keyExecSource = "arewefastyet_source"

	// keyVitessVersion is the name of the key that stores the git reference
	// or SHA on which benchmarks will be executed.
	keyVitessVersion = "vitess_git_version"

	keyExecutionType = "arewefastyet_execution_type"

	keyVtgatePlanner = "planner_version"

	// keyGoVersion defines the golang version to use for the execution.
	keyGoVersion = "golang_gover"

	stderrFile = "exec-stderr.log"
	stdoutFile = "exec-stdout.log"

	ErrorNotPrepared      = "exec is not prepared"
	ErrorExecutionTimeout = "execution timeout"
)

type Exec struct {
	UUID          uuid.UUID
	AnsibleConfig ansible.Config
	Source        string
	GitRef        string

	// Status defines the status of the execution (canceled, finished, failed, etc)
	Status string

	StartedAt  *time.Time
	FinishedAt *time.Time

	// Defines the type of execution (oltp, tpcc, micro, ...)
	TypeOf string

	// PullNB defines the pull request number linked to this execution.
	PullNB int

	// Configuration used to interact with the SQL database.
	configDB *psdb.Config

	// Client to communicate with the SQL database.
	clientDB *psdb.Client

	// Configuration used to authenticate and insert execution stats
	// data to a remote database system.
	statsRemoteDBConfig stats.RemoteDBConfig

	// rootDir represents the parent directory of the Exec.
	// From there, the Exec's unique directory named Exec.dirPath will
	// be created once Exec.Prepare is called.
	rootDir string

	// dirPath is Exec's unique directory where all reports, directories,
	// files, and logs are kept.
	dirPath string

	stdout io.Writer
	stderr io.Writer

	createdInDB bool
	prepared    bool
	configPath  string

	// VtgatePlannerVersion is the planner version that vtgate is going to use
	VtgatePlannerVersion string

	// GolangVersion is the go version to use while executing the benchmark on the remote host
	GolangVersion string

	// ServerAddress is the IP address on which the benchmark will be executed.
	ServerAddress string
}

const (
	SourceCron            = "cron"
	SourcePullRequest     = "cron_pr"
	SourcePullRequestBase = "cron_pr_base"
	SourceTag             = "cron_tags_"
	SourceReleaseBranch   = "cron_"
)

// SetStdout sets the standard output of Exec.
func (e *Exec) SetStdout(stdout *os.File) {
	e.stdout = stdout
	e.AnsibleConfig.SetStdout(stdout)
}

// SetStderr sets the standard error output of Exec.
func (e *Exec) SetStderr(stderr *os.File) {
	e.stderr = stderr
	e.AnsibleConfig.SetStderr(stderr)
}

// SetOutputToDefaultPath sets Exec's outputs to their default files (stdoutFile and
// stderrFile). If they can't be found in Exec.dirPath, they will be created.
func (e *Exec) SetOutputToDefaultPath() error {
	if !e.prepared {
		return errors.New(ErrorNotPrepared)
	}
	outFile, err := os.OpenFile(path.Join(e.dirPath, stdoutFile), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	errFile, err := os.OpenFile(path.Join(e.dirPath, stderrFile), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	e.stdout = outFile
	e.stderr = errFile
	e.AnsibleConfig.SetOutputs(outFile, errFile)
	return nil
}

// Prepare prepares the Exec for a future Execution.
func (e *Exec) Prepare() error {
	// Returns if the execution is already prepared
	if e.prepared {
		return nil
	}

	var err error
	defer func() {
		if !e.createdInDB {
			return
		}
		e.handleStepEnd(err)
	}()

	e.clientDB, err = e.configDB.NewClient()
	if err != nil {
		return err
	}

	// insert new exec in SQL
	if _, err = e.clientDB.Insert(
		"INSERT INTO execution(uuid, status, source, git_ref, type, pull_nb, go_version) VALUES(?, ?, ?, ?, ?, ?, ?)",
		e.UUID.String(),
		StatusCreated,
		e.Source,
		e.GitRef,
		e.TypeOf,
		e.PullNB,
		e.GolangVersion,
	); err != nil {
		return err
	}
	e.createdInDB = true

	err = e.prepareDirectories()
	if err != nil {
		return err
	}

	if e.configPath == "" {
		e.configPath = viper.ConfigFileUsed()
	}
	e.AnsibleConfig.ExtraVars = map[string]interface{}{}
	e.statsRemoteDBConfig.AddToAnsible(&e.AnsibleConfig)
	if e.PullNB != 0 {
		e.AnsibleConfig.ExtraVars["vitess_git_version_fetch_pr"] = "refs/pull/" + strconv.Itoa(e.PullNB) + "/head"
		e.AnsibleConfig.ExtraVars["vitess_git_version_pr_nb"] = e.PullNB
	}

	// Enable schema tracking only if we execute macrobenchmark main CRONs
	if e.Source == SourceCron && e.TypeOf != "micro" {
		e.AnsibleConfig.ExtraVars["vitess_schema_tracking"] = 1
	}

	e.prepared = true
	return nil
}

// ExecuteWithTimeout will call execution's Execute method with the given timeout.
func (e Exec) ExecuteWithTimeout(timeout time.Duration) (err error) {
	defer func() {
		e.handleStepEnd(err)
	}()
	errs := make(chan error)

	go func() {
		errs <- e.Execute()
	}()

	select {
	case err = <-errs:
		return
	case <-time.After(timeout):
		err = errors.New(ErrorExecutionTimeout)
		return
	}
}

// Execute will provision infra, configure Ansible files, and run the given Ansible config.
func (e *Exec) Execute() (err error) {
	defer func() {
		e.handleStepEnd(err)
	}()

	if !e.prepared {
		return errors.New(ErrorNotPrepared)
	}
	if _, err := e.clientDB.Insert("UPDATE execution SET started_at = CURRENT_TIME, status = ? WHERE uuid = ?", StatusStarted, e.UUID.String()); err != nil {
		return err
	}

	// TODO: optimize tokenization of Ansible files.
	err = ansible.AddIPsToFiles([]string{e.ServerAddress}, e.AnsibleConfig)
	if err != nil {
		return err
	}
	err = ansible.AddLocalConfigPathToFiles(e.configPath, e.AnsibleConfig)
	if err != nil {
		return err
	}

	e.prepareAnsibleForExecution()

	// Run the given config on Ansible
	err = ansible.Run(&e.AnsibleConfig)
	if err != nil {
		return err
	}
	return nil
}

func (e *Exec) prepareAnsibleForExecution() {
	e.AnsibleConfig.ExtraVars[keyExecUUID] = e.UUID.String()
	e.AnsibleConfig.ExtraVars[keyVitessVersion] = e.GitRef
	e.AnsibleConfig.ExtraVars[keyExecSource] = e.Source
	e.AnsibleConfig.ExtraVars[keyExecutionType] = e.TypeOf
	e.AnsibleConfig.ExtraVars[keyGoVersion] = e.GolangVersion

	// not adding the -planner_version flag to ansible if we did not specify it or if using the default value
	if e.VtgatePlannerVersion == string(macrobench.Gen4FallbackPlanner) {
		e.AnsibleConfig.ExtraVars[keyVtgatePlanner] = e.VtgatePlannerVersion
	}
}

func (e *Exec) Success() error {
	// checking if the execution has not already failed
	rows, err := e.clientDB.Select("SELECT uuid FROM execution WHERE uuid = ? AND status = ?", e.UUID.String(), StatusFailed)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return nil
	}
	_, err = e.clientDB.Insert("UPDATE execution SET finished_at = CURRENT_TIME, status = ? WHERE uuid = ?", StatusFinished, e.UUID.String())
	return err
}

func (e *Exec) handleStepEnd(err error) {
	if err != nil {
		_, _ = e.clientDB.Insert("UPDATE execution SET finished_at = CURRENT_TIME, status = ? WHERE uuid = ?", StatusFailed, e.UUID.String())
	}
}

// NewExec creates a new *Exec with an autogenerated uuid.UUID as well
// as a constructed infra.Infra.
func NewExec() (*Exec, error) {
	ex := Exec{
		UUID:  uuid.New(),

		// By default Exec prints os.Stdout and os.Stderr.
		// This can be changed later by explicitly using
		// exec.SetStdout and exec.SetStderr, or SetOutputToDefaultPath.
		stdout: os.Stdout,
		stderr: os.Stderr,

		configDB:   &psdb.Config{},
		clientDB:   nil,
		configPath: viper.ConfigFileUsed(),
	}
	return &ex, nil
}

// NewExecWithConfig will create a new Exec using the NewExec method, and will
// use viper.Viper to apply the configuration located at pathConfig.
func NewExecWithConfig(pathConfig string) (*Exec, error) {
	e, err := NewExec()
	if err != nil {
		return nil, err
	}
	v := viper.New()

	v.SetConfigFile(pathConfig)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	err = e.AddToViper(v)
	if err != nil {
		return nil, err
	}
	e.configPath = pathConfig
	return e, nil
}

func GetRecentExecutions(client storage.SQLClient) ([]*Exec, error) {
	var res []*Exec
	query := "SELECT uuid, status, git_ref, started_at, finished_at, source, type, pull_nb, go_version FROM execution ORDER BY started_at DESC LIMIT 50"
	result, err := client.Select(query)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var eUUID string
		exec := &Exec{}
		err = result.Scan(&eUUID, &exec.Status, &exec.GitRef, &exec.StartedAt, &exec.FinishedAt, &exec.Source, &exec.TypeOf, &exec.PullNB, &exec.GolangVersion)
		if err != nil {
			return nil, err
		}
		exec.UUID, err = uuid.Parse(eUUID)
		if err != nil {
			return nil, err
		}
		if exec.TypeOf != "micro" {
			macroResult, err := client.Select("SELECT m.vtgate_planner_version FROM macrobenchmark m, execution e WHERE e.uuid = m.exec_uuid AND e.uuid = ? LIMIT 1", eUUID)
			if err != nil {
				return nil, err
			}
			defer macroResult.Close()

			var plannerVersion string
			if macroResult.Next() {
				err = macroResult.Scan(&plannerVersion)
				if err != nil {
					return nil, err
				}
			}
			exec.VtgatePlannerVersion = plannerVersion
		}
		res = append(res, exec)
	}
	return res, nil
}

func GetFinishedExecution(client storage.SQLClient, gitRef, source, benchmarkType, plannerVersion string, pullNb int) (string, error) {
	var eUUID string
	var result *sql.Rows
	var err error
	query := ""
	if plannerVersion == "" {
		// no plannerVersion, meaning we are dealing with a micro benchmark
		query = "SELECT e.uuid FROM execution e WHERE e.source = ? AND e.status = ? AND e.type = ? AND e.git_ref = ? AND e.pull_nb = ? ORDER BY e.finished_at DESC LIMIT 1"
		result, err = client.Select(query, source, StatusFinished, benchmarkType, gitRef, pullNb)
	} else {
		// we have a plannerVersion, meaning we are dealing with a macro benchmark
		query = "SELECT e.uuid FROM execution e, macrobenchmark m WHERE e.uuid = m.exec_uuid AND m.vtgate_planner_version = ? AND e.source = ? AND e.status = ? AND e.type = ? AND e.git_ref = ? AND e.pull_nb = ? ORDER BY e.finished_at DESC LIMIT 1"
		result, err = client.Select(query, plannerVersion, source, StatusFinished, benchmarkType, gitRef, pullNb)
	}
	if err != nil {
		return "", err
	}
	defer result.Close()
	for result.Next() {
		err = result.Scan(&eUUID)
		if err != nil {
			return "", err
		}
	}
	return eUUID, nil
}

// GetPreviousFromSourceMicrobenchmark gets the previous execution from the same source for microbenchmarks
func GetPreviousFromSourceMicrobenchmark(client storage.SQLClient, source, gitRef string) (execUUID, gitRefOut string, err error) {
	query := "SELECT e.uuid, e.git_ref FROM execution e WHERE e.source = ? AND e.status = 'finished' AND " +
		"e.type = \"micro\" AND e.git_ref != ? ORDER BY e.started_at DESC LIMIT 1"
	result, err := client.Select(query, source, gitRef)
	if err != nil {
		return
	}
	defer result.Close()
	for result.Next() {
		err = result.Scan(&execUUID, &gitRefOut)
		if err != nil {
			return
		}
	}
	return
}

// GetPreviousFromSourceMacrobenchmark gets the previous execution from the same source with the sane plannerVersion for macrobenchmarks
func GetPreviousFromSourceMacrobenchmark(client storage.SQLClient, source, typeOf, plannerVersion, gitRef string) (execUUID, gitRefOut string, err error) {
	query := "SELECT e.uuid, e.git_ref FROM execution e, macrobenchmark m WHERE e.source = ? AND e.status = 'finished' AND " +
		"e.type = ? AND e.git_ref != ? AND m.exec_uuid = e.uuid AND m.vtgate_planner_version = ? ORDER BY e.started_at DESC LIMIT 1"
	result, err := client.Select(query, source, typeOf, gitRef, plannerVersion)
	if err != nil {
		return
	}
	defer result.Close()
	for result.Next() {
		err = result.Scan(&execUUID, &gitRefOut)
		if err != nil {
			return
		}
	}
	return
}

// GetLatestCronJobForMicrobenchmarks will fetch and return the commit sha for which
// the last cron job for microbenchmarks was run
func GetLatestCronJobForMicrobenchmarks(client storage.SQLClient) (gitSha string, err error) {
	query := "select git_ref from execution where source = \"cron\" and status = \"finished\" and type = \"micro\" order by started_at desc limit 1"
	rows, err := client.Select(query)
	if err != nil {
		return "", err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&gitSha)
		return gitSha, err
	}
	return "", nil
}

// GetLatestCronJobForMacrobenchmarks will fetch and return the commit sha for which
// the last cron job for macrobenchmarks was run
func GetLatestCronJobForMacrobenchmarks(client storage.SQLClient) (gitSha string, err error) {
	query := "select git_ref from execution where source = \"cron\" and status = \"finished\" and ( type = \"oltp\" or type = \"tpcc\" ) order by started_at desc limit 1"
	rows, err := client.Select(query)
	if err != nil {
		return "", err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&gitSha)
		return gitSha, err
	}
	return "", nil
}

func Exists(client storage.SQLClient, gitRef, source, typeOf, status string) (bool, error) {
	query := "SELECT uuid FROM execution WHERE status = ? AND git_ref = ? AND type = ? AND source = ?"
	result, err := client.Select(query, status, gitRef, typeOf, source)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return result.Next(), nil
}

func ExistsMacrobenchmark(client storage.SQLClient, gitRef, source, typeOf, status, planner string) (bool, error) {
	query := "SELECT uuid FROM execution e, macrobenchmark m WHERE e.status = ? AND e.git_ref = ? AND e.type = ? AND e.source = ? AND m.vtgate_planner_version = ? AND e.uuid = m.exec_uuid"
	result, err := client.Select(query, status, gitRef, typeOf, source, planner)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return result.Next(), nil
}

func ExistsMacrobenchmarkStartedToday(client storage.SQLClient, gitRef, source, typeOf, planner, status string) (bool, error) {
	query := fmt.Sprintf("SELECT uuid FROM execution e, macrobenchmark m WHERE e.status = '%s' AND e.git_ref = ? AND e.type = ? AND e.source = ? AND m.vtgate_planner_version = ? AND e.uuid = m.exec_uuid AND e.started_at >= CURDATE()", status)
	result, err := client.Select(query, gitRef, typeOf, source, planner)
	if err != nil {
		return false, err
	}
	exists := result.Next()
	result.Close()
	if exists {
		return true, nil
	}
	query = fmt.Sprintf("SELECT uuid FROM execution e WHERE e.status = '%s' AND e.git_ref = ? AND e.type = ? AND e.source = ? AND e.started_at >= CURDATE()", status)
	result, err = client.Select(query, gitRef, typeOf, source)
	if err != nil {
		return false, err
	}
	defer result.Close()
	for result.Next() {
		var exec_uuid string
		err = result.Scan(&exec_uuid)
		if err != nil {
			return false, err
		}
		query = "SELECT exec_uuid FROM macrobenchmark m WHERE m.exec_uuid = ?"
		resultMacro, err := client.Select(query, exec_uuid)
		if err != nil {
			return false, err
		}
		defer resultMacro.Close()

		next := resultMacro.Next()
		resultMacro.Close()
		if !next {
			return true, nil
		}
	}
	return false, nil
}
