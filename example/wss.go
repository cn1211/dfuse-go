package main

import (
	"fmt"
	"os"

	"github.com/chenyihui555/dfuse-go/entity"

	"github.com/chenyihui555/dfuse-go"
)

func main() {
	apiKey := os.Getenv("apikey")
	if apiKey == "" {
		panic("key is null")
	}

	opt := dfuse.Options{
		ApiKey:  apiKey,
		Network: dfuse.Jungle,
	}
	client := dfuse.NewClient(&opt)
	defer client.Wss().Close()

	if err := getTableRows(client); err != nil {
		fmt.Println("get table row err", err)
		return
	}

	select {}
}

func getTableRows(cli *dfuse.Client) error {
	var err error
	f := func(respType string, callback *dfuse.Callback) {
		switch respType {
		case dfuse.TableSnapshot:
			resp, err := callback.TableSnapshot()
			if err != nil {
				return
			}

			fmt.Printf("snapshot :%+v \n", resp)

		case dfuse.TableDelta:
			delta, err := callback.TableDelta()
			if err != nil {
				return
			}
			fmt.Printf("delta :%+v \n", delta.Data)
			fmt.Printf("key :%+v \n", delta.Data.DBOP.Key)
			fmt.Printf("new :%+v \n", delta.Data.DBOP.New)
			fmt.Printf("old :%+v \n", delta.Data.DBOP.Old)

		case dfuse.Listening:
			resp, err := callback.Listening()
			if err != nil {
				return
			}
			fmt.Printf("listening :%+v \n", resp)

		case dfuse.Progress:
			resp, err := callback.Progress()
			if err != nil {
				return
			}
			fmt.Printf("progress :%+v \n", resp)

		case dfuse.Error:
			resp, err := callback.Error()
			if err != nil {
				fmt.Printf("err :%+v \n", err)
				return
			}
			fmt.Printf("stream error info :%+v \n", resp)

		default:
			err = fmt.Errorf("invalid type %s", respType)
			fmt.Println("error", err)
		}
	}

	err = cli.Wss().GetTableRows("wss-ori-test", &entity.GetTableRows{
		Code:  "zheshimatch1",
		Scope: "pizza2usde",
		Table: "sell",
		Json:  true,
		Limit: 10,
	}, f, &entity.OptionReq{
		Fetch:            false,
		Listen:           true,
		StartBlock:       0,
		IrreversibleOnly: false,
		WithProgress:     10,
	})
	return err
}
