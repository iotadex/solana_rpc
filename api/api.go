package api

import (
	"context"
	"encoding/json"
	"net/http"
	"solana_rpc/config"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gin-gonic/gin"
)

// Send the test Token to user
func TokenAccountBalance(c *gin.Context) {
	client := rpc.New(config.SolanaRpc)
	owner := c.Query("account")
	collection := c.Query("collection")
	pubKey, err := solana.PublicKeyFromBase58(owner)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": ACCOUNT_ERROR,
			"err_msg":  "account error. " + c.Query("account"),
		})
		return
	}

	var data token.Account
	spl := c.Query("spl")
	if spl == "0" {
		var out *rpc.GetBalanceResult
		if out, err = client.GetBalance(context.TODO(), pubKey, "finalized"); err == nil {
			data.Amount = out.Value
			data.Owner = pubKey
		}
	} else if spl == "721" {
		data.Amount, err = GetNftAssert(owner, collection)
		data.Owner = pubKey
		data.Mint, _ = solana.PublicKeyFromBase58(collection)
	} else {
		err = client.GetAccountDataInto(context.TODO(), pubKey, &data)
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err_code": SYSTEM_ERROR,
			"err_msg":  "rpc error. " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":    true,
		"owner":     data.Owner.String(),
		"programid": data.Mint.String(),
		"amount":    data.Amount,
	})
}

type SearchAssertMethod struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		OwnerAddress string   `json:"ownerAddress"`
		Grouping     []string `json:"grouping"`
		Page         int      `json:"page"`
		Limit        int      `json:"limit"`
	} `json:"params"`
}

type SearchAssertResult struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      string `json:"id"`
	Result  struct {
		Total uint64 `json:"total"`
	}
}

var searchNftMethod SearchAssertMethod

func GetNftAssert(owner, collection string) (uint64, error) {
	grouping := make([]string, 0, 2)
	grouping = append(grouping, "collection")
	grouping = append(grouping, collection)
	searchNftMethod.Params.OwnerAddress = owner
	searchNftMethod.Params.Grouping = grouping

	res, err := HttpJsonPost(config.SolanaRpc, searchNftMethod)
	if err != nil {
		return 0, err
	}
	var nftResult SearchAssertResult
	if err = json.Unmarshal(res, &nftResult); err != nil {
		return 0, err
	}

	return nftResult.Result.Total, nil
}

func init() {
	searchNftMethod.Id = strconv.FormatInt(time.Now().Unix(), 10)
	searchNftMethod.Jsonrpc = "2.0"
	searchNftMethod.Method = "searchAssets"
	searchNftMethod.Params.Page = 1
	searchNftMethod.Params.Limit = 1
}
