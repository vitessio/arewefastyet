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
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/vitessio/arewefastyet/go/storage"
	"github.com/vitessio/arewefastyet/go/storage/psdb"
	"github.com/vitessio/arewefastyet/go/tools/git"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/exec/stats"
	"github.com/vitessio/arewefastyet/go/infra/ansible"
)

const (
	stderrFile = "exec-stderr.log"
	stdoutFile = "exec-stdout.log"

	ErrorNotPrepared      = "exec is not prepared"
	ErrorExecutionTimeout = "execution timeout"
)

type Exec struct {
	UUID          uuid.UUID
	RawUUID       string
	AnsibleConfig ansible.Config
	Source        string
	GitRef        string
	VitessVersion git.Version

	// NextBenchmarkIsTheSame is set to true if the next benchmark has the same config
	// as the current one. This allows us to do some optimization in Ansible and speed
	// up the entire benchmarking process.
	NextBenchmarkIsTheSame bool

	// PreviousBenchmarkIsTheSame is set to true if the previous benchmark had the same
	// config as this one. It allows us to skip the preparatory cleanup of the server.
	PreviousBenchmarkIsTheSame bool

	// Status defines the status of the execution (canceled, finished, failed, etc)
	Status string

	StartedAt  *time.Time
	FinishedAt *time.Time

	// Defines the workload of execution (oltp, tpcc, micro, ...)
	Workload string

	// PullNB defines the pull request number linked to this execution.
	PullNB            int
	PullBaseBranchRef string

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
	secretsPath string

	// VtgatePlannerVersion is the planner version that vtgate is going to use
	VtgatePlannerVersion string

	// GolangVersion is the go version to use while executing the benchmark on the remote host
	GolangVersion string

	// ServerAddress is the IP address on which the benchmark will be executed.
	ServerAddress string

	RepoDir string

	// The configuration of the Vitess components (only vttablet and vtgate for now) can be
	// customized through the configuration file. Some additional flags can be passed down
	// to those two binaries.
	//
	// This is achieved by the 'exec-vitess-config' flag. This flag is a map that contains a map:
	// 		map[string]map[string]string
	// The outer map represents the different releases of Vitess and the inner map the different
	// binaries. Finally, the key of the inner map corresponds to the flags to add.
	// The key for the inner map being the Vitess version it has to be formatted using '-' to
	// separate the major/minor/patch numbers. Like so: 'v15.0.0' becomes: '15-0-0'.
	//
	// We can have as many releases as we want in the outer map. Only the one that is the closest
	// to this execution's version will be used. For instance, if we define two versions, v14.0.0
	// and v15.0.0, but the execution runs on Vitess v14.0.3, then the configuration for v14.0.0
	// will be used.
	//
	// Here is an example of how the flag should be formatted in YAML:
	//
	// 			exec-vitess-config:
	//  		  14: # will match >= v14.0.0
	//    		    vtgate: --my_flag=0
	//  		  15: # will match >= v15.0.0 and override the flags defined in < 15.0.0
	//    		    vtgate: --my_flag=1
	// 				vttablet: --my_flag=1
	//  		  15-0-1: # will match >= v15.0.1 and override the flags defined in < 15.0.1
	//    		    vtgate: --my_flag=3
	// 				vttablet: --my_flag=3 --custom-flag=on
	//
	// Finally, rawVitessConfig stores the raw data from viper.Viper which is then computed during
	// the execution's Prepare step to become the final version stored in vitessConfig. The value
	// in vitessConfig is then used when adding extra vars to Ansible.
	rawVitessConfig rawSingleVitessVersionConfig
	vitessConfig    vitessConfig

	vitessSchemaPath string
}

const (
	MaximumBenchmarkWithSameConfig = 10

	SourceCron            = "cron"
	SourcePullRequest     = "cron_pr"
	SourcePullRequestBase = "cron_pr_base"
	SourceTag             = "cron_tags_"
	SourceReleaseBranch   = "cron_"
)

// NewExec creates a new *Exec given the string representation of an uuid.UUID.
// If no UUID is provided, a new one will be generated.
func NewExec(uuidRaw string) (*Exec, error) {
	if uuidRaw == "" {
		uuidRaw = uuid.NewString()
	}
	parsedUUID, err := uuid.Parse(uuidRaw)
	if err != nil {
		return nil, err
	}

	ex := Exec{
		UUID: parsedUUID,

		// By default Exec prints os.Stdout and os.Stderr.
		// This can be changed later by explicitly using
		// exec.SetStdout and exec.SetStderr, or SetOutputToDefaultPath.
		stdout: os.Stdout,
		stderr: os.Stderr,

		configDB:      &psdb.Config{},
		clientDB:      nil,
		configPath:    viper.ConfigFileUsed(),
		AnsibleConfig: ansible.NewConfig(),
	}
	return &ex, nil
}

// NewExecWithConfig will create a new Exec using the NewExec method, and will
// use viper.Viper to apply the configuration located at pathConfig.
func NewExecWithConfig(path, uuid string) (*Exec, error) {
	e, err := NewExec(uuid)
	if err != nil {
		return nil, err
	}

	nv := viper.New()
	err = nv.MergeConfigMap(viper.AllSettings())
	if err != nil {
		return nil, err
	}
	nv.SetConfigFile(path)
	err = nv.MergeInConfig()
	if err != nil {
		return nil, err
	}

	err = e.AddToViper(nv)
	if err != nil {
		return nil, err
	}
	e.configPath = path
	return e, nil
}

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
	if _, err = e.clientDB.Write(
		"INSERT INTO execution(uuid, status, source, git_ref, workload, pull_nb, go_version) VALUES(?, ?, ?, ?, ?, ?, ?)",
		e.UUID.String(),
		StatusCreated,
		e.Source,
		e.GitRef,
		e.Workload,
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
	if e.secretsPath == "" {
		e.secretsPath = viper.GetString("secrets")
	}

	// TODO: optimize tokenization of Ansible files.
	err = ansible.AddIPsToFiles([]string{e.ServerAddress}, e.AnsibleConfig)
	if err != nil {
		return err
	}

	err = prepareVitessConfiguration(e.rawVitessConfig, e.VitessVersion, &e.vitessConfig)
	if err != nil {
		return err
	}

	err = e.prepareAnsibleForExecution()
	if err != nil {
		return err
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
	if _, err := e.clientDB.Write("UPDATE execution SET started_at = CURRENT_TIME, status = ? WHERE uuid = ?", StatusStarted, e.UUID.String()); err != nil {
		return err
	}

	// Run the given config on Ansible
	err = ansible.Run(&e.AnsibleConfig)
	if err != nil {
		return err
	}
	return nil
}

// prepareAnsibleForExecution adds all the required values to run the benchmark with Ansible.
// These values are stored using a key/value map.
func (e *Exec) prepareAnsibleForExecution() error {
	// arewefastyet related values
	absConfigPath, err := filepath.Abs(e.configPath)
	if err != nil {
		return err
	}
	e.AnsibleConfig.AddExtraVar(ansible.KeyBenchmarkConfigPath, absConfigPath)

	absSecretsPath, err := filepath.Abs(e.secretsPath)
	if err != nil {
		return err
	}
	e.AnsibleConfig.AddExtraVar(ansible.KeyBenchmarkSecretsPath, absSecretsPath)
	e.AnsibleConfig.AddExtraVar(ansible.KeyExecUUID, e.UUID.String())
	e.AnsibleConfig.AddExtraVar(ansible.KeyExecutionWorkload, e.Workload)

	if e.PreviousBenchmarkIsTheSame {
		e.AnsibleConfig.AddExtraVar(ansible.KeyLastIsSame, true)
	}
	if e.NextBenchmarkIsTheSame {
		e.AnsibleConfig.AddExtraVar(ansible.KeyNextIsSame, true)
	}

	// vitess related values
	e.AnsibleConfig.AddExtraVar(ansible.KeyVitessVersion, e.GitRef)
	if e.PullNB != 0 {
		e.AnsibleConfig.AddExtraVar(ansible.KeyVitessVersionFetchPR, "refs/pull/"+strconv.Itoa(e.PullNB)+"/head")
		e.AnsibleConfig.AddExtraVar(ansible.KeyVitessVersionPRNumber, e.PullNB)
	}
	e.AnsibleConfig.AddExtraVar(ansible.KeyVtgatePlanner, e.VtgatePlannerVersion)
	e.AnsibleConfig.AddExtraVar(ansible.KeyVitessMajorVersion, e.VitessVersion.Major)
	e.AnsibleConfig.AddExtraVar(ansible.KeyVitessSchema, e.vitessSchemaPath)
	e.vitessConfig.addToAnsible(&e.AnsibleConfig)

	// runtime related values
	e.AnsibleConfig.AddExtraVar(ansible.KeyGoVersion, e.GolangVersion)

	// stats database related values
	e.statsRemoteDBConfig.AddToAnsible(&e.AnsibleConfig)
	return nil
}

func (e *Exec) Success() error {
	// checking if the execution has not already failed
	rows, err := e.clientDB.Read("SELECT uuid FROM execution WHERE uuid = ? AND status = ?", e.UUID.String(), StatusFailed)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return nil
	}
	_, err = e.clientDB.Write("UPDATE execution SET finished_at = CURRENT_TIME, status = ? WHERE uuid = ?", StatusFinished, e.UUID.String())
	return err
}

func (e *Exec) handleStepEnd(err error) {
	if err != nil {
		_, _ = e.clientDB.Write("UPDATE execution SET finished_at = CURRENT_TIME, status = ? WHERE uuid = ?", StatusFailed, e.UUID.String())
	}
}

func GetRecentExecutions(client storage.SQLClient) ([]*Exec, error) {
	var res []*Exec
	query := "SELECT uuid, status, git_ref, started_at, finished_at, source, workload, pull_nb, go_version FROM execution ORDER BY started_at DESC LIMIT 300"
	result, err := client.Read(query)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		exec := &Exec{}
		err = result.Scan(&exec.RawUUID, &exec.Status, &exec.GitRef, &exec.StartedAt, &exec.FinishedAt, &exec.Source, &exec.Workload, &exec.PullNB, &exec.GolangVersion)
		if err != nil {
			return nil, err
		}
		res = append(res, exec)
	}
	return res, nil
}

func GetFinishedExecution(client storage.SQLClient, gitRef, source, workload, plannerVersion string, pullNb int) (string, error) {
	var eUUID string
	var result *sql.Rows
	var err error
	query := ""
	if plannerVersion == "" {
		// no plannerVersion, meaning we are dealing with a micro benchmark
		query = "SELECT e.uuid FROM execution e WHERE e.source = ? AND e.status = ? AND e.workload = ? AND e.git_ref = ? AND e.pull_nb = ? ORDER BY e.finished_at DESC LIMIT 1"
		result, err = client.Read(query, source, StatusFinished, workload, gitRef, pullNb)
	} else {
		// we have a plannerVersion, meaning we are dealing with a macro benchmark
		query = "SELECT e.uuid FROM execution e, macrobenchmark m WHERE e.uuid = m.exec_uuid AND m.vtgate_planner_version = ? AND e.source = ? AND e.status = ? AND e.workload = ? AND e.git_ref = ? AND e.pull_nb = ? ORDER BY e.finished_at DESC LIMIT 1"
		result, err = client.Read(query, plannerVersion, source, StatusFinished, workload, gitRef, pullNb)
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

func IsLastExecutionFinished(client storage.SQLClient) (bool, error) {
	query := "SELECT e.status FROM execution e ORDER BY e.started_at DESC LIMIT 1"
	result, err := client.Read(query)
	if err != nil {
		return false, err
	}
	defer result.Close()
	var status string
	if result.Next() {
		err = result.Scan(&status)
		if err != nil {
			return false, err
		}
	}
	return status == StatusFinished, nil
}

// GetPreviousFromSourceMicrobenchmark gets the previous execution from the same source for microbenchmarks
func GetPreviousFromSourceMicrobenchmark(client storage.SQLClient, source, gitRef string) (execUUID, gitRefOut string, err error) {
	query := "SELECT e.uuid, e.git_ref FROM execution e WHERE e.source = ? AND e.status = 'finished' AND " +
		"e.workload = \"micro\" AND e.git_ref != ? ORDER BY e.started_at DESC LIMIT 1"
	result, err := client.Read(query, source, gitRef)
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
func GetPreviousFromSourceMacrobenchmark(client storage.SQLClient, source, workload, plannerVersion, gitRef string) (execUUID, gitRefOut string, err error) {
	query := "SELECT e.uuid, e.git_ref FROM execution e, macrobenchmark m WHERE e.source = ? AND e.status = 'finished' AND " +
		"e.workload = ? AND e.git_ref != ? AND m.exec_uuid = e.uuid AND m.vtgate_planner_version = ? ORDER BY e.started_at DESC LIMIT 1"
	result, err := client.Read(query, source, workload, gitRef, plannerVersion)
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

// GetLatestDailyJobForMicrobenchmarks will fetch and return the commit sha for which
// the last daily job for microbenchmarks was run
func GetLatestDailyJobForMicrobenchmarks(client storage.SQLClient) (gitSha string, err error) {
	query := "select git_ref from execution where source = \"cron\" and status = \"finished\" and workload = \"micro\" order by started_at desc limit 1"
	rows, err := client.Read(query)
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

// GetLatestDailyJobForMacrobenchmarks will fetch and return the commit sha for which
// the last daily job for macrobenchmarks was run
func GetLatestDailyJobForMacrobenchmarks(client storage.SQLClient) (gitSha string, err error) {
	query := "select git_ref from execution where source = \"cron\" and status = \"finished\" and ( workload != \"micro\" ) order by started_at desc limit 1"
	rows, err := client.Read(query)
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

func Exists(client storage.SQLClient, gitRef, source, workload, status string) (bool, error) {
	query := "SELECT uuid FROM execution WHERE status = ? AND git_ref = ? AND workload = ? AND source = ?"
	result, err := client.Read(query, status, gitRef, workload, source)
	if err != nil {
		return false, err
	}
	defer result.Close()
	return result.Next(), nil
}

func CountMacroBenchmark(client storage.SQLClient, gitRef, source, workload, status, planner string) (int, error) {
	query := "SELECT count(uuid) FROM execution e, macrobenchmark m WHERE e.status = ? AND e.git_ref = ? AND e.workload = ? AND e.source = ? AND m.vtgate_planner_version = ? AND e.uuid = m.exec_uuid"
	result, err := client.Read(query, status, gitRef, workload, source, planner)
	if err != nil {
		return 0, err
	}
	defer result.Close()
	var nb int
	if result.Next() {
		err = result.Scan(&nb)
		if err != nil {
			return 0, err
		}
	}
	return nb, nil
}

func DeleteExecution(client storage.SQLClient, gitRef, UUID, source string) error {
	query := fmt.Sprintf("DELETE FROM execution WHERE uuid LIKE '%%%s%%' AND git_ref LIKE '%%%s%%' AND source = '%s'", UUID, gitRef, source)
	_, err := client.Read(query)
	if err != nil {
		return err
	}
	return nil
}

type History struct {
	SHA                  string     `json:"sha"`
	Source               string     `json:"source"`
	WorkloadsBenchmarked int        `json:"workloads_benchmarked"`
	StartedAt            *time.Time `json:"started_at"`
}

func GetHistory(client storage.SQLClient) ([]*History, error) {
	query := `
			SELECT
				git_ref,
				source,
				COUNT(DISTINCT workload) AS distinct_workloads,
				MIN(min_started_at) AS min_started_at
			FROM (
				SELECT
					git_ref,
					source,
					workload,
					MIN(started_at) AS min_started_at
				FROM
					execution
				WHERE
					status = 'finished'
				GROUP BY
					git_ref,
					source,
					workload
				HAVING
					COUNT(*) >= 10
			) AS subquery
			GROUP BY
				git_ref,
				source
			ORDER BY
				min_started_at DESC;`

	result, err := client.Read(query)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	res := make([]*History, 0)
	for result.Next() {
		history := &History{}
		err = result.Scan(&history.SHA, &history.Source, &history.WorkloadsBenchmarked, &history.StartedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, history)
	}
	return res, nil
}
