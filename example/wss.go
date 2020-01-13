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

	fmt.Println("success")
	select {}
}

func getTableRows(cli *dfuse.Client) error {
	var err error
	f := func(respType string, payload dfuse.Payload) {
		switch respType {
		case dfuse.TableSnapshot:
			resp, err := payload.TableSnapshot()
			if err != nil {
				return
			}
			fmt.Printf("snapshot :%+v \n", resp)

		case dfuse.TableDelta:
			resp, err := payload.TableDelta()
			if err != nil {
				return
			}

			fmt.Printf("delta :%+v \n", resp)
			fmt.Printf("new :%+v \n", resp.Data.DBOP.New)
			fmt.Printf("old :%+v \n", resp.Data.DBOP.Old)

		case dfuse.Listening:
			resp, err := payload.Listening()
			if err != nil {
				return
			}
			fmt.Printf("listening :%+v \n", resp)

		case dfuse.Progress:
			resp, err := payload.Progress()
			if err != nil {
				return
			}
			fmt.Printf("progress :%+v \n", resp)

		case dfuse.Ping:
			fmt.Println("")

		case dfuse.Error:
			err := payload.Error()
			if err != nil {
				return
			}
			fmt.Printf("err :%+v \n", err)

		default:
			err = fmt.Errorf("invalid type :%s", respType)
		}
	}

	err = cli.Wss().GetTableRows(&entity.GetTableRows{
		Code:  "zheshimatch1",
		Scope: "gou2eos",
		Table: "oribuy",
		Json:  false,
		Limit: 10,
	}, f, &entity.OptionReq{
		Fetch:            true,
		Listen:           true,
		StartBlock:       0,
		IrreversibleOnly: false,
		WithProgress:     5,
	})
	return err
}

func getTransactionLifecycle() error {
	return nil
}
