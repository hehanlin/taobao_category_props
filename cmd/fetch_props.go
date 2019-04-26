// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/hehanlin/taobao_category/logic"
	"github.com/spf13/cobra"
)

// fetch_propsCmd represents the fetch_props command
var fetch_propsCmd = &cobra.Command{
	Use:   "fetch_props",
	Short: "传入淘宝发布类目id(cid), 返回该类目下的关键(描述)属性和销售(sku)属性",
	Long: `example: taobao_category fetch_props --cid 50010850 --token 7abb3b767e
	return: 一个包含所有该类目的属性csv文件`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := logic.Fetch_props(cid, token); err != nil {
			log.Fatal(err)
		}
	},
}

var (
	cid   int
	token string
)

func init() {
	RootCmd.AddCommand(fetch_propsCmd)

	fetch_propsCmd.PersistentFlags().IntVarP(&cid, "cid", "c", -1, "传入的类目id")
	fetch_propsCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "淘宝token, 登录淘宝账号后，cookie中的tb_token")

}
