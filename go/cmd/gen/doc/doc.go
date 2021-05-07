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

package doc

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func GenerateDoc() *cobra.Command {
	var docDir string
	cmd := &cobra.Command{
		Use:     "doc",
		Short:   "Generates documentation for the CLI",
		Long:    "Generates documentation for the CLI",
		Example: "arewefastyet gen doc --doc-dir ./arewefastyet/docs",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Generating documentation.")
			cmd.Root().DisableAutoGenTag = true
			err := doc.GenMarkdownTree(cmd.Root(), docDir)
			if err != nil {
				return err
			}
			log.Println("Documentation generated successfully.")
			return nil
		},
	}
	cmd.Flags().StringVar(&docDir, "doc-dir", "./docs", "Directory where the documentation will reside.")
	return cmd
}
