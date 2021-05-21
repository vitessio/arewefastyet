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
	"strconv"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/vitessio/arewefastyet/go/exec/stats"
	"github.com/vitessio/arewefastyet/go/infra"
	"github.com/vitessio/arewefastyet/go/infra/ansible"
	"github.com/vitessio/arewefastyet/go/infra/construct"
	"github.com/vitessio/arewefastyet/go/infra/equinix"
	"github.com/vitessio/arewefastyet/go/slack"
	"github.com/vitessio/arewefastyet/go/storage/influxdb"
	"github.com/vitessio/arewefastyet/go/storage/mysql"
	"github.com/vitessio/arewefastyet/go/tools/git"
	"github.com/vitessio/arewefastyet/go/tools/macrobench"
	"github.com/vitessio/arewefastyet/go/tools/microbench"
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

	stderrFile = "exec-stderr.log"
	stdoutFile = "exec-stdout.log"

	ErrorNotPrepared = "exec is not prepared"
)

type Exec struct {
	UUID          uuid.UUID
	InfraConfig   infra.Config
	AnsibleConfig ansible.Config
	Infra         infra.Infra
	Source        string
	GitRef        string

	// Defines the type of execution (oltp, tpcc, micro, ...)
	typeOf string

	// Defines the pull request number linked to this execution.
	pullNB int

	// Configuration used to interact with the SQL database.
	configDB *mysql.ConfigDB

	// Client to communicate with the SQL database.
	clientDB *mysql.Client

	// Configuration used to authenticate and insert execution stats
	// data to a remote database system.
	statsRemoteDBConfig stats.RemoteDBConfig

	// Configuration used to send message to Slack.
	slackConfig slack.Config

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

	e.clientDB, err = mysql.New(*e.configDB)
	if err != nil {
		return err
	}

	// insert new exec in SQL
	if _, err = e.clientDB.Insert("INSERT INTO execution(uuid, status, source, git_ref, type, pull_nb) VALUES(?, ?, ?, ?, ?, ?)", e.UUID.String(), StatusCreated, e.Source, e.GitRef, e.typeOf, e.pullNB); err != nil {
		return err
	}
	e.createdInDB = true

	e.Infra.SetTags(map[string]string{
		"execution_git_ref":         git.ShortenSHA(e.GitRef),
		"execution_source":          e.Source,
		"execution_type":            e.typeOf,
		"execution_planner_version": e.VtgatePlannerVersion,
	})

	err = e.prepareDirectories()
	if err != nil {
		return err
	}

	err = e.Infra.Prepare()
	if err != nil {
		return err
	}
	if e.configPath == "" {
		e.configPath = viper.ConfigFileUsed()
	}
	e.AnsibleConfig.ExtraVars = map[string]interface{}{}
	e.statsRemoteDBConfig.AddToAnsible(&e.AnsibleConfig)
	if e.pullNB != 0 {
		e.AnsibleConfig.ExtraVars["vitess_git_version_fetch_pr"] = "refs/pull/" + strconv.Itoa(e.pullNB) + "/head"
		e.AnsibleConfig.ExtraVars["vitess_git_version_pr_nb"] = e.pullNB
	}

	e.prepared = true
	return nil
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

	IPs, err := e.provision()
	if err != nil {
		return err
	}

	// TODO: optimize tokenization of Ansible files.
	err = ansible.AddIPsToFiles(IPs, e.AnsibleConfig)
	if err != nil {
		return err
	}
	err = ansible.AddLocalConfigPathToFiles(e.configPath, e.AnsibleConfig)
	if err != nil {
		return err
	}

	e.AnsibleConfig.ExtraVars[keyExecUUID] = e.UUID.String()
	e.AnsibleConfig.ExtraVars[keyVitessVersion] = e.GitRef
	e.AnsibleConfig.ExtraVars[keyExecSource] = e.Source
	e.AnsibleConfig.ExtraVars[keyExecutionType] = e.typeOf
	e.AnsibleConfig.ExtraVars[keyVtgatePlanner] = e.VtgatePlannerVersion

	// Infra will run the given config.
	err = e.Infra.Run(&e.AnsibleConfig)
	if err != nil {
		return err
	}
	return nil
}

func (e *Exec) Success() {
	// checking if the execution has not already failed
	rows, errSQL := e.clientDB.Select("SELECT uuid FROM execution WHERE uuid = ? AND status = ?", e.UUID.String(), StatusFailed)
	if errSQL != nil || rows.Next() {
		return
	}
	_, _ = e.clientDB.Insert("UPDATE execution SET finished_at = CURRENT_TIME, status = ? WHERE uuid = ?", StatusFinished, e.UUID.String())
}

func (e *Exec) handleStepEnd(err error) {
	if err != nil {
		_, _ = e.clientDB.Insert("UPDATE execution SET finished_at = CURRENT_TIME, status = ? WHERE uuid = ?", StatusFailed, e.UUID.String())
	}
}

func (e Exec) SendNotificationForRegression() (err error) {
	defer func() {
		e.handleStepEnd(err)
	}()

	var previousExec, previousGitRef string
	var ignoreNonRegression bool
	if e.pullNB > 0 {
		previousExec, previousGitRef, err = e.getPreviousForPR()
		ignoreNonRegression = true
	} else {
		previousExec, previousGitRef, err = e.getPreviousFromSameSource()
	}
	if err != nil {
		return err
	}
	if previousExec == "" {
		return nil
	}

	// regression header, appender to header in the event of a regression
	regressionHeader := `*Observed a regression.*
`

	// header of the message, before the regression explanation
	header := ``
	if e.pullNB > 0 {
		header += fmt.Sprintf(`Benchmarked PR #<https://github.com/vitessio/vitess/pull/%d>.`, e.pullNB)
	} else {
		header += `Comparing: recent commit <https://github.com/vitessio/vitess/commit/` + e.GitRef + `|` + git.ShortenSHA(e.GitRef) + `> with old commit <https://github.com/vitessio/vitess/commit/` + previousGitRef + `|` + git.ShortenSHA(previousGitRef) + `>.`
	}
	header += `Benchmark UUIDs, recent: ` + e.UUID.String()[:7] + ` old: ` + previousExec[:7] + `.
Comparison can be seen at : ` + getComparisonLink(e.GitRef, previousGitRef) + `

`

	if e.typeOf == "micro" {
		microBenchmarks, err := microbench.Compare(e.clientDB, e.GitRef, previousGitRef)
		if err != nil {
			return err
		}
		regression := microBenchmarks.Regression()
		err = e.sendMessageIfRegression(ignoreNonRegression, regression, header, regressionHeader)
		if err != nil {
			return err
		}
	} else if e.typeOf == "oltp" || e.typeOf == "tpcc" {
		influxCfg := influxdb.Config{
			Host:     e.statsRemoteDBConfig.Host,
			Port:     e.statsRemoteDBConfig.Port,
			User:     e.statsRemoteDBConfig.User,
			Password: e.statsRemoteDBConfig.Password,
			Database: e.statsRemoteDBConfig.DbName,
		}
		influxClient, err := influxCfg.NewClient()
		if err != nil {
			return err
		}

		macrosMatrices, err := macrobench.CompareMacroBenchmarks(e.clientDB, influxClient, e.GitRef, previousGitRef)
		if err != nil {
			return err
		}

		macroResults := macrosMatrices[macrobench.Type(e.typeOf)].(macrobench.ComparisonArray)
		if len(macroResults) == 0 {
			return fmt.Errorf("no macrobenchmark result")
		}

		regression := macroResults[0].Regression()
		err = e.sendMessageIfRegression(ignoreNonRegression, regression, header, regressionHeader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e Exec) sendMessageIfRegression(ignoreNonRegression bool, regression, header, regressionHeader string) error {
	if regression != "" || ignoreNonRegression {
		hd := header
		if regression != "" {
			hd = regressionHeader + header
		}
		err := e.sendSlackMessage(regression, hd)
		if err != nil {
			return err
		}
	}
	return nil
}

func getComparisonLink(leftSHA, rightSHA string) string {
	return "https://benchmark.vitess.io/compare?r=" + leftSHA + "&c=" + rightSHA
}

func (e Exec) sendSlackMessage(regression, header string) error {
	content := header + regression
	msg := slack.TextMessage{Content: content}
	err := msg.Send(e.slackConfig)
	if err != nil {
		return err
	}
	return nil
}

// CleanUp cleans and removes all things required only during the execution flow
// and not after it is done.
func (e Exec) CleanUp() (err error) {
	defer func() {
		e.handleStepEnd(err)
	}()
	err = e.Infra.CleanUp()
	if err != nil {
		return err
	}
	return nil
}

// NewExec creates a new *Exec with an autogenerated uuid.UUID as well
// as a constructed infra.Infra.
func NewExec() (*Exec, error) {
	// todo: dynamic choice for infra provider
	inf, err := construct.NewInfra(equinix.Name)
	if err != nil {
		return nil, err
	}

	ex := Exec{
		UUID:  uuid.New(),
		Infra: inf,

		// By default Exec prints os.Stdout and os.Stderr.
		// This can be changed later by explicitly using
		// exec.SetStdout and exec.SetStderr, or SetOutputToDefaultPath.
		stdout: os.Stdout,
		stderr: os.Stderr,

		configDB:   &mysql.ConfigDB{},
		clientDB:   nil,
		configPath: viper.ConfigFileUsed(),
	}

	// ex.AnsibleConfig.SetOutputs(ex.stdout, ex.stderr)
	ex.Infra.SetConfig(&ex.InfraConfig)
	ex.Infra.SetExecUUID(ex.UUID)

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

func (e Exec) getPreviousForPR() (execUUID, gitRef string, err error) {
	// getting the PR's base ref
	_, baseRef, err := git.GetPullRequestHeadAndBase(fmt.Sprintf("https://api.github.com/repos/vitessio/vitess/pulls/%d", +e.pullNB))
	if err != nil {
		return "", "", err
	}

	var result *sql.Rows
	if e.typeOf == "micro" {
		// compare with base ref
		query := "SELECT e.uuid, e.git_ref FROM execution e WHERE e.status = 'finished' AND " +
			"e.type = ? AND e.git_ref != ? ORDER BY e.started_at DESC LIMIT 1"
		result, err = e.clientDB.Select(query, e.typeOf, baseRef)
	} else {
		// compare with base ref and same vtgate_planner_version
		query := "SELECT e.uuid, e.git_ref FROM execution e, macrobenchmark m WHERE e.status = 'finished' AND " +
			"e.type = ? AND e.git_ref != ? AND e.uuid = m.exec_uuid AND m.vtgate_planner_version = ? ORDER BY e.started_at DESC LIMIT 1"
		result, err = e.clientDB.Select(query, e.typeOf, baseRef, e.VtgatePlannerVersion)
	}
	if err != nil {
		return
	}
	for result != nil && result.Next() {
		err = result.Scan(&execUUID, &gitRef)
		if err != nil {
			return
		}
	}
	return
}

func (e Exec) getPreviousFromSameSource() (execUUID, gitRef string, err error) {
	query := "SELECT e.uuid, e.git_ref FROM execution e WHERE e.source = ? AND e.status = 'finished' AND " +
		"e.type = ? AND e.git_ref != ? ORDER BY e.started_at DESC LIMIT 1"
	result, err := e.clientDB.Select(query, e.Source, e.typeOf, e.GitRef)
	if err != nil {
		return
	}
	for result.Next() {
		err = result.Scan(&execUUID, &gitRef)
		if err != nil {
			return
		}
	}
	return
}

// GetLatestCronJobForMicrobenchmarks will fetch and return the commit sha for which
// the last cron job for microbenchmarks was run
func GetLatestCronJobForMicrobenchmarks(client *mysql.Client) (gitSha string, err error) {
	query := "select git_ref from execution where source = \"cron\" and status = \"finished\" and type = \"micro\" order by started_at desc limit 1"
	rows, err := client.Select(query)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		err = rows.Scan(&gitSha)
		return gitSha, err
	}
	return "", nil
}

// GetLatestCronJobForMacrobenchmarks will fetch and return the commit sha for which
// the last cron job for macrobenchmarks was run
func GetLatestCronJobForMacrobenchmarks(client *mysql.Client) (gitSha string, err error) {
	query := "select git_ref from execution where source = \"cron\" and status = \"finished\" and ( type = \"oltp\" or type = \"tpcc\" ) order by started_at desc limit 1"
	rows, err := client.Select(query)
	if err != nil {
		return "", err
	}

	for rows.Next() {
		err = rows.Scan(&gitSha)
		return gitSha, err
	}
	return "", nil
}

func Exists(clientDB *mysql.Client, gitRef, source, typeOf, status string) (bool, error) {
	query := "SELECT uuid FROM execution WHERE status = ? AND git_ref = ? AND type = ? AND source = ?"
	result, err := clientDB.Select(query, status, gitRef, typeOf, source)
	if err != nil {
		return false, err
	}
	return result.Next(), nil
}

func ExistsMacrobenchmark(clientDB *mysql.Client, gitRef, source, typeOf, status, planner string) (bool, error) {
	query := "SELECT uuid FROM execution e, macrobenchmark m WHERE e.status = ? AND e.git_ref = ? AND e.type = ? AND e.source = ? AND m.vtgate_planner_version = ? AND e.uuid = m.exec_uuid"
	result, err := clientDB.Select(query, status, gitRef, typeOf, source, planner)
	if err != nil {
		return false, err
	}
	return result.Next(), nil
}
