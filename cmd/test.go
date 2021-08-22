/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"time"

	"eos-arb/repo/arbitrage"
	"eos-arb/repo/arbitrage/dex_adapter"
	"eos-arb/repo/arbitrage/models"
	"eos-arb/repo/arbitrage/printer"
	"eos-arb/repo/config"
	"eos-arb/repo/constants/obj_const"
	"eos-arb/repo/constants/str_const"
	"eos-arb/repo/db_util"
	"eos-arb/repo/eos_api"

	"github.com/cihub/seelog"
	"github.com/go-redis/redis/v8"
	"github.com/panyanyany/eos-go"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		cfg := config.LoadConfig()
		db := db_util.InitDb(cfg.Db.Name, cfg.Db.Pass)

		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		api := eos_api.NewEosApi(rdb, time.Hour*0)
		api.Actor = eos.AccountName(cfg.Accounts["default"].Name)
		api.PrivateKey = cfg.Accounts["default"].PrivateKey
		//api.EnableCache = false

		defibox := &dex_adapter.Defibox{
			BaseAdapter: dex_adapter.BaseAdapter{Api: api},
			Repo:        nil,
		}
		defis := &dex_adapter.Defis{
			BaseAdapter: dex_adapter.BaseAdapter{Api: api},
			Repo:        nil,
		}
		dolphin := &dex_adapter.Dolphin{
			BaseAdapter: dex_adapter.BaseAdapter{Api: api},
			Repo:        nil,
		}
		dexes := []dex_adapter.IDex{
			defibox,
			defis,
			dolphin,
		}
		for _, dex := range dexes {
			err = dex.Init(db)
			if err != nil {
				err = fmt.Errorf("dex.Init(): %w", err)
				seelog.Error(err)
				return
			}
		}

		arb := new(arbitrage.Repo)
		arb.Api = api
		arb.Printer = &printer.PathPrinter{ShowContract: false, ShowCmd: false, From: eos.Name(api.GetActor()), ShowDex: true}

		for {
			seelog.Debugf("start loop")
			pairs := arb.GetAllPairs(dexes)
			seelog.Debugf("pairs: %v", len(pairs))
			pairs = arb.FilterPairs(pairs)

			defiPairs := []*models.Pair{}
			for _, pair := range pairs {
				if pair.Dex == defis {
					defiPairs = append(defiPairs, pair)
				}
			}

			// defis-eos
			arb.RunTask(arbitrage.RunTaskInput{
				BaseSymbol:      obj_const.SEos,
				BaseContract:    eos.AccountName(str_const.CEosio),
				BaseAmount:      1,
				MinProfitAmount: 4,
				Pairs:           defiPairs,
			})
			// defis-usdt
			arb.RunTask(arbitrage.RunTaskInput{
				BaseSymbol:      obj_const.SUsdt,
				BaseContract:    eos.AccountName(str_const.CTether),
				BaseAmount:      1,
				MinProfitAmount: 15,
				Pairs:           defiPairs,
			})
			// defis-usdc
			//arb.RunTask(arbitrage.RunTaskInput{
			//	BaseSymbol:      obj_const.SUsdc,
			//	BaseContract:    eos.AccountName(str_const.CUsdx),
			//	BaseAmount:      51,
			//	MinProfitAmount: 204,
			//	Pairs:           defiPairs,
			//})

			// eos
			arb.RunTask(arbitrage.RunTaskInput{
				BaseSymbol:      obj_const.SEos,
				BaseContract:    eos.AccountName(str_const.CEosio),
				BaseAmount:      1000,
				MinProfitAmount: 4,
				Pairs:           pairs,
			})
			// usdt
			arb.RunTask(arbitrage.RunTaskInput{
				BaseSymbol:      obj_const.SUsdt,
				BaseContract:    eos.AccountName(str_const.CTether),
				BaseAmount:      2000,
				MinProfitAmount: 15,
				Pairs:           pairs,
			})
			// usn
			//arb.RunTask(arbitrage.RunTaskInput{
			//	BaseSymbol:      obj_const.SUsn,
			//	BaseContract:    eos.AccountName(str_const.CUsn),
			//	BaseAmount:      2000,
			//	MinProfitAmount: 15,
			//	Pairs:           pairs,
			//})
		}
	},
}

func RunDefis(arb *arbitrage.Repo, dex dex_adapter.IDex, pairs []*models.Pair, printer printer.PathPrinter) {
	var chances [][]*models.PathJob
	var err error

	//printer := printer.PathPrinter{ShowContract: false, ShowCmd: false, From: eos.Name(arb.Api.GetActor())}
	pairs = arb.FilterPairs(pairs)

	chances, err = arb.GetChances(arbitrage.GetChancesInput{
		Pairs:     pairs,
		PathDepth: 4,
		BaseAsset: eos.ExtendedAsset{
			Asset:    eos.Asset{Amount: 1, Symbol: obj_const.SEos},
			Contract: "eosio.token",
		},
		MinProfit: eos.ExtendedAsset{
			Asset:    eos.Asset{Amount: 1, Symbol: obj_const.SEos},
			Contract: "eosio.token",
		},
	})
	if err != nil {
		err = fmt.Errorf("arb.GetChances: %w", err)
		seelog.Error(err)
	}
	//for _, paths := range chances {
	//	printer.PrintPaths(paths)
	//}
	//seelog.Infof("")
	chances = arb.FilterChances(chances, eos.Asset{Amount: 4, Symbol: obj_const.SEos})
	for _, paths := range chances {
		printer.PrintPaths(paths)
	}
	//return
	minProfit := eos.Asset{Amount: 4, Symbol: obj_const.SEos}
	arb.RunChances(chances, minProfit)
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
