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
	f := func(respType, message string) {
		if err != nil {
			fmt.Println("err:", err)
			return
		}

		switch respType {
		case dfuse.TableSnapshot:
			resp, err := cli.Wss().TableSnapshot()
			if err != nil {
				return
			}

			fmt.Printf("snapshot :%+v \n", resp)

		case dfuse.TableDelta:
			resp, err := cli.Wss().TableDelta()
			if err != nil {
				return
			}

			fmt.Printf("delta :%+v \n", resp.Data)
			fmt.Printf("new :%+v \n", resp.Data.DBOP.New)
			fmt.Printf("old :%+v \n", resp.Data.DBOP.Old)

		//case dfuse.Listening:
		//	resp, err := cli.Wss().Listening()
		//	if err != nil {
		//		return
		//	}
		//	fmt.Printf("listening :%+v \n", resp)
		//
		//case dfuse.Progress:
		//	resp, err := cli.Wss().Progress()
		//	if err != nil {
		//		return
		//	}
		//	fmt.Printf("progress :%+v \n", resp)
		//
		//case dfuse.Ping:
		//	fmt.Println("")

		case dfuse.Error:
			err := cli.Wss().Error()
			fmt.Printf("err :%+v \n", err)

		default:
			err = fmt.Errorf("invalid type %s", respType)
			fmt.Println("error", err)
		}
	}

	err = cli.Wss().GetTableRows("wss-test", &entity.GetTableRows{
		Code:  "zheshimatch1",
		Scope: "pizza2usde",
		Table: "order",
		Json:  true,
		Limit: 10,
	}, f, &entity.OptionReq{
		Fetch:            false,
		Listen:           true,
		StartBlock:       0,
		IrreversibleOnly: false,
		WithProgress:     5,
	})
	return err
}
