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

package slack

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagToken   = "slack-token"
	flagChannel = "slack-channel"
)

// Config used for Slack.
type Config struct {
	Token   string
	Channel string
}

// AddToCommand will add Config's CLI flags to the given *cobra.Command.
func (c *Config) AddToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.Token, flagToken, "", "Token used to authenticate Slack")
	cmd.Flags().StringVar(&c.Token, flagChannel, "", "Slack channel on which to post messages")

	_ = viper.BindPFlag(flagToken, cmd.Flags().Lookup(flagToken))
	_ = viper.BindPFlag(flagChannel, cmd.Flags().Lookup(flagChannel))
}
